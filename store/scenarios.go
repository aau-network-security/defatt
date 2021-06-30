package store

import (
	"errors"
	"io/ioutil"

	"github.com/rs/zerolog/log"
	yaml "gopkg.in/yaml.v2"
)

var (
	scenarios         = map[int]Scenario{}
	ErrUnkownScenario = errors.New("No scenario with that id")
)

type Scenario struct {
	ID         uint32
	Networks   map[string]network `yaml:"networks"`
	Story      string             `yaml:"story"`
	Duration   uint32             `yaml:"duration"`
	Difficulty string             `yaml:"difficulty"`
}

type network struct {
	Vlan  string
	Chals []string `yaml:"challenges"`
}

// GetScenariosFromFile will parse the given file into a map of Scenario
func GetScenariosFromFile(path string) error {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(f, scenarios)
	if err != nil {
		return err
	}

	log.Debug().Int("amount", len(scenarios)).Msg("read senarios from file")

	return nil
}

func GetScenarioByID(id int) (Scenario, error) {
	scenario, ok := scenarios[id]
	if !ok {
		return Scenario{}, ErrUnkownScenario
	}

	log.Debug().Uint32("ID", scenario.ID).Msg("got scenario from store")
	return scenario, nil
}

func GetAllScenarios() []Scenario {
	var scenarioArr []Scenario
	for k, v := range scenarios {
		scenario := Scenario{
			ID:         uint32(k),
			Networks:   v.Networks,
			Story:      v.Story,
			Duration:   v.Duration,
			Difficulty: v.Difficulty,
		}
		scenarioArr = append(scenarioArr, scenario)
	}

	return scenarioArr
}
