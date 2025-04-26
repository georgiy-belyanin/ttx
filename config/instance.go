package config

import (
	"regexp"

	"dario.cat/mergo"
)

type Instance struct {
	Replicaset  string
	Group       string
	Name        string
	Labels      map[string]string
	ConnectUris []string
}

func (config *Config) FindInstance(name string) (*Instance, error) {
	instances, err := config.Instances()
	if err != nil {
		return nil, err
	}

	for _, instance := range instances {
		if instance.Name == name {
			return &instance, nil
		}
	}

	return nil, nil
}

func (ic *InstanceConfig) ConnectUris() []string {
	connectUris := make([]string, 0)

	iproto := ic.Iproto
	if iproto == nil {
		return connectUris
	}

	advertise := iproto.Advertise
	if advertise != nil {
		advertisedClientUri := ic.Iproto.Advertise.Client
		if advertisedClientUri != nil {
			connectUris = append(connectUris, *advertisedClientUri)
		}
	}

	listen := iproto.Listen
	for _, listen := range listen {
		listenUri := listen.Uri
		if listenUri != nil {
			connectUris = append(connectUris, *listenUri)
		}
	}

	return connectUris
}

var substituteRe *regexp.Regexp = nil

func substituteVars(s string, instanceName string) string {
	if substituteRe == nil {
		var err error
		substituteRe, err = regexp.Compile(`\{\{\s*instance_name\s*\}\}`)
		if err != nil {
			return instanceName
		}
	}

	return substituteRe.ReplaceAllString(s, instanceName)
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

				err = mergo.Merge(&instanceConfig, DefaultInstanceConfig)
				if err != nil {
					return []Instance{}, err
				}

				connectUris := instanceConfig.ConnectUris()
				for i, connectUri := range connectUris {
					connectUris[i] = substituteVars(connectUri, instanceName)
				}

				instance := Instance{
					Name:        instanceName,
					Replicaset:  replicasetName,
					Group:       groupName,
					ConnectUris: connectUris,
				}

				instances = append(instances, instance)
			}
		}
	}

	return instances, nil
}
