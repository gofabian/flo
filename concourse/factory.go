package concourse

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofabian/flo/drone"
	"github.com/gofabian/flo/git"
)

const workspaceName = "workspace"

func CreatePipeline(dronePipeline *drone.Pipeline) (*Pipeline, error) {
	gitRepository, err := git.GetRepository()
	if err != nil {
		return nil, err
	}

	gitResource := Resource{
		Name: dronePipeline.Name + "-git",
		Type: "git",
		Source: map[string]string{
			"uri":    gitRepository.URL,
			"branch": gitRepository.Branch,
		},
	}

	concourseSteps := make([]Step, len(dronePipeline.Steps)+1)
	concourseSteps[0] = Step{
		Get:     gitResource.Name,
		Trigger: true,
	}

	previousWorkspace := gitResource.Name
	for i, droneStep := range dronePipeline.Steps {
		nextWorkspace := fmt.Sprintf("%s%d", workspaceName, i+1)

		concourseSteps[i+1] = Step{
			Task: droneStep.Name,
			Config: &Task{
				Platform: Linux,
				ImageResource: ImageResource{
					Type:   "registry-image",
					Source: createSourceFromImage(droneStep.Image),
				},
				Run:     createCommand(&droneStep),
				Inputs:  []Input{{Name: workspaceName}},
				Outputs: []Output{{Name: workspaceName}},
			},
			InputMapping: map[string]string{
				workspaceName: previousWorkspace,
			},
			OutputMapping: map[string]string{
				workspaceName: nextWorkspace,
			},
		}

		previousWorkspace = nextWorkspace
	}

	job := Job{
		Name: dronePipeline.Name,
		Plan: concourseSteps,
	}
	pipeline := Pipeline{
		Resources: []Resource{gitResource},
		Jobs:      []Job{job},
	}
	return &pipeline, nil
}

func createSourceFromImage(image string) ImageSource {
	imageElements := strings.SplitN(image, ":", 2)
	repository := imageElements[0]
	var tag string
	if len(imageElements) > 1 {
		tag = imageElements[1]
	}
	return ImageSource{
		Repository: repository,
		Tag:        tag,
	}
}

func createCommand(droneStep *drone.Step) *Command {
	switch len(droneStep.Commands) {
	case 0:
		return createPluginCommand(droneStep)
	case 1:
		return createSingleCommand(droneStep)
	default:
		return createMultiCommand(droneStep)
	}
}

func createPluginCommand(droneStep *drone.Step) *Command {
	// todo: plugin task
	panic(errors.New("Drone plugins not implemented"))
}

func createSingleCommand(droneStep *drone.Step) *Command {
	elements := strings.SplitN(droneStep.Commands[0], " ", 2)

	switch len(elements) {
	case 0:
		return &Command{
			Dir:  workspaceName,
			Path: "",
		}
	default:
		return &Command{
			Dir:  workspaceName,
			Path: elements[0],
			Args: elements[1:],
		}
	}
}

func createMultiCommand(droneStep *drone.Step) *Command {
	script := strings.Join(droneStep.Commands, "\n")
	return &Command{
		Dir:  workspaceName,
		Path: "sh",
		Args: []string{"-exc", script},
	}
}
