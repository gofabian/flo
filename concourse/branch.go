package concourse

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/gofabian/flo/drone"
)

type pipeline struct {
	Name  string
	Steps []step
}

type step struct {
	Name        string
	Repository  string
	Tag         string
	Command     string
	CommandArgs []string
	Commands    []string
}

func CreateBranchPipeline(selfUpdateJob bool, dronePipeline *drone.Pipeline, writer io.Writer) error {
	var templateName string
	if selfUpdateJob && dronePipeline != nil {
		templateName = "full-pipeline"
	} else if selfUpdateJob {
		templateName = "self-update-pipeline"
	} else if dronePipeline != nil {
		templateName = "build-pipeline"
	} else {
		return fmt.Errorf("missing template")
	}

	cfg := createTemplateConfig(dronePipeline)

	t, err := template.New(templateName).Parse(branchPipelineTemplate)
	if err != nil {
		return fmt.Errorf("cannot parse template %s: %w", templateName, err)
	}

	err = t.Execute(writer, cfg)
	if err != nil {
		return fmt.Errorf("cannot execute template %s: %w", templateName, err)
	}
	return nil
}

func createTemplateConfig(dronePipeline *drone.Pipeline) *pipeline {
	cfg := &pipeline{}

	if dronePipeline == nil {
		return cfg
	}

	cfg.Name = dronePipeline.Name
	cfg.Steps = make([]step, len(dronePipeline.Steps))

	for i, droneStep := range dronePipeline.Steps {
		repository, tag := splitImage(droneStep.Image)
		cfg.Steps[i].Name = droneStep.Name
		cfg.Steps[i].Repository = repository
		cfg.Steps[i].Tag = tag
		if len(droneStep.Commands) == 1 {
			args := strings.Split(droneStep.Commands[0], " ")
			cfg.Steps[i].Command = args[0]
			if len(args) > 1 {
				cfg.Steps[i].CommandArgs = args[1:]
			}
		} else {
			cfg.Steps[i].Commands = droneStep.Commands
		}
	}
	return cfg
}

func splitImage(image string) (string, string) {
	elements := strings.SplitN(image, ":", 2)
	if len(elements) <= 1 {
		return image, ""
	}
	return elements[0], elements[1]
}
