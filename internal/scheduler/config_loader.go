package scheduler

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig(
	path string,
) (*SchedulerConfig, error) {

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg SchedulerConfig

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}