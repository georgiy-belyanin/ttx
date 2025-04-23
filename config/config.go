package config

import yaml "gopkg.in/yaml.v3"

type ReplicasetConfig struct {
	Instances      map[string]InstanceConfig `yaml:"instances"`
	InstanceConfig InstanceConfig
}

func (j *ReplicasetConfig) UnmarshalYAML(value *yaml.Node) error {
	var raw map[string]any
	if err := value.Decode(&raw); err != nil {
		return err
	}
	type Plain ReplicasetConfig
	var plain Plain
	if err := value.Decode(&plain); err != nil {
		return err
	}
	if err := value.Decode(&plain.InstanceConfig); err != nil {
		return err
	}
	*j = ReplicasetConfig(plain)
	return nil
}

type GroupConfig struct {
	Replicasets    map[string]ReplicasetConfig `yaml:"replicasets"`
	InstanceConfig InstanceConfig
}

func (j *GroupConfig) UnmarshalYAML(value *yaml.Node) error {
	var raw map[string]any
	if err := value.Decode(&raw); err != nil {
		return err
	}
	type Plain GroupConfig
	var plain Plain
	if err := value.Decode(&plain); err != nil {
		return err
	}
	if err := value.Decode(&plain.InstanceConfig); err != nil {
		return err
	}
	*j = GroupConfig(plain)
	return nil
}

type Config struct {
	Groups         map[string]GroupConfig `yaml:"groups"`
	InstanceConfig InstanceConfig
}

func (j *Config) UnmarshalYAML(value *yaml.Node) error {
	var raw map[string]any
	if err := value.Decode(&raw); err != nil {
		return err
	}
	type Plain Config
	var plain Plain
	if err := value.Decode(&plain); err != nil {
		return err
	}
	if err := value.Decode(&plain.InstanceConfig); err != nil {
		return err
	}
	*j = Config(plain)
	return nil
}
