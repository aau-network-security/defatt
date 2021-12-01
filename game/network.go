package game

import (
	"fmt"

	"github.com/aau-network-security/defatt/models"
	"github.com/rs/zerolog/log"
)

func (env *environment) initializeOVSBridge(bridgeName string) error {
	if err := env.controller.Ovs.VSwitch.AddBridge(bridgeName); err != nil {
		log.Error().Err(err).Msg("creating OVS bridge")
		return err
	}
	return nil
}

func (env *environment) createNetworks(bridge string, nets []models.Network) error {
	for _, network := range nets {
		if err := env.createPort(bridge, network.Name, network.Tag); err != nil {
			return err
		}
	}

	return nil
}

func (env *environment) createPort(bridge string, input string, vlan int) error {
	name := fmt.Sprintf("%s_%s", bridge, input)
	if err := env.controller.IPService.AddTunTap(name, "tap"); err != nil {
		log.Error().Err(err).Str("port", name).Msg("creating interface")
		return err
	}

	if err := env.controller.IFConfig.TapUp(name); err != nil {
		log.Error().Err(err).Str("port", name).Msg("setting interface UP")
		return err
	}

	if err := env.controller.Ovs.VSwitch.AddPortTagged(bridge, name, fmt.Sprint(vlan)); err != nil {
		log.Error().Err(err).Str("port", name).Msgf("adding port to switch")
		return err
	}

	return nil
}
