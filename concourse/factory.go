package concourse

import (
	"errors"
	"strings"

	"github.com/gofabian/flo/drone"
)

func CreatePipeline(dronePipeline *drone.Pipeline) *Pipeline {
	gitResource := Resource{
		Name: dronePipeline.Name + "-git",
		Source: map[string]string{
			"uri": "",
		},
	}

	concourseSteps := make([]Step, len(dronePipeline.Steps)+1)
	concourseSteps[0] = Step{
		Get:     gitResource.Name,
		Trigger: true,
	}

	for i, droneStep := range dronePipeline.Steps {
		concourseSteps[i+1] = Step{
			Task: droneStep.Name,
			Config: &Task{
				Platform: Linux,
				ImageResource: ImageResource{
					Type:   "registry-image",
					Source: createSourceFromImage(droneStep.Image),
				},
				Run: createCommand(&droneStep),
			},
		}
	}

	job := Job{
		Name: dronePipeline.Name,
		Plan: concourseSteps,
	}

	return &Pipeline{
		Resources: []Resource{gitResource},
		Jobs:      []Job{job},
	}
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
			Path: "",
			Args: elements,
		}
	default:
		return &Command{
			Path: elements[0],
			Args: elements[1:],
		}
	}
}

func createMultiCommand(droneStep *drone.Step) *Command {
	script := strings.Join(droneStep.Commands, "\n")
	return &Command{
		Path: "sh",
		Args: []string{"-exc", script},
	}
}
