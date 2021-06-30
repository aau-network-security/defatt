package wg

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"strconv"
	"strings"

	vpn "github.com/aau-network-security/defatt/app/daemon/vpn-proto"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

var (
	UnreachableVPNServiceErr = errors.New("Wireguard service is not running !")
	UnauthorizedErr          = errors.New("Unauthorized attempt to use VPN service ")
	NoTokenErrMsg            = "token contains an invalid number of segments"
	UnauthorizeErrMsg        = "unauthorized"
	AUTH_KEY                 = "wg"
)

type WireGuardConfig struct {
	Endpoint string
	Port     uint64
	AuthKey  string
	SignKey  string
	Enabled  bool
	CertFile string
	CertKey  string
	CAFile   string
	Dir      string // client configuration file will reside
}

type Creds struct {
	Token    string
	Insecure bool
}

func (c Creds) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"token": string(c.Token),
	}, nil
}

func (c Creds) RequireTransportSecurity() bool {
	return !c.Insecure
}

func NewGRPCVPNClient(wgConn WireGuardConfig) (vpn.WireguardClient, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		AUTH_KEY: wgConn.AuthKey,
	})
	tokenString, err := token.SignedString([]byte(wgConn.SignKey))
	if err != nil {
		return nil, TranslateRPCErr(err)
	}

	authCreds := Creds{Token: tokenString}

	if wgConn.Enabled {
		log.Debug().Bool("TLS", wgConn.Enabled).Msg(" secure connection enabled for creating secure db client")

		// Load the client certificates from disk
		certificate, err := tls.LoadX509KeyPair(wgConn.CertFile, wgConn.CertKey)
		log.Info().Str("Certfile", wgConn.CertFile).
			Str("CertKey", wgConn.CertKey).Msg("Certs files")
		if err != nil {
			log.Printf("could not load client key pair: %s", err)
		}

		// Create a certificate pool from the certificate authority
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(wgConn.CAFile)
		if err != nil {
			log.Printf("VPNCONN could not read ca certificate: %s", err)
		}

		// Append the certificates from the CA
		// This is chain.pem for letsencrypt
		// can be found in same place with existing certificates
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			log.Error().Msg("failed to append ca certs")
		}

		creds := credentials.NewTLS(&tls.Config{
			// no need to give specific Grpc address
			// if it is given certificates should be generated
			// per given address
			Certificates: []tls.Certificate{certificate},
			RootCAs:      certPool,
		})

		dialOpts := []grpc.DialOption{
			grpc.WithTransportCredentials(creds),
			grpc.WithPerRPCCredentials(authCreds),
		}

		conn, err := grpc.Dial(wgConn.Endpoint+":"+strconv.FormatUint(wgConn.Port, 10), dialOpts...)
		if err != nil {
			log.Error().Msgf("Error on dialing vpn service: %v", err)
			return nil, TranslateRPCErr(err)
		}
		c := vpn.NewWireguardClient(conn)
		return c, nil
	}

	authCreds.Insecure = true
	conn, err := grpc.Dial(wgConn.Endpoint+":"+strconv.FormatUint(wgConn.Port, 10), grpc.WithInsecure(), grpc.WithPerRPCCredentials(authCreds))
	if err != nil {
		return nil, TranslateRPCErr(err)
	}
	c := vpn.NewWireguardClient(conn)
	return c, nil
}

func TranslateRPCErr(err error) error {
	st, ok := status.FromError(err)
	if ok {
		msg := st.Message()
		switch {
		case UnauthorizeErrMsg == msg:
			return UnauthorizedErr

		case NoTokenErrMsg == msg:
			return UnauthorizedErr

		case strings.Contains(msg, "TransientFailure"):

			return UnreachableVPNServiceErr
		}

		return err
	}

	return err
}
