package dhcp

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"strings"
	"text/template"

	"github.com/aau-network-security/defat/controller"
	"github.com/aau-network-security/defat/virtual/docker"
	"github.com/aau-network-security/openvswitch/ovs"
	"github.com/rs/zerolog/log"
)

type Networks struct {
	Subnets    []Subnet
	DNS        string
	MAC        string
	FixAddress string
}

type Subnet struct {
	Interface string
	Vlan      string
	Network   string
	Min       string
	Max       string
	Router    string
}

type Server struct {
	cont       docker.Container
	confFile   string
	ipList     map[string]string
	macAddress string
}

type LanSpec struct {
	NetI   string
	LANTag string
	Bridge string
}

func createDHCPFile(nets Networks) string {

	var tpl bytes.Buffer

	dir, err := os.Getwd() // get working directory
	if err != nil {
		log.Error().Msgf("Error getting the working dir %v", err)
	}
	fullPathToTemplate := fmt.Sprintf("%s%s", dir, "/dnet/dhcp/dhcpd.conf.tmpl")

	tmpl := template.Must(template.ParseFiles(fullPathToTemplate))

	//tmpl := template.Must(template.ParseFiles("/home/ubuntu/vlad/sec03/defatt/dnet/dhcp/dhcpd.conf.tmpl"))

	tmpl.Execute(&tpl, nets)
	return tpl.String()
}

func addToSwitch(c *controller.NetController, net Subnet, bridge, cid string) error {
	if err := c.Ovs.Docker.AddPort(bridge, net.Interface, cid,
		// exclusive for dhcp
		ovs.DockerOptions{
			IPAddress: fmt.Sprintf("%s/24", net.Router),
		}); err != nil {
		log.Error().Msgf("Error on ovs-docker addport %v", err)
		return err
	}

	if err := c.Ovs.Docker.SetVlan(bridge, net.Interface, cid, net.Vlan); err != nil {
		log.Error().Msgf("Error on ovs-docker SetVlan %v", err)
		return err
	}

	return nil
}

func generateMac() net.HardwareAddr {
	buf := make([]byte, 6)
	var mac net.HardwareAddr

	_, err := rand.Read(buf)
	if err != nil {
	}

	// Set the local bit
	buf[0] |= 2

	mac = append(mac, buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])

	return mac
}

func stringHardwareAddress(hardwareAddr net.HardwareAddr) string {
	s := hardwareAddr.String()
	if len(s) == 0 {
		return "-"
	}
	return s
}

//New creates a DHCP server which will be listening on the interfaces given as the argument
func New(ctx context.Context, ifaces map[string]string, bridge string, c *controller.NetController) (*Server, error) {
	ipList := make(map[string]string)
	var networks Networks
	ipPool := controller.NewIPPoolFromHost()
	var sNet Subnet

	//var macAddress net.HardwareAddr
	macAddress := generateMac()
	//macAddressString := string(macAddress)
	macAddressString := stringHardwareAddress(macAddress)
	//
	for vl, vt := range ifaces {

		//make and if and Else Here and put

		if strings.Contains(vl, "monitor") {

			sNet.Interface = vl
			sNet.Vlan = vt
			sNet.Network = "10.10.10" + ".0"
			sNet.Min = "10.10.10" + ".6"
			sNet.Max = "10.10.10" + ".180"
			sNet.Router = "10.10.10" + ".1"
			networks.Subnets = append(networks.Subnets, sNet)

			networks.FixAddress = "10.10.10.200"
			networks.MAC = macAddressString
			continue

			//networks.MAC = ""
			//ipList[sNet.Vlan] = randIP

		}

		randIP, _ := ipPool.Get()
		sNet.Interface = vl
		sNet.Vlan = vt
		sNet.Network = randIP + ".0"
		sNet.Min = randIP + ".6"
		sNet.Max = randIP + ".254"
		sNet.Router = randIP + ".1"
		networks.Subnets = append(networks.Subnets, sNet)
		ipList[sNet.Vlan] = randIP
	}

	f, err := ioutil.TempFile("", "dhcpd-conf")
	if err != nil {
		return nil, err
	}
	confFile := f.Name()

	// todo: Assign values for MAC and Fixed Address
	// networks.FixAddress
	//  networks.MAC

	confStr := createDHCPFile(networks)
	_, err = f.WriteString(confStr)
	if err != nil {
		return nil, err
	}
	cont := docker.NewContainer(docker.ContainerConfig{
		Image: "lanestolen/dhcp",
		Mounts: []string{
			fmt.Sprintf("%s:/etc/dhcp/dhcpd.conf", confFile),
		},
		UsedPorts: []string{"67/udp"},
		Resources: &docker.Resources{
			MemoryMB: 50,
			CPU:      0.3,
		},
		Labels: map[string]string{
			"nap": "lan_dhcpd",
		},
		UseBridge: false,
	})
	if err := cont.Create(ctx); err != nil {
		log.Error().Msgf("Error in creating container  %v", err)
		return nil, err
	}
	if err := cont.Start(ctx); err != nil {
		log.Error().Msgf("Error in starting container  %v", err)
	}
	cid := cont.ID()
	for _, net := range networks.Subnets {
		if err := addToSwitch(c, net, bridge, cid); err != nil {
			log.Error().Msgf("Error on addToSwitch in dhcp %v ", err)
		}
		log.Info().Str("Interface", net.Interface).
			Str("Vlan", net.Vlan).
			Str("Network", net.Network).
			Str("Router", net.Router).
			Msgf("DHCP server initialized")

	}

	return &Server{
		cont:       cont,
		confFile:   confFile,
		ipList:     ipList,
		macAddress: macAddressString,
	}, nil
}

func (dhcp *Server) Run(ctx context.Context) error {
	cmds := []string{"dhcpd"}
	cid := dhcp.cont.ID()
	if err := dhcp.cont.Execute(ctx, cmds, cid); err != nil {
		log.Error().Msgf("Error in executing given DHCP command  %v", err)
	}
	return nil
}

//Stop should not be used as a command as it will close the container and there by remove the added interfaces so the container will break,
//more over the stop command is not removing the temp file from the file system
func (dhcp *Server) Stop() error {
	return dhcp.cont.Stop()
}

func (dhcp *Server) Close() error {
	if err := os.Remove(dhcp.confFile); err != nil {
		return err
	}

	if err := dhcp.cont.Close(); err != nil {
		return err
	}

	return nil
}

// might require mutex when using with goroutines
func (dhcp *Server) GetVlanIP(vlan string) string {
	return dhcp.ipList[vlan]
}

func (dhcp *Server) GetMAC() string {
	return dhcp.macAddress
}
