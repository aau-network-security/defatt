package daemon

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"

	pb "github.com/aau-network-security/defatt/app/daemon/proto"
	"github.com/aau-network-security/defatt/config"
	"github.com/aau-network-security/defatt/controller"
	"github.com/aau-network-security/defatt/database"
	vpn "github.com/aau-network-security/defatt/dnet/wg"
	"github.com/aau-network-security/defatt/frontend"
	"github.com/aau-network-security/defatt/game"
	"github.com/aau-network-security/defatt/store"
	"github.com/aau-network-security/defatt/virtual/docker"
	"github.com/aau-network-security/defatt/virtual/vbox"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var (
	PortIsAllocatedError = errors.New("Given gRPC port is already allocated")
	GrpcOptsErr          = errors.New("failed to retrieve server options")
	reTag                = regexp.MustCompile(`^[a-z]{4}$`)
	version              string
	displayTimeFormat    = time.RFC3339
)

type daemon struct {
	config     *config.Config
	auth       Authenticator
	users      store.UsersFile
	closers    []io.Closer
	vlib       vbox.Library
	web        *frontend.Web
	controller *controller.NetController
	games      map[store.Tag]game.Game
	pb.UnimplementedDaemonServer
}

type MissingConfigErr struct {
	Option string
}

type MngtPortErr struct {
	port string
}

type contextStream struct {
	grpc.ServerStream
	ctx context.Context
}

func New(conf *config.Config) (*daemon, error) {

	uf, err := store.NewUserFile(conf.DefatConfig.UsersFile)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to read users file: %s", conf.DefatConfig.UsersFile))
	}
	vlib := vbox.NewLibrary(conf.VmConfig.OvaDir)

	if len(uf.ListUsers()) == 0 && len(uf.ListSignupKeys()) == 0 {
		k := store.NewSignupKey()
		k.WillBeSuperUser = true

		if err := uf.CreateSignupKey(k); err != nil {
			return nil, err
		}

		log.Info().Msg("No users or signup keys found, creating a key")
	}

	web, err := frontend.New(fmt.Sprintf(":%d", conf.DefatConfig.FrontendPort), fmt.Sprintf(":%d", conf.DefatConfig.FrontendPortTLS), conf.DefatConfig.Endpoint, conf.DefatConfig.CertConf.CertFile, conf.DefatConfig.CertConf.CertKey)
	if err != nil {
		return nil, err
	}

	database.New(context.TODO(), conf.DefatConfig.DatabaseFile)

	contr := controller.New()

	keys := uf.ListSignupKeys()
	if len(uf.ListUsers()) == 0 && len(keys) > 0 {
		log.Info().Msg("No users found, printing keys")
		for _, k := range keys {
			log.Info().Str("key", k.String()).Msg("Found key")
		}
	}

	return &daemon{
		config:     conf,
		auth:       NewAuthenticator(uf, conf.DefatConfig.SigningKey),
		users:      uf,
		closers:    []io.Closer{},
		vlib:       vlib,
		controller: contr,
		web:        web,
	}, nil
}
func (m *MissingConfigErr) Error() string {
	return fmt.Sprintf("%s cannot be empty", m.Option)
}

func (m *MngtPortErr) Error() string {
	return fmt.Sprintf("failed to listen on management port %s", m.port)
}

func (d *daemon) Run() error {
	defer database.Close()
	go func() {
		if err := d.web.Run(); err != nil {
			log.Error().Err(err).Msg("error while running frontend")
		}
	}()

	gRPCPort := fmt.Sprintf(":%d", d.config.DefatConfig.Port)
	// start gRPC daemon
	lis, err := net.Listen("tcp", gRPCPort)
	if err != nil {
		return &MngtPortErr{gRPCPort}
	}
	log.Info().Str("port", gRPCPort).Msg("gRPC daemon has been started!")

	opts, err := d.grpcOpts()
	if err != nil {
		return errors.Wrap(GrpcOptsErr, err.Error())
	}
	s := d.GetServer(opts...)
	pb.RegisterDaemonServer(s, d)

	reflection.Register(s)
	log.Info().Msg("Reflection Registration is called.... ")

	return s.Serve(lis)
}

func (d *daemon) Close() error {
	var errs error
	var wg sync.WaitGroup

	for _, c := range d.closers {
		wg.Add(1)
		go func(c io.Closer) {
			if err := c.Close(); err != nil && errs == nil {
				errs = err
			}
			wg.Done()
		}(c)
	}

	wg.Wait()

	if err := docker.DefaultLinkBridge.Close(); err != nil {
		return err
	}

	return errs
}

func (d *daemon) GetServer(opts ...grpc.ServerOption) *grpc.Server {
	nonAuth := []string{"LoginUser", "SignupUser"}
	var logger *zerolog.Logger

	streamInterceptor := func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx, authErr := d.auth.AuthenticateContext(stream.Context())
		ctx = withAuditLogger(ctx, logger)
		stream = &contextStream{stream, ctx}

		header := metadata.Pairs("daemon-version", version)
		stream.SendHeader(header)

		for _, endpoint := range nonAuth {
			if strings.HasSuffix(info.FullMethod, endpoint) {
				return handler(srv, stream)
			}
		}

		if authErr != nil {
			return authErr
		}

		return handler(srv, stream)
	}

	unaryInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx, authErr := d.auth.AuthenticateContext(ctx)
		ctx = withAuditLogger(ctx, logger)

		header := metadata.Pairs("daemon-version", version)
		grpc.SendHeader(ctx, header)

		for _, endpoint := range nonAuth {
			if strings.HasSuffix(info.FullMethod, endpoint) {
				return handler(ctx, req)
			}
		}

		if authErr != nil {
			return nil, authErr
		}

		return handler(ctx, req)
	}

	opts = append([]grpc.ServerOption{
		grpc.StreamInterceptor(streamInterceptor),
		grpc.UnaryInterceptor(unaryInterceptor),
	}, opts...)
	return grpc.NewServer(opts...)
}

