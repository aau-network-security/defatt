package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"

	"github.com/aau-network-security/defat/virtual/docker"
	dockerclient "github.com/fsouza/go-dockerclient"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

type Config struct {
	VmConfig           VmConfig                         `yaml:"vm-config"`
	WireguardService   WgConnConf                       `yaml:"wireguard-service"`
	DefatConfig        DefattConf                       `yaml:"defat-config"`
	DockerRepositories []dockerclient.AuthConfiguration `yaml:"docker-repositories"`
}

type VmConfig struct {
	OvaDir string `yaml:"ova-dir"`
	//OvaDir string `yaml:"ova-test"` //use for local test
}

type DefattConf struct {
	Endpoint   string            `yaml:"endpoint"`
	Port       uint64            `yaml:"port"`
	SigningKey string            `yaml:"sign-key"`
	UsersFile  string            `yaml:"users-file"`
	CertConf   CertificateConfig `yaml:"tls"`
	DatabaseFile    string            `yaml:"db-file"`
}

type WgConnConf struct {
	Endpoint string            `yaml:"endpoint"`
	Port     uint64            `yaml:"port"`
	AuthKey  string            `yaml:"auth-key"`
	SignKey  string            `yaml:"sign-key"`
	Dir      string            `yaml:"client-conf-dir"`
	CertConf CertificateConfig `yaml:"tls"`
}

type CertificateConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Directory string `yaml:"directory"`
	CertFile  string `yaml:"certfile"`
	CertKey   string `yaml:"certkey"`
	CAFile    string `yaml:"cafile"`
}

func NewConfig(path string) (*Config, error) {
	f, err := ioutil.ReadFile(path)

	if err != nil {
		log.Error().Msgf("Reading config file err: %v", err)
		return nil, err
	}

	var c Config
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

	if c.WireguardService.CertConf.Enabled {
		if c.WireguardService.CertConf.Directory == "" {
			usr, err := user.Current()
			if err != nil {
				return nil, errors.New("Invalid user")
			}
			c.WireguardService.CertConf.Directory = filepath.Join(usr.HomeDir, ".local", "share", "certmagic")
		}
	}

	return &c, nil
}
