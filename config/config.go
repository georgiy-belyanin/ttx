package config

type Instance struct {
}

type Replicaset struct {
	Instances map[string]Instance `yaml:"instances"`
}

type Group struct {
	Replicasets map[string]Replicaset `yaml:"replicasets"`
}

type Config struct {
	Groups map[string]Group `yaml:"groups"`
}

func (c *Config) InstanceNames() []string {
	instanceNames := make([]string, 0)
	for _, group := range c.Groups {
		for _, replicaset := range group.Replicasets {
			for instanceName := range replicaset.Instances {
				instanceNames = append(instanceNames, instanceName)
			}
		}
	}
	return instanceNames
}
