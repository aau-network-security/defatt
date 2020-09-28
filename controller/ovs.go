package controller

import (
	"github.com/aau-network-security/openvswitch/ovs"
	"github.com/rs/zerolog/log"
)

type OvsManagement struct {
	*ovs.Client
	*ovs.VSwitchSetService
	*NetClient
}

func (c *OvsManagement) CreateBridge(bridgeName string) error {
	if err := c.VSwitch.AddBridge(bridgeName); err != nil {
		log.Error().Msgf("Error on creating bridge on openvswitch with name %s error %v", bridgeName, err)
		return err
	}
	return nil
}

func (c *OvsManagement) addPort(bridgeName, port string) error {

	if err := c.VSwitch.AddPort(bridgeName, port); err != nil {
		log.Error().Msgf("Error on adding port on openvswitch %v", err)
		return err
	}
	return nil
}
