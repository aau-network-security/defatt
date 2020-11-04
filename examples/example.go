package examples

import (
	"context"
	"fmt"

	"github.com/mrturkmencom/defat/controller"
	"github.com/mrturkmencom/defat/model"
	"github.com/mrturkmencom/defat/virtual/docker"
	"github.com/mrturkmencom/defat/virtual/vbox"
	"github.com/rs/zerolog/log"
)

var (
	bridgeName = "SW"
)

type Vlan struct {
	Subnet string
	RandIP []int
}

//todo:
// Following function is PoC of adding docker containers to
// ovs-bridges, will be changed when we have dhcp server

// create a docker with none network
// start the docker container with openvswitch vlan
// guideline from IBM is followed; https://developer.ibm.com/recipes/tutorials/using-ovs-bridge-for-docker-networking/

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
	config, err := model.NewConfig("/home/ubuntu/defat/config/config.yml")
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
				//log.Error().Msgf("Failed to start virtual machine on vlan %s", main.vlans[0])
			}
		}
	}

}

func RunDocker(config docker.ContainerConfig, cli *controller.OvsManagement, ip int, vlanProp Vlan) error {
	ctx := context.Background()
	container := docker.NewContainer(config)
	if err := container.Create(ctx); err != nil {
		log.Error().Msgf("Error in creating container  %v", err)
		return err
	}
	if err := container.Start(ctx); err != nil {
		log.Error().Msgf("Error in creating container  %v", err)
		return err
	}

	cid := container.ID()
	if cid == "" {
		return fmt.Errorf("Container id could be fetched correctly")
	}

	fullIP := fmt.Sprintf("%s.%d/24", vlanProp.Subnet, ip)
	if err := cli.OvsDService.AddPort(controller.OvsDockerInfo{OvsBridge: bridgeName, Eth: "eth0", Container: cid, NetI: controller.NETInfo{IpAddr: fullIP}}); err != nil {
		log.Error().Msgf("Error on adding port on docker %v", err)
		return err
	}

	if err := cli.OvsDService.SetVlan(controller.OvsDockerInfo{OvsBridge: bridgeName, Eth: "eth0", Container: cid, Vlan: "20"}); err != nil {
		log.Error().Msgf("Error on ovs-docker SetVlan %v", err)
		return err
	}
	return nil

}
