package game

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/aau-network-security/defatt/store"
	"github.com/aau-network-security/defatt/virtual/docker"
	"github.com/aau-network-security/defatt/virtual/vbox"
	"github.com/aau-network-security/openvswitch/ovs"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var (
	ErrVirtualInstanceNil = errors.New("failed to create virtual instance")
)

func (env *environment) initializeScenario(ctx context.Context, bridge string, scenario store.Scenario) error {
	var wg sync.WaitGroup

	for _, host := range scenario.Hosts {
		if host.Type == "docker" {
			wg.Add(1)
			go env.attachDocker(ctx, &wg, bridge, host.Image, host.Networks)
			continue
		}
		if host.Type == "vbox" {
			wg.Add(1)
			go env.attachVM(ctx, &wg, host.Name, bridge, host.Image, host.Networks)
			continue
		}
		log.Error().Msg("Unsupported challenge host")
	}
	wg.Wait()

	return nil
}

func (env *environment) attachDocker(ctx context.Context, wg *sync.WaitGroup, bridge string, image string, nets []string) error {
	defer wg.Done()

	container := docker.NewContainer(docker.ContainerConfig{
		Image: image,
		Labels: map[string]string{
			"nap-game":      bridge,
			"game-networks": strings.Join(nets, ","),
		}})

	if err := container.Create(ctx); err != nil {
		log.Error().Err(err).Msg("creating container")
		return err
	}

	if err := container.Start(ctx); err != nil {
		log.Error().Err(err).Msg("starting container")
		return err
	}

	cid := container.ID()
	if cid == "" {
		log.Error().Msg("getting ID for container")
		return ErrGettingContainerID
	}
	for i, network := range nets {
		if err := env.controller.Ovs.Docker.AddPort(bridge, fmt.Sprintf("eth%d", i), cid, ovs.DockerOptions{DHCP: true, VlanTag: network}); err != nil {
			log.Error().Err(err).Str("container", cid).Msg("adding port to container")
			return err
		}
	}

	if container == nil {
		return ErrVirtualInstanceNil
	}

	env.instances = append(env.instances, container)

	return nil
}

func (env *environment) attachVM(ctx context.Context, wg *sync.WaitGroup, name, bridge, image string, nets []string) error {
	var ifaceNames []string
	var fullIfaceName string
	defer wg.Done()
	for _, network := range nets {
		ifacesuffix := uuid.New().String()[0:5]
		ifaceName := fmt.Sprintf("%s_%s", network, ifacesuffix)
		vlan, err := strconv.Atoi(network)
		if err != nil {
			return err
		}
		if err := env.createPort(bridge, ifaceName, vlan); err != nil {
			return err
		}
		fullIfaceName = fmt.Sprintf("%s_%s_%s", bridge, network, ifacesuffix)
		ifaceNames = append(ifaceNames, fullIfaceName)
	}
	fmt.Println(name)

	if name == "mailserver" {

		vm, err := env.vlib.GetCopy(ctx,
			bridge,
			vbox.InstanceConfig{Image: "mails.ova",
				CPU:      1,
				MemoryMB: 8192},

			vbox.SetBridge(ifaceNames, true),
			vbox.SetMAC("04d30454fe15", 2),
		)

		if err != nil {
			fmt.Println(err)
			log.Err(err).Msgf("problem creating the mailserver: %v", err)
		}

		if vm == nil {
			fmt.Print("NU s-a creat masina mail \n")
		}
		env.instances = append(env.instances, vm)
		if err := vm.Start(ctx); err != nil {
			log.Error().Err(err).Msg("starting mailserver virtual machine")
			return err
		}
	} else if name == "DCcon" {
		vm, err := env.vlib.GetCopy(ctx,
			bridge,
			vbox.InstanceConfig{Image: "win10NoWDMail2.ova",
				CPU:      2,
				MemoryMB: 8192},
			vbox.SetBridge(ifaceNames, true),
			vbox.SetMAC("04d3b0c757c7", 2),
		)
		if err != nil {
			fmt.Println(err)
			log.Err(err).Msgf("problem creating the mailserver: %v", err)
		}

		if vm == nil {
			fmt.Print("NU s-a creat masina DC \n")
		}
		env.instances = append(env.instances, vm)
		if err := vm.Start(ctx); err != nil {
			log.Error().Err(err).Msg("starting mailserver virtual machine")
			return err
		}
	} else {
		vm, err := env.vlib.GetCopy(ctx,
			bridge,
			vbox.InstanceConfig{Image: image,
				CPU:      1,
				MemoryMB: 2048},
			vbox.SetBridge(ifaceNames, true),
		)

		if err != nil {
			log.Error().Err(err).Msg("VM not created ")
			return err
		}
		if vm == nil {
			return ErrVMNotCreated
		}
		env.instances = append(env.instances, vm)
		if err := vm.Start(ctx); err != nil {
			log.Error().Err(err).Msg("starting virtual machine")
			return err
		}

	}
	//if name == "DCcon" {
	//	vm, err := env.vlib.GetCopy(ctx,
	//		bridge,
	//		vbox.InstanceConfig{Image: image,
	//			CPU:      1,
	//			MemoryMB: 2048},
	//		vbox.SetBridge(ifaceNames, true),
	//		vbox.SetMAC("04d3b0c757c7", 1),
	//
	//	)
	//	log.Error().Err(err).Msg("VM not created ")
	//	return err
	//
	//	if vm == nil {
	//		return ErrVMNotCreated
	//	}
	//	env.instances = append(env.instances, vm)
	//	if err := vm.Start(ctx); err != nil {
	//		log.Error().Err(err).Msg("starting mailserver virtual machine")
	//		return err
	//	}
	//}

	return nil
}
