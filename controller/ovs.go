package controller

import (
	"github.com/aau-network-security/openvswitch/ovs"
	"github.com/rs/zerolog/log"
)

type OvsManagement struct {
	*ovs.Client
	*ovs.VSwitchSetService
	*OvsDocker
	*NetClient
}

func (c *OvsManagement) CreateBridge(bridgeName string) error {
	if err := c.VSwitch.AddBridge(bridgeName); err != nil {
		log.Error().Msgf("Error on creating bridge on openvswitch with name %s error %v", bridgeName, err)
		return err
	}
	return nil
}

func (c *OvsManagement) RemoveBridge(bridgeName string) error {
	if err := c.VSwitch.DeleteBridge(bridgeName); err != nil {
		log.Error().Msgf("Error on deleting bridge bridge on openvswitch with name %s error %v", bridgeName, err)
		return err
	}
	return nil
}

//Here why are not also the functions for Addport SetTag, Vlan...blalblaba
