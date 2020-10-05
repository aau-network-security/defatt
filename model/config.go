package model

import (
	"fmt"
	"io/ioutil"

	dockerclient "github.com/fsouza/go-dockerclient"
	"github.com/mrturkmencom/defat/virtual/docker"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

type VConfig struct {
	VmConfig struct {
		OvaDir string `yaml:"ova-dir,omitempty"`
	} `yaml:"vm-config"`
	DockerRepositories []dockerclient.AuthConfiguration `yaml:"docker-repositories,omitempty"`
}

func NewConfig(path string) (*VConfig, error) {
	f, err := ioutil.ReadFile(path)

	if err != nil {
		log.Error().Msgf("Reading config file err: %v", err)
		return nil, err
	}

	var c VConfig
	err = yaml.Unmarshal(f, &c)
	if err != nil {
		log.Error().Msgf("Unmarshall error %v \n", err)
		return nil, err
	}
	for _, repo := range c.DockerRepositories {
		docker.Registries[repo.ServerAddress] = repo
	}
	if c.VmConfig.OvaDir == "" {
		return nil, fmt.Errorf("Specify vm directory, err: %v", err)
	}
	return &c, nil
}
