package config

import (
	"dario.cat/mergo"
)

type Instance struct {
	Replicaset string
	Group      string
	Name       string
	Labels     map[string]string
}

func (config *Config) Instances() ([]Instance, error) {
	instances := make([]Instance, 0)

	for groupName, groupConfig := range config.Groups {
		for replicasetName, replicasetConfig := range groupConfig.Replicasets {
			for instanceName, instanceConfig := range replicasetConfig.Instances {
				err := mergo.Merge(&instanceConfig, replicasetConfig.InstanceConfig)
				if err != nil {
					return []Instance{}, err
				}

				err = mergo.Merge(&instanceConfig, groupConfig.InstanceConfig)
				if err != nil {
					return []Instance{}, err
				}

				err = mergo.Merge(&instanceConfig, config.InstanceConfig)
				if err != nil {
					return []Instance{}, err
				}
				// TODO: also apply defaults.

				instance := Instance{
					Name:       instanceName,
					Replicaset: replicasetName,
					Group:      groupName,
				}

				instances = append(instances, instance)
			}
		}
	}

	return instances, nil
}
