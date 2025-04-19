package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

func LoadYamlFile(fileName string) (*Config, error) {
	var config Config

	configYaml, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(configYaml, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
