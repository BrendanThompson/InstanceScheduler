package azure

import (
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
	yaml "gopkg.in/yaml.v3"
)

type Tags struct {
	InstanceSchedulingEnabled     string `yaml:"enabled"`
	InstanceSchedulingSchedule    string `yaml:"schedule"`
	InstanceSchedulingPatchWindow string `yaml:"patchWindow"`
}

func NewTagsFromConfig(path string) (*Tags, error) {
	var tags Tags
	var err error

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &tags)
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("Loaded tags: %+v", tags)

	return &tags, nil
}

func (t *Tags) LoadValues(tags map[string]*string) (bool, string, string) {
	var enabled bool
	var schedule string
	var patchWindow string

	for key, value := range tags {
		switch key {
		case t.InstanceSchedulingEnabled:
			enabled, _ = strconv.ParseBool(*value)
		case t.InstanceSchedulingSchedule:
			schedule = *value
		case t.InstanceSchedulingPatchWindow:
			patchWindow = *value
		}
	}

	return enabled, schedule, patchWindow
}
