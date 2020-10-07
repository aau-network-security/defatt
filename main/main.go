package main

import (
	"context"
	"fmt"

	"github.com/aau-network-security/openvswitch/ovs"
	"github.com/mrturkmencom/defat/controller"
	"github.com/mrturkmencom/defat/model"
	"github.com/mrturkmencom/defat/virtual/docker"
	"github.com/mrturkmencom/defat/virtual/vbox"
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

	tapTags := map[string]string{
		"tap0": "10",
		"tap2": "20",
		"tap4": "30",
	}

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
	if err := c.IFConfig.OvsBridgeUp(bridgeName, "192.168.0.1", "255.255.0.0"); err != nil {
		log.Error().Msgf("Error %v", err)
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
		if err := c.VSwitch.AddPortTagged(bridgeName, t, tag); err != nil {
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

	//randomized ips could be changed according to
	//upcoming requests and requirements
	ipPool := controller.NewIPPoolFromHost()
	//example of random ip addresses
	for i := 0; i < 10; i++ {
		randomIp, _ := ipPool.Get()
		fmt.Printf("Auto generated random ip is: %s.0/24\n", randomIp)
	}

	// create a docker with none network
	// start the docker container with openvswitch vlan
	// guideline from IBM is followed; https://developer.ibm.com/recipes/tutorials/using-ovs-bridge-for-docker-networking/
	addDockerToOVS(c)
}

//todo:
// Following function is PoC of adding docker containers to
// ovs-bridges, will be changed when we have dhcp server
func addDockerToOVS(c *controller.OvsManagement) {
	container := docker.NewContainer(docker.ContainerConfig{
		Image: "mrturkmen/joomla-rce",
		EnvVars: map[string]string{
			"APP_FLAG": "Testing app flag",
		},
		UseBridge: false,
	})
	if err := container.Create(context.Background()); err != nil {
		log.Error().Msgf("Error in creating container  %v", err)
	}

	if err := container.Start(context.Background()); err != nil {
		log.Error().Msgf("Error in creating container  %v", err)
	}
	cid := container.ID()
	if cid == "" {
		panic(fmt.Errorf("ERROR DOCKER CONTAINER DOES NOT HAVE ID FOR ITSELFS"))
	}
	// attach ovs-docker
	if err := c.OvsDService.AddPort(controller.OvsDockerInfo{OvsBridge: bridgeName, Eth: "eth0", Container: cid,
		NetI: controller.NETInfo{
			IpAddr:  "192.168.1.1/16",
			Gateway: "192.168.0.1",
		}}); err != nil {
		log.Error().Msgf("Error on ovs-docker addport %v", err)
	}

	container_x := docker.NewContainer(docker.ContainerConfig{
		Image:     "guacamole/guacd:1.0.0",
		UseBridge: false,
	})

	if err := container_x.Create(context.Background()); err != nil {
		log.Error().Msgf("Error in creating container  %v", err)
	}

	if err := container_x.Start(context.Background()); err != nil {
		log.Error().Msgf("Error in starting container  %v", err)
	}

	cid_x := container_x.ID()
	if cid_x == "" {
		panic(fmt.Errorf("ERROR DOCKER CONTAINER DOES NOT HAVE ID FOR ITSELFS"))
	}

	if err := c.OvsDService.AddPort(controller.OvsDockerInfo{OvsBridge: bridgeName, Eth: "eth0", Container: cid_x,
		NetI: controller.NETInfo{
			IpAddr:  "192.168.1.2/16",
			Gateway: "192.168.0.1",
		}}); err != nil {
		log.Error().Msgf("Error on ovs-docker addport %v", err)
	}

}

// todo:
// Following function is proof of concept which shows,
// functions are working nicely, it will be changed
// when dhcp server is ready.
func addVMsToOVS(c *controller.OvsManagement) {
	var vms map[string][]string
	//parse configuration file
	config, err := model.NewConfig("<path-to-config>")
	if err != nil {
		log.Error().Msgf("Error on reading configuration file %v", err)
	}
	log.Debug().Msgf("NewConfig read from given place...")
	// import and run a vm on an openvswitch vlan
	log.Debug().Msgf("Creating vbox library for vbox...")
	vlib := vbox.NewLibrary(config.VmConfig.OvaDir)
	if vlib == nil {
		log.Error().Msgf("Library could not be created properly...")
	}
	//map structure will have ids of vms and attached vlans to those vlans
	//in each vm, we are enabling promiscuous  mode
	vms = map[string][]string{
		"vm1": {"vlan10", "vlan30"},
		"vm2": {"vlan20"},
		"vm3": {"vlan30"},
	}
	for _, vl := range vms {
		log.Info().Msgf("VL content is : %v", vl)
		vm, err := vlib.GetCopy(context.Background(),
			vbox.InstanceConfig{Image: "kali.ova",
				CPU:      2,
				MemoryMB: 4096},
			vbox.SetBridge(vl),
		)
		if err != nil {
			log.Error().Msgf("Error while getting copy of VM")
		}
		if vm != nil {
			log.Debug().Msgf("VM %s has following vlans attached %v ", vm.Info().Id, vl)
			vms[vm.Info().Id] = vl
			log.Debug().Msgf("VM [ %s ] is starting .... ", vm.Info().Id)
			if err := vm.Start(context.Background()); err != nil {
				log.Error().Msgf("Failed to start virtual machine on vlan %s", vlans[0])
			}
		}
	}

}
