package controller

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"os/exec"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/aau-network-security/openvswitch/ovs"
)

type NetController struct {
	// IPService is wrapping functionalities of `ip` tool
	IPService *IPService

	// IfConfigServce is wrapping functionalities of `ifconfig` tool
	IFConfig *IFConfigService

	// Ovs Client
	Ovs *ovs.Client

	// IpPool

	IPPool *IPPool

	// Used to enable root command
	sudo bool

	// flags to service
	flags []string

	// enable debug or not
	debug bool

	// Implementation of ExecFunc.
	execFunc ExecFunc

	// Implementation of PipeFunc.
	pipeFunc PipeFunc
}

type Error struct {
	Out []byte
	Err error
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Err, string(e.Out))
}

func (c *NetController) exec(cmd string, args ...string) ([]byte, error) {
	flags := append(c.flags, args...)

	// If needed, prefix sudo.
	if c.sudo {
		flags = append([]string{cmd}, flags...)
		cmd = "sudo"
	}

	c.debugf("exec %s %v", cmd, flags)
	out, err := c.execFunc(cmd, flags...)
	if out != nil {
		out = bytes.TrimSpace(out)
		c.debugf("exec: %q", string(out))
	}
	if err != nil {
		// Wrap errors in Error type for further introspection
		return nil, &Error{
			Out: out,
			Err: err,
		}
	}
	return out, nil
}

func (c *NetController) debugf(format string, i ...interface{}) {
	if !c.debug {
		return
	}
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Debug().Msgf("defat: "+format, i...)
}

//
type ExecFunc func(cmd string, args ...string) ([]byte, error)

// An OptionFunc is a function which can apply configuration to a NetController.
type OptionFunc func(c *NetController)

func New(options ...OptionFunc) *NetController {
	// Always execute and pipe using shell when created with New.
	c := &NetController{
		flags:    make([]string, 0),
		execFunc: shellExec,
		debug:    true,
		sudo:     true,
		pipeFunc: shellPipe,
	}
	for _, o := range options {
		o(c)
	}

	ip := &IPService{
		c: c,
	}
	ifconf := &IFConfigService{
		c: c,
	}

	c.Ovs = ovs.New(ovs.Sudo(), ovs.Debug(false))

	c.IPService = ip

	c.IFConfig = ifconf

	return c
}

// Sudo specifies that "sudo" should be prefixed to all controller commands.
func Sudo() OptionFunc {
	return func(c *NetController) {
		c.sudo = true
	}
}

// shellPipe is a PipeFunc which shells out to the binary cmd using the arguments
// args, and writing to the command's stdin using stdin.
func shellPipe(stdin io.Reader, cmd string, args ...string) ([]byte, error) {
	command := exec.Command(cmd, args...)

	stdout, err := command.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := command.StderrPipe()
	if err != nil {
		return nil, err
	}

	wc, err := command.StdinPipe()
	if err != nil {
		return nil, err
	}

	if err := command.Start(); err != nil {
		return nil, err
	}

	if _, err := io.Copy(wc, stdin); err != nil {
		return nil, err
	}

	// Reference: https://golang.org/pkg/os/exec/#Cmd.StdinPipe
	if err := wc.Close(); err != nil {
		return nil, err
	}

	mr := io.MultiReader(stdout, stderr)
	b, err := ioutil.ReadAll(mr)
	if err != nil {
		return nil, err
	}

	return b, command.Wait()
}

// shellExec is an ExecFunc which shells out to the binary cmd using the
// arguments args, and returns its combined stdout and stderr and any errors
// which may have occurred.
func shellExec(cmd string, args ...string) ([]byte, error) {
	return exec.Command(cmd, args...).CombinedOutput()
}

// A PipeFunc is a function which accepts an input stdin stream, command,
// and arguments, and returns command output and an error.
type PipeFunc func(stdin io.Reader, cmd string, args ...string) ([]byte, error)
