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
	ipList   map[string]string
}

func createDHCPFile(nets Networks) string {
	var tpl bytes.Buffer
	tmpl := template.Must(template.ParseFiles("/home/ubuntu/defat/dnet/dhcp/dhcpd.conf.tmpl"))
	tmpl.Execute(&tpl, nets)
	return tpl.String()
}
func addToSwitch(c *controller.OvsManagement, net Subnet, bridge, cid string) {
	if err := c.OvsDService.AddPort(controller.OvsDockerInfo{OvsBridge: bridge,
		Eth:       net.Interface,
		Container: cid,
		// exclusive for dhcp
		NetI: controller.NETInfo{
			IpAddr: fmt.Sprintf("%s/24", net.Router),
		}}); err != nil {
		log.Error().Msgf("Error on ovs-docker addport %v", err)
	}

	if err := c.OvsDService.SetVlan(controller.OvsDockerInfo{OvsBridge: bridge, Eth: net.Interface, Container: cid, Vlan: net.Vlan}); err != nil {
		log.Error().Msgf("Error on ovs-docker SetVlan %v", err)
	}
}

//New creates a DHCP server which will be listening on the interfaces given as the argument
func New(ctx context.Context, ifaces map[string]string, bridge string, c *controller.OvsManagement) (*Server, error) {
	ipList := make(map[string]string)
	var networks Networks
	ipPool := controller.NewIPPoolFromHost()
	for vl, vt := range ifaces {
		var sNet Subnet
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
		addToSwitch(c, net, bridge, cid)
	}

	return &Server{
		cont:     cont,
		confFile: confFile,
		ipList:   ipList,
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

//Stop should not be used as a command as it will close the container and there by remove the added interfaces so the container will break, more over the stop command is not removing the temp file from the file system
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
