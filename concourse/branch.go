package concourse

import (
	"fmt"
	"io"
	"strings"
	"text/template"
)

type pipeline struct {
	Name      string
	Steps     []step
	DroneFile string
}

type step struct {
	Name        string
	Repository  string
	Tag         string
	Command     string
	CommandArgs []string
	Commands    []string
}

func CreateBranchPipeline(cfg *Config, writer io.Writer) error {
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

	templateCfg := createTemplateConfig(cfg)

	t, err := template.New(templateName).Parse(branchPipelineTemplate)
	if err != nil {
		return fmt.Errorf("cannot parse template %s: %w", templateName, err)
	}

	err = t.Execute(writer, templateCfg)
	if err != nil {
		return fmt.Errorf("cannot execute template %s: %w", templateName, err)
	}
	return nil
}

func createTemplateConfig(cfg *Config) *pipeline {
	templateCfg := &pipeline{}
	templateCfg.DroneFile = cfg.Input

	if !cfg.BuildJob {
		return templateCfg
	}

	templateCfg.Name = cfg.DronePipeline.Name
	templateCfg.Steps = make([]step, len(cfg.DronePipeline.Steps))

	for i, droneStep := range cfg.DronePipeline.Steps {
		repository, tag := splitImage(droneStep.Image)
		templateCfg.Steps[i].Name = droneStep.Name
		templateCfg.Steps[i].Repository = repository
		templateCfg.Steps[i].Tag = tag
		if len(droneStep.Commands) == 1 {
			args := strings.Split(droneStep.Commands[0], " ")
			templateCfg.Steps[i].Command = args[0]
			if len(args) > 1 {
				templateCfg.Steps[i].CommandArgs = args[1:]
			}
		} else {
			templateCfg.Steps[i].Commands = droneStep.Commands
		}
	}
	return templateCfg
}

func splitImage(image string) (string, string) {
	elements := strings.SplitN(image, ":", 2)
	if len(elements) <= 1 {
		return image, ""
	}
	return elements[0], elements[1]
}
