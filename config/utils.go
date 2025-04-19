package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
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

var names = []string{
	"config.yml",
	"config.yaml",
	"source.yml",
	"source.yaml",
}

func FindYamlFileAtPath(path string) (string, error) {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}

	for _, name := range names {
		for _, dirEntry := range dirEntries {
			if dirEntry.IsDir() {
				continue
			}

			fileName := dirEntry.Name()
			if dirEntry.Name() != name {
				continue
			}

			return filepath.Join(path, fileName), nil
		}
	}

	return FindYamlFileAtPath(filepath.Join(path, ".."))
}
