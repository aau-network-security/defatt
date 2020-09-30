package main

import (
	"fmt"

	"github.com/aau-network-security/openvswitch/ovs"
	"github.com/mrturkmencom/defat/controller"
	"github.com/rs/zerolog/log"
)

// Following function calls are equivalent to the code in bash script
// https://github.com/aau-network-security/NAP/blob/master/setup-openswitch.sh

// This is just a PoC

func main() {

	vlans := []string{"vlan10", "vlan20", "vlan30"}
	taps := []string{"tap0", "tap2", "tap4"}
	tapTags := map[string]string{
		"tap0": "10",
		"tap2": "20",
		"tap4": "30",
	}
	bridgeName := "SW"
	c := &controller.OvsManagement{
		Client: ovs.New(
			ovs.Sudo(),
			ovs.Debug(false),
		),
		NetClient: controller.New(controller.Sudo()),
	}

	//ovs-vsctl add-br SW
	if err := c.CreateBridge(bridgeName); err != nil {
		log.Error().Msgf("Error on creating OVS bridge %v", err)
	}

	for _, vl := range vlans {
		//ovs-vsctl add-port SW vlan10 tag=10 -- set interface vlan10 type=internal
		//ovs-vsctl add-port SW vlan20 tag=20 -- set interface vlan20 type=internal
		//ovs-vsctl add-port SW vlan30 tag=30 -- set interface vlan30 type=internal
		if err := c.VSwitch.AddPort(bridgeName, vl); err != nil {
			log.Error().Msgf("Error on adding port with tag err %v", err)
		}
		log.Info().Msgf("AddPort Set Interface Options %s", vl)
		if err := c.VSwitch.Set.Interface(vl, ovs.InterfaceOptions{Type: ovs.InterfaceTypeInternal}); err != nil {
			log.Error().Msgf("Error on matching interface error %v", err)
		}

	}

	for _, t := range taps {
		//ip tuntap add tap0 mode tap
		//ifconfig tap0 up
		//ip tuntap add tap2 mode tap
		// ifconfig tap2 up
		//ip tuntap add tap4 mode tap
		//ifconfig tap4 up
		if err := c.IPService.AddTunTap(t, "tap"); err != nil {
			log.Error().Msgf("Error happened on adding tuntap %v", err)
		}
		if err := c.IFConfig.TapUp(t); err != nil {
			log.Error().Msgf("Error happened on making up tap %s %v", t, err)
		}
	}

	for t, tag := range tapTags {
		//ovs-vsctl add-port SW tap0 tag=10
		//ovs-vsctl add-port SW tap2 tag=20
		//ovs-vsctl add-port SW tap4 tag=30
		if err := c.VSwitch.AddPortTagged(bridgeName, t, fmt.Sprintf("tag=%s", tag)); err != nil {
			log.Error().Msgf("Error on adding port with tag err %v", err)
		}
	}

	for _, v := range vlans {
		//ifconfig vlan10 up
		//ifconfig vlan20 up
		//ifconfig vlan30 up
		if err := c.IFConfig.TapUp(v); err != nil {
			log.Error().Msgf("Error happened on making up tap %s %v", v, err)
		}
	}

	log.Info().Msgf("Taps are created and upped")

	interfaces, err := c.VSwitch.ListBridges()
	if err != nil {
		log.Error().Msgf("Error on listing bridge %v", err)
	}
	for _, i := range interfaces {
		fmt.Printf("Created interface:  %s\n", i)
	}

	// randomized ips could be changed according to
	// upcoming requests and requirements
	ipPool := controller.NewIPPoolFromHost()
	// example of random ip addresses
	for i := 0; i < 10; i++ {
		randomIp, _ := ipPool.Get()
		fmt.Printf("Auto generated random ip is: %s.0/24\n", randomIp)
	}

}
