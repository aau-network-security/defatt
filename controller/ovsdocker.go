package controller

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

// precondition : to have ovs-docker in your path as executable

type OvsDocker struct {
	c *NetClient
}

type OvsDockerInfo struct {
	OvsBridge string
	Eth       string
	Container string // either container name if provided or container id
	Vlan      string
	NetI      NETInfo
}

type NETInfo struct {
	IpAddr string // in CIDR/4 format ; 10.0.9.1/24
	// additional fields when providing macaddress
	MACAddress string // todo: optional; could be discussed later stages
	Gateway    string // todo: optional; could be discussed later stages
	MTU        int    // todo: optional; could be discussedd later stages
}

type OvsDockerOpts func([]string, *OvsDocker) error

// exec executes an ExecFunc using 'ip'.

func (ovsd *OvsDocker) exec(args ...string) ([]byte, error) {
	return ovsd.c.exec("ovs-docker", args...)
}

// AddPort adds INTERFACE inside CONTAINER and connects it as a port
//in Open vSwitch BRIDGE. Optionally, sets ADDRESS on
//INTERFACE. ADDRESS can include a '/' to represent network
//prefix length

func (ovsd *OvsDocker) AddPort(dInfo OvsDockerInfo) error {
	//ovs-docker add-port br-int eth1 c474a0e2830e
	//--ipaddress=192.168.1.2/24 --gateway=192.168.1.1
	//--macaddress="a2:c3:0d:49:7f:f8" --mtu=1450
	log.Debug().Msgf("Adding port %s to container %s ...", dInfo.OvsBridge, dInfo.Container)
	cmds := []string{"add-port", dInfo.OvsBridge, dInfo.Eth, dInfo.Container, fmt.Sprintf("--ipaddress=%s", dInfo.NetI.IpAddr)}
	_, err := ovsd.c.OvsDService.exec(cmds...)
	return err
}

// DelPort Deletes INTERFACE inside CONTAINER and
// removes its connection to Open vSwitch BRIDGE
func (ovsd *OvsDocker) DelPort(dInfo OvsDockerInfo) error {
	// ovs-docker del-port br-int eth1 c474a0e2830e
	log.Debug().Msgf("Deleting port %s from container %s", dInfo.OvsBridge, dInfo.Container)
	cmds := []string{"del-port", dInfo.OvsBridge, dInfo.Eth, dInfo.Container}
	_, err := ovsd.c.OvsDService.exec(cmds...)
	return err
}

// DelPorts removes all Open vSwitch interfaces from CONTAINER

func (ovsd *OvsDocker) DelPorts(dInfo OvsDockerInfo) error {
	// ovs-docker del-ports br-int c474a0e2830es
	log.Debug().Msgf("Removing all Open vSwitch interfaces from %s", dInfo.Container)
	cmds := []string{"del-ports", dInfo.OvsBridge, dInfo.Container}
	_, err := ovsd.c.OvsDService.exec(cmds...)
	return err
}

// SetVlan configures the INTERFACE of CONTAINER attached to BRIDGE to become an access port of VLAN
func (ovsd *OvsDocker) SetVlan(dInfo OvsDockerInfo) error {
	// ovs-docker set-vlan br-int eth1 c474a0e2830e 5
	log.Debug().Msgf("Configures the INTERFACE of CONTAINER attached to BRIDGE to become an access port of VLAN")
	cmds := []string{"set-vlan", dInfo.OvsBridge, dInfo.Eth, dInfo.Container, dInfo.Vlan}
	_, err := ovsd.c.OvsDService.exec(cmds...)
	return err
}
