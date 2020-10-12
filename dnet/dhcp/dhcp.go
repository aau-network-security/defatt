package dhcp

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/mrturkmencom/defat/controller"
	"github.com/mrturkmencom/defat/virtual/docker"
	"github.com/rs/zerolog/log"
)

// Could be put inside some envirionment struct moving on

type Networks struct {
	Subnets []Subnet
	DNS     string
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
	cont     docker.Container
	confFile string
}

func createDHCPFile(nets Networks) string {
	var tpl bytes.Buffer
	tmpl := template.Must(template.ParseFiles("dhcpd.conf.tmpl"))
	tmpl.Execute(&tpl, nets)
	return tpl.String()
}
func addToSwitch(c *controller.OvsManagement, net Subnet, bridge, cid string) {
	if err := c.OvsDService.AddPort(controller.OvsDockerInfo{OvsBridge: bridge, Eth: net.Interface, Container: cid,
		NetI: controller.NETInfo{
			IpAddr: fmt.Sprintf("%s/24", net.Router),
		}}); err != nil {
		log.Error().Msgf("Error on ovs-docker addport %v", err)
	}

	if err := c.OvsDService.SetVlan(controller.OvsDockerInfo{OvsBridge: bridge, Eth: net.Interface, Container: cid, Vlan: net.Vlan}); err != nil {
		log.Error().Msgf("Error on ovs-docker addport %v", err)
	}
}

//New creates a DHCP server which will be listening on the interfaces given as the argument
func New(ifaces map[string]string, bridge string, c *controller.OvsManagement) (*Server, error) {
	var networks Networks
	ipPool := controller.NewIPPoolFromHost()
	for vl, vt := range ifaces {
		var net Subnet
		randIP, _ := ipPool.Get()
		net.Interface = vl
		net.Vlan = vt
		net.Network = randIP + ".0"
		net.Min = randIP + ".6"
		net.Max = randIP + ".254"
		net.Router = randIP + ".1"
		networks.Subnets = append(networks.Subnets, net)
	}
	f, err := ioutil.TempFile("", "dhcpd-conf")
	if err != nil {
		return nil, err
	}
	confFile := f.Name()

	confStr := createDHCPFile(networks)
	_, err = f.WriteString(confStr)
	if err != nil {
		return nil, err
	}
	cont := docker.NewContainer(docker.ContainerConfig{
		Image: "lanestolen/dhcp", // no need to add tag since it is not updated for 5 months.
		Mounts: []string{
			fmt.Sprintf("%s:/etc/dhcp/dhcpd.conf", confFile),
		},
		UsedPorts: []string{"67/udp"},
		Resources: &docker.Resources{
			MemoryMB: 50,
			CPU:      0.3,
		},
		Labels: map[string]string{
			"nap": "lab_dhcpd",
		},
		UseBridge: false,
	})
	if err := cont.Create(context.Background()); err != nil {
		log.Error().Msgf("Error in creating container  %v", err)
	}
	if err := cont.Start(context.Background()); err != nil {
		log.Error().Msgf("Error in starting container  %v", err)
	}
	cid := cont.ID()
	for _, net := range networks.Subnets {
		addToSwitch(c, net, bridge, cid)
	}

	return &Server{
		cont:     cont,
		confFile: confFile,
	}, nil
}

func (dhcp *Server) Run() error {
	runCMD := []string{"dhcpd"}
	cid := dhcp.cont.ID()
	if err := dhcp.cont.Execute(runCMD, cid); err != nil {
		log.Error().Msgf("Error in starting DHCP  %v", err)
	}
	return nil
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
