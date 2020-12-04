package concourse

type Pipeline struct {
	Jobs          []Job
	Resources     []Resource     `yaml:"resource,omitempty"`
	ResourceTypes []ResourceType `yaml:"resource_types,omitempty"`
	Groups        []GroupConfig  `yaml:"groups,omitempty"`
}

type Job struct {
	Name             string
	Plan             []Step
	OldName          string           `yaml:"old_name,omitempty"`
	Serial           *bool            `yaml:"serial,omitempty"`
	Public           *bool            `yaml:"public,omitempty"`
	ConditionalSteps ConditionalSteps `yaml:"-,inline"`
}

type Resource struct {
	Name       string
	Type       string
	Source     map[string]string
	OldName    string            `yaml:"old_name,omitempty"`
	Icon       string            `yaml:"icon,omitempty"`
	Version    map[string]string `yaml:"version,omitempty"`
	CheckEvery string            `yaml:"check_every,omitempty"`
	Tags       []string          `yaml:"tags,omitempty"`
	Public     *bool             `yaml:"public,omitempty"`
}

type ResourceType struct {
	Name       string
	Type       string
	Source     map[string]string
	Privileged *bool             `yaml:"privileged,omitempty"`
	Params     map[string]string `yaml:"params,omitempty"`
	CheckEvery string            `yaml:"check_every,omitempty"`
	Tags       []string          `yaml:"tags,omitempty"`
	Defaults   map[string]string `yaml:"defaults,omitempty"`
}

type GroupConfig struct {
	Name string
	Jobs []string `yaml:"jobs,omitempty"`
}
