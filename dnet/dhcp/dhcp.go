package dhcp

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"strconv"

	"github.com/aau-network-security/defatt/dnet/dhcp/proto"
	"github.com/aau-network-security/defatt/dnet/wg"
	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	AUTH_KEY = "wg"
)

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

func NewDHCPClient(ctx context.Context, wgConn wg.WireGuardConfig, port uint) (proto.DHCPClient, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		AUTH_KEY: wgConn.AuthKey,
	})
	tokenString, err := token.SignedString([]byte(wgConn.SignKey))
	if err != nil {
		return nil, err
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
			grpc.WithBlock(),
		}

		conn, err := grpc.DialContext(ctx, wgConn.Endpoint+":"+strconv.FormatUint(uint64(port), 10), dialOpts...)
		if err != nil {
			log.Error().Msgf("Error on dialing vpn service: %v", err)
			return nil, err
		}
		c := proto.NewDHCPClient(conn)
		return c, nil
	}

	authCreds.Insecure = true
	conn, err := grpc.DialContext(ctx, wgConn.Endpoint+":"+strconv.FormatUint(uint64(port), 10), grpc.WithInsecure(), grpc.WithPerRPCCredentials(authCreds), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	c := proto.NewDHCPClient(conn)
	return c, nil
}
