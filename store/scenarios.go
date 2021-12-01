package store

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/aau-network-security/defatt/models"
	"github.com/rs/zerolog/log"
	yaml "gopkg.in/yaml.v2"
)

var (
	ErrUnkownScenario = errors.New("no scenario with that id")
)

type Scenario struct {
	Name       string           `yaml:"name"`
	Topic      string           `yaml:"topic"`
	StoryRed   string           `yaml:"story-red"`
	StoryBlue  string           `yaml:"story-blue"`
	Duration   uint32           `yaml:"duration"`
	Difficulty string           `yaml:"difficulty"`
	Networks   []models.Network `yaml:"networks"`
	Hosts      []models.Host    `yaml:"hosts"`
}

// LoadScenarios will all files in a directory into a map of Scenario
func LoadScenarios(path string) (map[int]Scenario, error) {
	scenarios := make(map[int]Scenario)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for i, file := range files {
		var scenario Scenario
		log.Debug().Str("file", path+"/"+file.Name()).Msg("readig scenario from file")
		f, err := ioutil.ReadFile(path + "/" + file.Name())
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(f, &scenario)
		if err != nil {
			return nil, err
		}
		scenarios[i] = scenario
	}

	log.Debug().Int("amount", len(scenarios)).Msg("read senarios from file")

	return scenarios, nil
}
