// Copyright (c) 2018-2019 Aalborg University
// Use of this source code is governed by a GPLv3
// license that can be found in the LICENSE file.

package dns

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aau-network-security/defatt/store"
	"io/ioutil"
	"os"
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

type Domains struct {
	Zonefile string
	URL      string
}

type RR struct {
	Type      string
	RData     string
	IPAddress string
	Domain    string
}

func createCorefile(domains Domains) string {
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

//func AttachToSwitch(c *controller.NetController, contID string, bridge string, ipList map[string]string ) error{
//
//	return nil
//
//}

func New(bridge string, ipList map[string]string, scenario store.Scenario) (*Server, error) {

	var domains Domains
	var records RR

	domains.URL = scenario.FQDN
	//stripTLD := strings.SplitAfter(scenario.FQDN, ".")
	stripTLD := strings.Split(scenario.FQDN, ".")
	fmt.Printf("aici trebuie sa fie doar domaniul: %s\n", stripTLD[0])
	domains.Zonefile = stripTLD[0]

	c, err := ioutil.TempFile("", "Corefile")
	if err != nil {
		return nil, err
	}

	Corefile := c.Name()
	fmt.Printf("Asta este numele coreFile: %s\n", Corefile)
	CorefileStr := createCorefile(domains)

	_, err = c.WriteString(CorefileStr)
	if err != nil {
		return nil, err
	}

	for _, network := range ipList {

		ipAddrs := strings.TrimSuffix(network, ".0/24")
		ipAddrs = ipAddrs + ".2"
		records.IPAddress = ipAddrs

		break
	}

	records.Domain = scenario.FQDN

	z, err := ioutil.TempFile("", "zonefile")
	if err != nil {
		return nil, err
	}

	zonefile := z.Name()
	fmt.Printf("Asta este numele zonefile: %s\n", zonefile)

	zonefileStr := createZonefile(records)

	_, err = z.WriteString(zonefileStr)
	if err != nil {
		return nil, err
	}

	cont := docker.NewContainer(docker.ContainerConfig{
		Image: "coredns/coredns:latest",
		Mounts: []string{
			fmt.Sprintf("%s:/Corefile", Corefile),
			fmt.Sprintf("%s:/root/db.%s", zonefile, domains.Zonefile),
			fmt.Sprintf("%s:/root/db.blue.monitor", "db.blue.monitor"),
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
			"nap-game": bridge,
		},
	})

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
