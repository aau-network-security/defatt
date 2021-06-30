package cli

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"

	pb "github.com/aau-network-security/defat/app/daemon/proto"
	color "github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

const (
	NoTokenErrMsg     = "token contains an invalid number of segments"
	UnauthorizeErrMsg = "unauthorized"
)

var (
	UnreachableDaemonErr = errors.New("Daemon seems to be unreachable")
	UnauthorizedErr      = errors.New("You seem to not be logged in")
)

type Creds struct {
	Token    string
	Insecure bool
}

type Client struct {
	TokenFile string
	Token     string
	conn      *grpc.ClientConn
	rpcClient pb.DaemonClient
}

func Execute() {
	c, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	var rootCmd = &cobra.Command{Use: "defat"}
	rootCmd.AddCommand(
		c.StartGame(),
		c.ListScenarios(),
		c.ListChallengesOnScenario(),
		c.CmdUser(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (c Creds) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"token": string(c.Token),
	}, nil
}

func (c Creds) RequireTransportSecurity() bool {
	return !c.Insecure
}

func NewClient() (*Client, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("Unable to find home directory")
	}

	tokenFile := filepath.Join(usr.HomeDir, ".defat_token")
	c := &Client{
		TokenFile: tokenFile,
	}

	if _, err := os.Stat(tokenFile); err == nil {
		if err := c.LoadToken(); err != nil {
			return nil, err
		}
	}

	host := os.Getenv("DEFAT_HOST")
	//todo i have change it for testing purpose
	if host == "" {
		host = "sec03.lab.es.aau.dk"
	}

	port := os.Getenv("DEFAT_PORT")
	if port == "" {
		port = "5454"
	}

	authCreds := Creds{Token: c.Token}
	dialOpts := []grpc.DialOption{}

	ssl_off := os.Getenv("DEFAT_SSL_OFF")
	endpoint := fmt.Sprintf("%s:%s", host, port)
	var creds credentials.TransportCredentials
	if strings.ToLower(ssl_off) == "true" {
		authCreds.Insecure = true
		dialOpts = append(dialOpts,
			grpc.WithInsecure(),
			grpc.WithPerRPCCredentials(authCreds))
	} else {
		if host == "localhost" {
			devCertPool := x509.NewCertPool()
			creds = setCertConfig(devCertPool)
		} else {
			certPool, _ := x509.SystemCertPool()
			creds = setCertConfig(certPool)
		}
		dialOpts = append(dialOpts,
			grpc.WithTransportCredentials(creds),
			grpc.WithPerRPCCredentials(authCreds))
	}

	conn, err := grpc.Dial(endpoint, dialOpts...)
	if err != nil {
		return nil, err
	}
	c.conn = conn
	c.rpcClient = pb.NewDaemonClient(conn)

	return c, nil
}

func setCertConfig(certPool *x509.CertPool) credentials.TransportCredentials {
	creds := credentials.NewTLS(&tls.Config{})
	creds = credentials.NewTLS(&tls.Config{
		RootCAs: certPool,
	})
	return creds
}

func (c *Client) LoadToken() error {
	raw, err := ioutil.ReadFile(c.TokenFile)
	if err != nil {
		return err
	}

	c.Token = string(raw)
	return nil
}

func (c *Client) SaveToken() error {
	return ioutil.WriteFile(c.TokenFile, []byte(c.Token), 0644)
}

func (c *Client) Close() {
	c.conn.Close()
}

// Downloads necessary localhost certificates
// for local development
func downloadCerts(certMap map[string]string) error {
	_, err := os.Stat("localcerts")
	if os.IsNotExist(err) {
		errDir := os.MkdirAll("localcerts", 0755)
		if errDir != nil {
			log.Fatal(err)
		}
		for k, v := range certMap {
			// Get the data
			resp, err := http.Get(v)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			// Create the file

			out, err := os.Create("localcerts/" + k)
			if err != nil {
				return err
			}

			defer out.Close()
			// Write the body to file
			_, err = io.Copy(out, resp.Body)
			if err != nil {
				return err
			}
		}
	}
	return nil

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
			return UnreachableDaemonErr
		}

		return err
	}

	return err
}

func PrintError(err error) {
	err = TranslateRPCErr(err)
	fmt.Printf("%s %s\n", color.Red("<!>"), color.Red(err.Error()))
}

func PrintWarning(s string) {
	fmt.Printf("%s %s\n", color.Brown("<?>"), color.Brown(s))
}

func ReadSecret(inputHint string) (string, error) {
	fmt.Printf(inputHint)
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Printf("\n")
	if err != nil {
		return "", err
	}

	return string(bytePassword), nil
}
