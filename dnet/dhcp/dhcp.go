// Copyright (c) 2018-2019 Aalborg University
// Use of this source code is governed by a GPLv3
// license that can be found in the LICENSE file.

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
)

// Could be put inside some envirionment struct moving on

type Networks struct {
	Subnets []Subnet
}

type Subnet struct {
	Network string
	Min     string
	Max     string
	Router  string
}

func createDHCPFile(amount int) string {
	var networks Networks
	var tpl bytes.Buffer
	ipPool := controller.NewIPPoolFromHost()
	for i := 0; i < amount; i++ {
		var net Subnet
		randIP, _ := ipPool.Get()
		net.Network = randIP + ".0"
		net.Min = randIP + ".6"
		net.Max = randIP + ".254"
		net.Router = randIP + ".1"
		networks.Subnets = append(networks.Subnets, net)
	}
	tmpl := template.Must(template.ParseFiles("dhcpd.conf.tmpl"))
	tmpl.Execute(&tpl, networks)
	return tpl.String()
}

type Server struct {
	cont     docker.Container
	confFile string
}

func New(format func(n int) string) (*Server, error) {
	f, err := ioutil.TempFile("", "dhcpd-conf")
	if err != nil {
		return nil, err
	}
	confFile := f.Name()

	confStr := createDHCPFile(3)
	_, err = f.WriteString(confStr)
	if err != nil {
		return nil, err
	}
	cont := docker.NewContainer(docker.ContainerConfig{
		Image: "networkboot/dhcpd", // no need to add tag since it is not updated for 5 months.
		Mounts: []string{
			fmt.Sprintf("%s:/data/dhcpd.conf", confFile),
		},
		UsedPorts: []string{"67/udp"},
		Resources: &docker.Resources{
			MemoryMB: 50,
			CPU:      0.3,
		},
		Cmd: []string{""},
		Labels: map[string]string{
			"nap": "lab_dhcpd",
		},
	})

	return &Server{
		cont:     cont,
		confFile: confFile,
	}, nil
}

func (dhcp *Server) Container() docker.Container {
	return dhcp.cont
}

func (dhcp *Server) Run(ctx context.Context) error {
	return dhcp.cont.Run(ctx)
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

func (dhcp *Server) Stop() error {
	return dhcp.cont.Stop()
}
