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
	VmConfig struct {
		OvaDir string `yaml:"ova-dir,omitempty"`
	} `yaml:"vm-config"`
	WireguardService   WgConnConf                       `yaml:"wireguard-service,omitempty"`
	DockerRepositories []dockerclient.AuthConfiguration `yaml:"docker-repositories,omitempty"`
	DefatConfig        DefattConf                       `yaml:"defatt-config, omitempty"`
}

type DefattConf struct {
	Endpoint   string            `yaml:"endpoint,omitempty"`
	Port       uint64            `yaml:"port,omitempty"`
	SigningKey string            `yaml:"sign-key,omitempty"`
	UsersFile  string            `yaml:"users-file,omitempty"`
	CertConf   CertificateConfig `yaml:"tls, omitempty"`
}

type WgConnConf struct {
	Endpoint string            `yaml:"endpoint,omitempty"`
	Port     uint64            `yaml:"port,omitempty"`
	AuthKey  string            `yaml:"auth-key,omitempty"`
	SignKey  string            `yaml:"sign-key,omitempty"`
	Dir      string            `yaml:"client-conf-dir,omitempty"`
	CertConf CertificateConfig `yaml:"tls,omitempty"`
}

type CertificateConfig struct {
	Enabled   bool   `yaml:"enabled,omitempty"`
	Directory string `yaml:"directory,omitempty"`
	CertFile  string `yaml:"certfile,omitempty"`
	CertKey   string `yaml:"certkey,omitempty"`
	CAFile    string `yaml:"cafile,omitempty"`
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
