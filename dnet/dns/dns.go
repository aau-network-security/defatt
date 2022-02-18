// Copyright (c) 2018-2019 Aalborg University
// Use of this source code is governed by a GPLv3
// license that can be found in the LICENSE file.

package dns

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aau-network-security/defatt/controller"
	"github.com/aau-network-security/defatt/store"
	"github.com/aau-network-security/openvswitch/ovs"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/aau-network-security/defatt/virtual/docker"
	"github.com/rs/zerolog/log"
)


//var (
//	//go:embed Corefile.tmpl
//	Corefile embed.FS
//
//	//go:embed zonefile.tmpl
//	zonefile embed.FS
//)







//const (

//	coreFileContent = `. {
//    file zonefile
//    prometheus     # enable metrics
//    errors         # show errors
//    log            # enable query logs
//}
//`
//	zonePrefixContent = `$ORIGIN .
//@   3600 IN SOA sns.dns.icann.org. noc.dns.icann.org. (
//                2017042745 ; serial
//                7200       ; refresh (2 hours)
//                3600       ; retry (1 hour)
//                1209600    ; expire (2 weeks)
//                3600       ; minimum (1 hour)
//                )
//
//`
//)



type Server struct {
	cont     docker.Container
	corefile string
	zonefile string
	ipList   map[string]string
}

type Domains struct{
	records []RR
	Zonefile string

}

type RR struct {
	Name  string
	Type  string
	RData string
	IPAddress string
	Domain string
}

func createCorefile(domains Domains) (string) {
	var tpl bytes.Buffer

	dir, err := os.Getwd() // get working directory
	if err != nil {
		log.Error().Msgf("Error getting the working dir for CoreFile %v", err)
	}
	fullPathToTemplate := fmt.Sprintf("%s%s", dir, "/dnet/dns/Corefile.tmpl")

	tmpl := template.Must(template.ParseFiles(fullPathToTemplate))


	tmpl.Execute(&tpl, domains)
	return tpl.String()
}

func createZonefile(datas RR) string {

	var ztpl bytes.Buffer

	dir, err := os.Getwd() // get working directory
	if err != nil {
		log.Error().Msgf("Error getting the working dir for zonefile %v", err)
	}
	fullPathToTemplate := fmt.Sprintf("%s%s", dir, "/dnet/dns/zonefile.tmpl")

	tmpl := template.Must(template.ParseFiles(fullPathToTemplate))

	//tmpl := template.Must(template.ParseFiles("/home/ubuntu/vlad/sec03/defatt/dnet/dhcp/dhcpd.conf.tmpl"))

	tmpl.Execute(&ztpl, datas)
	return ztpl.String()
}



func attachToSwitch(c *controller.NetController, contID string, bridge string, ipList map[string]string ) error{
	i:=1
	for _, network := range ipList {
		if network == "10.10.10.0/24" {
			continue
		} else {
			ipAddrs := strings.TrimSuffix(network, ".0/24")
			ipAddrs = ipAddrs + ".3/24"

			fmt.Println(ipAddrs)
			//fmt.Sprintf("eth%d", vlan)
			tag := i * 10

			sTag := strconv.Itoa(tag)

			fmt.Println(sTag)
			if err := c.Ovs.Docker.AddPort(bridge, fmt.Sprintf("eth%d", i), contID, ovs.DockerOptions{VlanTag: sTag, IPAddress: ipAddrs}); err != nil {

				log.Error().Err(err).Str("container", contID).Msg("adding port to DNS container")
				return err
			}
			i++
			fmt.Println(i)

		}
	}

	return nil

}



func New(control *controller.NetController, bridge string, ipList map[string]string, scenario store.Scenario) (*Server, error) {


	var domains Domains
	var records RR

	records.Name = scenario.DNS
	stripTLD := strings.SplitAfter(scenario.DNS, ".")
	domains.Zonefile = stripTLD[0]
	c, err := ioutil.TempFile("", "Corefile")
	if err != nil {
		return nil, err
	}

	Corefile := c.Name()

	CorefileStr := createCorefile(domains)

	_,err = c.WriteString(CorefileStr)
	if err != nil{
		return nil, err
	}

	for _, network := range ipList{
		records.IPAddress = network
		break
	}

	records.Domain = scenario.DNS

	z, err := ioutil.TempFile("", "zonefile")
	if err != nil {
		return nil, err
	}

	zonefile := z.Name()


	zonefileStr :=createZonefile(records)

	_,err = c.WriteString(zonefileStr)
	if err != nil{
		return nil, err
	}

	cont := docker.NewContainer(docker.ContainerConfig{
		Image: "coredns/coredns:1.8.6",
		Mounts: []string{
			fmt.Sprintf("%s:/Corefile", Corefile),
			fmt.Sprintf("%s:/zonefile", zonefile),
		},
		UsedPorts: []string{
			"53/tcp",
			"53/udp",
		},
		Resources: &docker.Resources{
			MemoryMB: 50,
			CPU:      0.3,
		},
		Cmd: []string{"--conf", "Corefile"},
		Labels: map[string]string{
			"nap-game":      bridge,
		},
	})
	contID := cont.ID()

	if err := attachToSwitch(control,contID, bridge, ipList); err != nil {
		log.Error().Msgf("Error on addToSwitch in DNS %v ", err)
	}

	return &Server{
		cont:     cont,
		corefile: Corefile,
		zonefile: zonefile,
	}, nil
}

func (s *Server) Container() docker.Container {
	return s.cont
}

func (s *Server) Run(ctx context.Context) error {
	return s.cont.Run(ctx)
}

func (s *Server) Close() error {
	if err := os.Remove(s.corefile); err != nil {
		log.Warn().Msgf("error while removing DNS configuration file: %s", err)
	}

	if err := s.cont.Close(); err != nil {
		log.Warn().Msgf("error while closing DNS container: %s", err)
	}

	return nil
}

func (s *Server) Stop() error {
	return s.cont.Stop()
}