func withAuditLogger(ctx context.Context, logger *zerolog.Logger) context.Context {
	if logger == nil {
		return ctx
	}

	u, ok := ctx.Value(us{}).(store.User)
	if !ok {
		return logger.WithContext(ctx)
	}

	ls := logger.With().
		Str("user", u.Username).
		Bool("is-super-user", u.SuperUser).
		Logger()
	logger = &ls

	return logger.WithContext(ctx)
}

func (s *contextStream) Context() context.Context {
	return s.ctx
}

func combineErrors(errors []error) []string {
	var errorString []string
	for _, e := range errors {
		errorString = append(errorString, e.Error())
	}
	return errorString
}

func (d *daemon) CreateGame(ctx context.Context, req *pb.CreateGameRequest) (*pb.CreateGameResponse, error) {
	if err := d.createGame(req.Tag, req.Name, int(req.ScenarioNo)); err != nil {
		return &pb.CreateGameResponse{}, err
	}
	return &pb.CreateGameResponse{Message: "Game is created"}, nil
}
func (d *daemon) StopGame(ctx context.Context, req *pb.StopGameRequest) (*pb.StopGameResponse, error) {

	game := req.Tag

	g, ok := d.games[store.Tag(game)]
	if !ok {
		return &pb.StopGameResponse{}, nil
	}
	if err := g.Close(); err != nil {
		return &pb.StopGameResponse{}, nil
	}
	return &pb.StopGameResponse{Message: fmt.Sprintf("Game [ %s ] is closed ", game)}, nil
}
func (d *daemon) ListGames(ctx context.Context, req *pb.EmptyRequest) (*pb.ListGamesResponse, error) {
	// todo: List Running games

	return nil, status.Errorf(codes.Unimplemented, "method ListGames not implemented")
}

func (d *daemon) ListScenarios(ctx context.Context, req *pb.EmptyRequest) (*pb.ListScenariosResponse, error) {
	var respScenarios []*pb.ListScenariosResponse_Scenario
	scenarios := store.GetAllScenarios()

	for _, v := range scenarios {
		var scenario pb.ListScenariosResponse_Scenario
		scenario.Id = v.ID
		scenario.Duration = v.Duration
		scenario.Difficulty = v.Difficulty
		scenario.Story = v.Story
		for k, value := range v.Networks {
			var network pb.Subnet
			network.Vlan = k
			network.Challenges = value.Chals
			scenario.Networks = append(scenario.Networks, &network)
		}

		respScenarios = append(respScenarios, &scenario)
	}

	return &pb.ListScenariosResponse{Scenarios: respScenarios}, nil
}

func (d *daemon) createGame(tag, name string, sceanarioNo int) error {
	wgConfig := d.config.WireguardService
	if !reTag.MatchString(tag) {
		return status.Errorf(codes.InvalidArgument, "Gametag does not follow allowed convention - should only be four lowercase letters")
	}

	gameConf := game.GConfig{
		ID:   uuid.New().String()[0:8],
		Name: name,
		Tag:  tag,
		Host: d.config.DefatConfig.Endpoint,
		WgConfig: vpn.WireGuardConfig{
			Endpoint: wgConfig.Endpoint,
			Port:     wgConfig.Port,
			AuthKey:  wgConfig.AuthKey,
			SignKey:  wgConfig.SignKey,
		},
	}

	env, err := game.NewEnvironment(&gameConf, d.vlib)
	if err != nil {
		return err
	}

	if err := env.Start(context.TODO(), tag, name, sceanarioNo); err != nil {
		return err
	}

	d.web.AddGame(env)

	return nil

}

func (d *daemon) ListScenChals(ctx context.Context, req *pb.ListScenarioChallengesReq) (*pb.ListScenarioChallengesResp, error) {
	// nt := store.GetScenarioByID(int(req.ScenarioId))
	// var networks []*pb.Network
	// for _, r := range nt.Networks {
	// 	networks = append(networks, &pb.Network{
	// 		Challenges: r.Chals,
	// 		Vlan:       r.Vlan,
	// 	})
	// }
	return &pb.ListScenarioChallengesResp{Chals: nil}, nil
}

func (d *daemon) grpcOpts() ([]grpc.ServerOption, error) {
	if d.config.DefatConfig.CertConf.Enabled {
		// Load cert pairs
		certificate, err := tls.LoadX509KeyPair(d.config.DefatConfig.CertConf.CertFile, d.config.DefatConfig.CertConf.CertKey)
		if err != nil {
			return nil, fmt.Errorf("could not load server key pair: %s", err)
		}

		// Create a certificate pool from the certificate authority
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(d.config.DefatConfig.CertConf.CAFile)
		if err != nil {
			return nil, fmt.Errorf("Defatt Grpc could not read ca certificate: %s", err)
		}
		// CA file for let's encrypt is located under domain conf as `chain.pem`
		// pass chain.pem location
		// Append the client certificates from the CA
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			return nil, errors.New("failed to append client certs")
		}

		// Create the TLS credentials
		creds := credentials.NewTLS(&tls.Config{
			// no need to RequireAndVerifyClientCert
			Certificates: []tls.Certificate{certificate},
			ClientCAs:    certPool,
		})

		return []grpc.ServerOption{grpc.Creds(creds)}, nil
	}
	return []grpc.ServerOption{}, nil
}
