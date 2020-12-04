package concourse

type ConditionalSteps struct {
	OnSuccess *Step `yaml:"on_success,omitempty"`
	OnFailure *Step `yaml:"on_failure,omitempty"`
	OnAbort   *Step `yaml:"on_abort,omitempty"`
	OnError   *Step `yaml:"on_error,omitempty"`
	Ensure    *Step `yaml:"ensure,omitempty"`
}

type Step struct {
	// get step
	Get      string            `yaml:"get,omitempty"`
	Resource string            `yaml:"resource,omitempty"`
	Passed   []string          `yaml:"passed,omitempty"`
	Params   map[string]string `yaml:"params,omitempty"`
	Trigger  bool              `yaml:"trigger,omitempty"`
	Version  map[string]string `yaml:"version,omitempty"`

	// task step
	Task       string            `yaml:"task,omitempty"`
	Config     *Task             `yaml:"config,omitempty"`
	File       string            `yaml:"file,omitempty"`
	Image      string            `yaml:"image,omitempty"`
	Privileged bool              `yaml:"privileged,omitempty"`
	Vars       map[string]string `yaml:"vars,omitempty"`
	//Params        map[string]string `yaml:"params,omitempty"`
	InputMapping  map[string]string `yaml:"input_mapping,omitempty"`
	OutputMapping map[string]string `yaml:"output_mapping,omitempty"`

	// set_pipeline step
	SetPipeline string `yaml:"set_pipeline,omitempty"`
	//File        string
	//Vars        map[string]string `yaml:"vars,omitempty"`
	VarFiles []string `yaml:"var_files,omitempty"`
	Team     string   `yaml:"team,omitempty"`

	Timeout          string           `yaml:"timeout,omitempty"`
	Attempts         *int             `yaml:"attempts,omitempty"`
	Tags             []string         `yaml:"tags,omitempty"`
	ConditionalSteps ConditionalSteps `yaml:"-,inline"`
}

type Task struct {
	Platform      Platform
	ImageResource ImageResource     `yaml:"image_resource"`
	Inputs        []Input           `yaml:"inputs,omitempty"`
	Outputs       []Output          `yaml:"outputs,omitempty"`
	Caches        []Cache           `yaml:"caches,omitempty"`
	Params        map[string]string `yaml:"params,omitempty"`
	Run           *Command          `yaml:"run,omitempty"`
	RootfsURI     string            `yaml:"rootfs_uri,omitempty"`
}

type Platform string

const (
	Linux Platform = "linux"
)

type ImageResource struct {
	Type    string
	Source  ImageSource       `yaml:"source,flow"`
	Params  map[string]string `yaml:"params,omitempty"`
	Version map[string]string `yaml:"version,omitempty"`
}

type ImageSource struct {
	Repository string
	Tag        string `yaml:"tag,omitempty"`
}

type Input struct {
	Name     string
	Path     string `yaml:"path,omitempty"`
	Optional *bool  `yaml:"optional,omitempty"`
}

type Output struct {
	Name string
	Path string `yaml:"path,omitempty"`
}

type Cache struct {
	Path string
}

type Command struct {
	Path string
	Args []string `yaml:"args,omitempty"`
	Dir  string   `yaml:"dir,omitempty"`
	User string   `yaml:"user,omitempty"`
}
