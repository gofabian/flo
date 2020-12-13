package concourse

import (
	"fmt"
	"io"
	"text/template"
)

type repository struct {
	DroneFile string
	Branches  []branch
}

type branch struct {
	Name           string
	HarmonizedName string
	DroneFile      string
}

func CreateRepositoryPipeline(cfg *Config, writer io.Writer) error {
	var templateName string
	if cfg.SelfUpdateJob && cfg.BuildJob {
		templateName = "full-pipeline"
	} else if cfg.SelfUpdateJob {
		templateName = "self-update-pipeline"
	} else if cfg.BuildJob {
		templateName = "build-pipeline"
	} else {
		return fmt.Errorf("missing template")
	}

	templateCfg := repository{}
	templateCfg.DroneFile = cfg.Input
	templateCfg.Branches = make([]branch, len(cfg.Branches))

	for i, branch := range cfg.Branches {
		templateCfg.Branches[i].Name = branch
		templateCfg.Branches[i].HarmonizedName = HarmonizeName(branch)
		templateCfg.Branches[i].DroneFile = cfg.Input
	}

	t, err := template.New(templateName).Parse(repositoryPipelineTemplate)
	if err != nil {
		return fmt.Errorf("cannot parse template %s: %w", templateName, err)
	}

	err = t.Execute(writer, templateCfg)
	if err != nil {
		return fmt.Errorf("cannot execute template %s: %w", templateName, err)
	}
	return nil
}
