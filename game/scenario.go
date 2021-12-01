package game

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/aau-network-security/defatt/store"
	"github.com/aau-network-security/defatt/virtual/docker"
	"github.com/aau-network-security/defatt/virtual/vbox"
	"github.com/aau-network-security/openvswitch/ovs"
	"github.com/rs/zerolog/log"
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
			"nap": "challenges",
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

	return nil
}

func (env *environment) attachVM(ctx context.Context, wg *sync.WaitGroup, name, bridge, image string, nets []string) error {
	var ifaceNames []string
	defer wg.Done()
	for _, network := range nets {
		ifaceName := fmt.Sprintf("%s_%s", network, name[0:5])
		vlan, err := strconv.Atoi(network)
		if err != nil {
			return err
		}
		if err := env.createPort(bridge, ifaceName, vlan); err != nil {
			return err
		}
		fullIfaceName := fmt.Sprintf("%s_%s_%s", bridge, network, name[0:5])
		ifaceNames = append(ifaceNames, fullIfaceName)
	}

	vm, err := env.vlib.GetCopy(ctx,
		bridge,
		vbox.InstanceConfig{Image: image,
			CPU:      1,
			MemoryMB: 2048},
		vbox.SetBridge(ifaceNames, true),
	)
	if err != nil {
		return err
	}
	if vm == nil {
		return ErrVMNotCreated
	}
	if err := vm.Start(ctx); err != nil {
		log.Error().Err(err).Msgf("starting virtual machine")
		return err
	}

	return nil
}