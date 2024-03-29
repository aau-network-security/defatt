package examples

import (
	"context"
	"fmt"

	"github.com/aau-network-security/defat/config"
	"github.com/aau-network-security/defat/controller"
	"github.com/aau-network-security/defat/virtual"
	"github.com/aau-network-security/defat/virtual/docker"
	"github.com/aau-network-security/defat/virtual/vbox"
	"github.com/aau-network-security/openvswitch/ovs"
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

func addDockerToOVS(c *controller.NetController, vlan string) error {

	container := docker.NewContainer(docker.ContainerConfig{
		Image: "mrturkmen/joomla-rce",
		EnvVars: map[string]string{
			"APP_FLAG": "Testing app flag",
		},
		UseBridge: false,
	})
	if err := container.Create(context.Background()); err != nil {
		log.Error().Msgf("Error in creating container  %v", err)
		return err
	}

	if err := container.Start(context.Background()); err != nil {
		log.Error().Msgf("Error in creating container  %v", err)
		return err
	}
	cid := container.ID()
	if cid == "" {
		panic(fmt.Errorf("ERROR DOCKER CONTAINER DOES NOT HAVE ID FOR ITSELFS"))
	}
	// attach ovs-docker
	if err := c.Ovs.Docker.AddPort(bridgeName, "eth0", cid, ovs.DockerOptions{
		DHCP:    true,
		VlanTag: vlan,
	}); err != nil {
		log.Error().Msgf("Error on ovs-docker addport %v", err)
		return err
	}

	container_x := docker.NewContainer(docker.ContainerConfig{
		Image:     "guacamole/guacd:1.0.0",
		UseBridge: false,
	})

	if err := container_x.Create(context.Background()); err != nil {
		log.Error().Msgf("Error in creating container  %v", err)
		return err
	}

	if err := container_x.Start(context.Background()); err != nil {
		log.Error().Msgf("Error in starting container  %v", err)
		return err
	}

	cid_x := container_x.ID()
	if cid_x == "" {
		panic(fmt.Errorf("ERROR DOCKER CONTAINER DOES NOT HAVE ID FOR ITSELFS"))
	}

	if err := c.Ovs.Docker.AddPort(bridgeName, "eth0", cid_x, ovs.DockerOptions{
		DHCP:    true,
		VlanTag: vlan,
	}); err != nil {
		log.Error().Msgf("Error on ovs-docker addport %v", err)
		return err
	}
	return nil
}

// todo:
// Following function is proof of concept which shows,
// functions are working nicely, it will be changed
// when dhcp server is ready.
func AddVMsToOvs() error {
	//var vms map[string][]string
	//parse configuration file
	conf, err := config.NewConfig("/Users/rvm/Downloads/AAUJOB/NAP/2021/defatt/config/config.yml")
	if err != nil {
		log.Error().Msgf("Error on reading configuration file %v", err)
		return err
	}
	log.Debug().Msgf("NewConfig read from given place...")
	// import and run a vm on an openvswitch vlan
	log.Debug().Msgf("Creating vbox library for vbox...")
	vlib := vbox.NewLibrary(conf.VmConfig.OvaDir)
	if vlib == nil {
		log.Error().Msgf("Library could not be created properly...")
		return fmt.Errorf("Error on new library")
	}
	//map structure will have ids of vms and attached vlans to those vlans
	//in each vm, we are enabling promiscuous  mode

	networks := []string{"vlan10", "vlan20", "vlan30"}

	log.Info().Msgf("VL content is : %v", networks)
	// if the wireguard vm is connected to all vlans it means that it is endpoint for  blue teams
	// if it is connected to only specific vlans it is for red teams
	vm, err := vlib.GetCopy(context.Background(),
		vbox.InstanceConfig{Image: "haaukins.ova",
			CPU:      1,
			MemoryMB: 256},
		vbox.MapVMPort([]virtual.NatPortSettings{
			{
				// this is for gRPC service
				HostPort:    "5353",
				GuestPort:   "5353",
				ServiceName: "wgservice",
				Protocol:    "tcp",
			},
			{
				// this is for VPN Connection
				HostPort:    "51820",
				GuestPort:   "51820",
				ServiceName: "wireguard",
				Protocol:    "udp",
			},
		}),
		// SetBridge parameter cleanFirst should be enabled when wireguard/router instance
		// is attaching to openvswitch network
		vbox.SetBridge(networks, false),

	)

	if err != nil {
		log.Error().Msgf("Error while getting copy of VM")
		return err
	}
	if vm != nil {
		log.Debug().Msgf("VM [ %s ] is starting .... ", vm.Info().Id)
		if err := vm.Start(context.Background()); err != nil {
			log.Error().Msgf("Failed to start virtual machine on vlan ")
			return err
		}
	}
	return nil
}

func RunDocker(config docker.ContainerConfig, cli *controller.NetController, vlan string) error {
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

	if err := cli.Ovs.Docker.AddPort(bridgeName, "eth0", cid, ovs.DockerOptions{DHCP: true, VlanTag: vlan}); err != nil {
		log.Error().Msgf("Error on adding port on docker %v", err)
		return err
	}

	if err := cli.Ovs.Docker.SetVlan(bridgeName, "eth0", cid, vlan); err != nil {
		log.Error().Msgf("Error on ovs-docker SetVlan %v", err)
		return err
	}
	return nil

}
