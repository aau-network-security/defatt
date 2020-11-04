package controller

import (
	"math/rand"
	"net"
	"regexp"
	"time"

	"github.com/rs/zerolog/log"
)

type IPTables struct {
	c *NetClient
}

//echo "Setting firewall rules"

//TODO: for the Defatt
//	iptables -P FORWARD DROP
//	iptables -A FORWARD -i enp0s8 -o enp0s9 -j ACCEPT
//	iptables -A FORWARD -i enp0s9 -o enp0s8 -j ACCEPT
//	iptables -A FORWARD -i enp0s8 -o enp0s10 -j ACCEPT
//	iptables -A FORWARD -i enp0s10 -o enp0s8 -m state ! --state NEW -j ACCEPT

var (
	//Name of the chaines
	inputC   = Chain("INPUT")
	forwardC = Chain("FORWARD")
	outputC  = Chain("OUTPUT")

	//name of the policy
	acceptP = Policy("ACCEPT")
	dropP   = Policy("DROP")

	appendA  = Action("-A") // append action
	flushA   = Action("-F") // flush action
	deleteA  = Action("-D") // delete action
	insertA  = Action("-I") // insert action
	defaultA = Action("-P") // set default rule

)

type Action string

type Chain string

type Policy string

//TODO: Add logging messages

//drop all rules from selected chain
func (ipTab *IPTables) DropExistingRule(c Chain) error {
	//iptables -F FORWARD
	log.Debug().Msgf("Dropping all rules from the %s chain", c)
	cmds := []string{string(flushA), string(c)}
	_, err := ipTab.c.IPTables.exec(cmds...)
	return err

}

func (ipTab *IPTables) SetDefaultRule(c Chain, p Policy) error {
	log.Debug().Msgf("Setting default policy to %s in  %s chain", c, p)
	//	iptables -P FORWARD DROP
	cmds := []string{string(defaultA), string(c), string(p)}
	_, err := ipTab.exec(cmds...)
	return err
}

//iptables -A FORWARD -i enp0s8 -o enp0s9 -j ACCEPT
func (ipTab *IPTables) SetAcceptRule(in, out string) error {
	log.Debug().Msgf("Setting in %s chain to forward traffic comming from %s to %s", in, out)
	cmds := []string{string(appendA), string(forwardC), "-i", in, "-o", out, "-j", string(acceptP)}
	_, err := ipTab.c.IPTables.exec(cmds...)
	return err
}

//iptables -A FORWARD -i enp0s10 -o enp0s8 -m state ! --state NEW -j ACCEPT

func (ipTab *IPTables) CheckWhoCreatesConn(in, out string) error {
	log.Debug().Msgf("Allow connection that is not new between %s and %s", in, out)
	cmds := []string{string(appendA), string(forwardC), "-i", in, "-o", out, "-m", "state", "!", "--state", "NEW", "-j", string(acceptP)}
	_, err := ipTab.exec(cmds...)
	return err
}

// exec executes an ExecFunc using 'iptables'.
func (ipTab *IPTables) exec(args ...string) ([]byte, error) {
	return ipTab.c.exec("iptables", args...)
}

func GetSystemInterfaces() ([]string, error) {
	var interfaces []string
	ifaces, err := net.Interfaces()
	//regex to find the wireguard interface
	re := regexp.MustCompile("wg")

	for _, value := range ifaces {
		//avoid wg interface
		if re.MatchString(value.Name) == true {
			continue
		}
		interfaces = append(interfaces, value.Name)
	}

	return interfaces, err
}

func PickRandomInterface() string {
	var ifaceName string
	var randomIndex int
	getInterfaceName, err := GetSystemInterfaces()
	if err == nil {
		//select everytime different interface for system
		rand.Seed(time.Now().UnixNano())
		min := 1
		max := len(getInterfaceName)
		randomIndex = rand.Intn(max-min+1) + min
		ifaceName = getInterfaceName[randomIndex]
	}

	return ifaceName
}

//TODO: how to make it persistent
//TODO: echo "net.ipv4.ip_forward=1" >> /etc/sysctl.conf
//		sysctl -p
