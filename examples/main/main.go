package main

import (
	"context"
	"fmt"

	wg "github.com/aau-network-security/defat/app/daemon/vpn-proto"
	"github.com/aau-network-security/defat/controller"
	vpn "github.com/aau-network-security/defat/dnet/wg"
	"github.com/aau-network-security/defat/examples"
	"github.com/aau-network-security/defat/virtual/docker"
	"github.com/aau-network-security/openvswitch/ovs"
	"github.com/rs/zerolog/log"
)

// Following function calls are equivalent to the code in bash script
// https://github.com/aau-network-security/NAP/blob/master/setup-openswitch.sh

// This is just a PoC

var (
	vlans      = []string{"vlan10", "vlan20", "vlan30"}
	taps       = []string{"tap0", "tap2", "tap4"}
	bridgeName = "SW"
)

func main() {

	wgClient, err := vpn.NewGRPCVPNClient(vpn.WireGuardConfig{
		// Fill out the required fields...
		Endpoint: "",
		Port:     0,
		AuthKey:  "",
		SignKey:  "",
		Enabled:  false,
		CertFile: "",
		CertKey:  "",
		CAFile:   "",
		Dir:      "",
	})
	if err != nil {
		log.Printf("Error VPN connection could not be initialized ! ")
	}
	vlanTags := map[string]string{
		"vlan10": "10",
		"vlan20": "20",
		"vlan30": "30",
	}

	tapTags := map[string]string{
		"tap0": "10",
		"tap2": "20",
		"tap4": "30",
	}

	c := controller.New(controller.Sudo())
	//ovs-vsctl add-br SW
	if err := c.Ovs.VSwitch.AddBridge(bridgeName); err != nil {
		log.Error().Msgf("Error on creating OVS bridge %v", err)
	}
	// if err := c.IFConfig.OvsBridgeUp(bridgeName, "192.168.0.1", "255.255.0.0"); err != nil {
	// 	log.Error().Msgf("Error %v", err)
	// }

	for vl, vt := range vlanTags {
		//ovs-vsctl add-port SW vlan10 tag=10 -- set interface vlan10 type=internal
		//ovs-vsctl add-port SW vlan20 tag=20 -- set interface vlan20 type=internal
		//ovs-vsctl add-port SW vlan30 tag=30 -- set interface vlan30 type=internal
		if err := c.Ovs.VSwitch.AddPortTagged(bridgeName, vl, vt); err != nil {
			log.Error().Msgf("Error on adding port with tag err %v", err)
		}
		log.Info().Msgf("AddPort Set Interface Options %s", vl)
		if err := c.Ovs.VSwitch.Set.Interface(vl, ovs.InterfaceOptions{Type: ovs.InterfaceTypeInternal}); err != nil {
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
		if err := c.Ovs.VSwitch.AddPortTagged(bridgeName, t, tag); err != nil {
			log.Error().Msgf("Error on adding port with tag err %v", err)
		}
	}

	ipPool := controller.NewIPPoolFromHost()

	for _, v := range vlans {
		//ifconfig vlan10 up
		//ifconfig vlan20 up
		//ifconfig vlan30 up
		if err := c.IFConfig.TapUp(v); err != nil {
			log.Error().Msgf("Error happened on making up tap %s %v", v, err)
		}

	}

	log.Info().Msgf("Taps are created and upped")

	interfaces, err := c.Ovs.VSwitch.ListBridges()
	if err != nil {
		log.Error().Msgf("Error on listing bridge %v", err)
	}
	for _, i := range interfaces {
		fmt.Printf("Created interface:  %s\n", i)
	}
	//
	//server, err := dhcp.New(context.Background(), vlanTags, bridgeName, c)
	//if err != nil {
	//	log.Error().Msgf("Error creating DHCP server %v", err)
	//}
	//if err := server.Run(context.Background()); err != nil {
	//	log.Error().Msgf("Error in starting DHCP  %v", err)
	//}

	dockerContainers := make(map[string]docker.ContainerConfig)
	dockerContainers["joomla"] = docker.ContainerConfig{
		Image: "mrturkmen/joomla-rce",
		EnvVars: map[string]string{
			"APP_FLAG": "Testing app flag",
		},
		UseBridge: false,
	}

	dockerContainers["httpd"] = docker.ContainerConfig{Image: "httpd", UseBridge: false}

	log.Debug().Msgf("AddVMSToOvs is running....")
	go func() {
		if err := examples.AddVMsToOvs(); err != nil {
			log.Error().Msgf("Error on vm returned %v", err)
		}
	}()

	for i, config := range dockerContainers {
		fmt.Printf("Run command for container %s\n", i)
		log.Info().Msgf("Executing commands for container %s", i)
		if err := examples.RunDocker(config, c, "20"); err != nil {
			fmt.Printf("Error returned %v", err)
		}
	}

	////for _, v := range vlans {
	//	fmt.Printf("Running initiallize I .... \n")
	randomizedIP, _ := ipPool.Get()
	//port, err := strconv.Atoi(v[len(v)-2:])
	//if err != nil {
	//	fmt.Printf("Error convert string to int: %v", err)
	//}
	listenport := 4020
	resp, err := wgClient.InitializeI(context.TODO(), &wg.IReq{
		Address:    fmt.Sprintf("%s.1", randomizedIP),
		ListenPort: uint32(listenport),
		SaveConfig: false,
		Eth:        "eth0",
		IName:      "vlan20" + "_vpn",
	})
	if err != nil {
		fmt.Printf("ERROR ON INITIALIZING VPN ENDPOINT %d , ERR: %v\n", listenport, err)
	}
	if resp != nil {
		fmt.Printf("Message from wg service %s\n", resp.Message)
	}
	//}

}
