package concourse

import (
	"fmt"
	"strings"

	"github.com/gofabian/flo/drone"
	"github.com/gofabian/flo/git"
)

func CreatePipeline(dronePipeline *drone.Pipeline) (*Pipeline, error) {
	gitResource, err := createGitResource(dronePipeline)
	if err != nil {
		return nil, err
	}
	buildJob, err := CreateBuildJob(dronePipeline, gitResource)
	if err != nil {
		return nil, err
	}
	maintenanceJob, err := CreateMaintenanceJob(dronePipeline, gitResource)
	if err != nil {
		return nil, err
	}
	pipeline := Pipeline{
		Resources: []Resource{*gitResource},
		Jobs:      []Job{*maintenanceJob, *buildJob},
	}
	return &pipeline, nil
}

func CreateMaintenanceJob(dronePipeline *drone.Pipeline, gitResource *Resource) (*Job, error) {
	checkoutStep := Step{
		Get:     gitResource.Name,
		Trigger: true,
	}
	floAdapterStep := Step{
		Task: "flo-adapter",
		Config: &Task{
			Platform:      Linux,
			ImageResource: *createImageResource("gofabian/flo:0"),
			Run: &Command{
				Dir:  "workspace",
				Path: "sh",
				Args: []string{
					"-exc",
					strings.Join([]string{
						"flo convert-pipeline -i .drone.yml > ../flo/pipeline.yml",
					}, "\n"),
				},
			},
			Inputs: []Input{{Name: "workspace"}},
			Outputs: []Output{
				{Name: "workspace"},
				{Name: "flo"},
			},
		},
		InputMapping: map[string]string{
			"workspace": gitResource.Name,
		},
	}

	setPipelineStep := Step{
		SetPipeline: "self",
		File:        "flo/pipeline.yml",
		Vars: map[string]string{
			"GIT_BRANCH": "((GIT_BRANCH))",
		},
	}

	allSteps := []Step{checkoutStep, floAdapterStep, setPipelineStep}
	job := Job{
		Name: fmt.Sprintf("maintenance-%s", dronePipeline.Name),
		Plan: allSteps,
	}
	return &job, nil
}

func CreateBuildJob(dronePipeline *drone.Pipeline, gitResource *Resource) (*Job, error) {
	checkoutStep := Step{
		Get:     gitResource.Name,
		Trigger: true,
		Passed:  []string{fmt.Sprintf("maintenance-%s", dronePipeline.Name)},
	}
	taskSteps := createTaskSteps(gitResource.Name, dronePipeline)
	allSteps := append([]Step{checkoutStep}, taskSteps...)

	job := Job{
		Name: dronePipeline.Name,
		Plan: allSteps,
	}
	return &job, nil
}

func createGitResource(dronePipeline *drone.Pipeline) (*Resource, error) {
	gitRepository, err := git.GetRepository()
	if err != nil {
		return nil, err
	}

	gitResource := Resource{
		Name: dronePipeline.Name + "-git",
		Type: "git",
		Source: map[string]string{
			"uri":    gitRepository.URL,
			"branch": "((GIT_BRANCH))", //gitRepository.Branch,
		},
	}
	return &gitResource, nil
}

func createTaskSteps(gitWorkspace string, dronePipeline *drone.Pipeline) []Step {
	taskSteps := make([]Step, len(dronePipeline.Steps))

	previousWorkspace := gitWorkspace
	for i, droneStep := range dronePipeline.Steps {
		nextWorkspace := fmt.Sprintf("workspace%d", i+1)

		taskSteps[i] = Step{
			Task: droneStep.Name,
			Config: &Task{
				Platform:      Linux,
				ImageResource: *createImageResource(droneStep.Image),
				Run:           createCommand(droneStep.Commands),
				Inputs:        []Input{{Name: "workspace"}},
				Outputs:       []Output{{Name: "workspace"}},
			},
			InputMapping: map[string]string{
				"workspace": previousWorkspace,
			},
			OutputMapping: map[string]string{
				"workspace": nextWorkspace,
			},
		}

		previousWorkspace = nextWorkspace
	}

	return taskSteps
}

func createImageResource(image string) *ImageResource {
	return &ImageResource{
		Type:   "registry-image",
		Source: *createSourceFromImage(image),
	}
}

func createSourceFromImage(image string) *ImageSource {
	imageElements := strings.SplitN(image, ":", 2)
	repository := imageElements[0]
	var tag string
	if len(imageElements) > 1 {
		tag = imageElements[1]
	}
	return &ImageSource{
		Repository: repository,
		Tag:        tag,
	}
}

func createCommand(script []string) *Command {
	switch len(script) {
	case 0:
		return nil
	case 1:
		return createSingleCommand(script)
	default:
		return createMultiCommand(script)
	}
}

func createSingleCommand(script []string) *Command {
	elements := strings.SplitN(script[0], " ", 2)

	switch len(elements) {
	case 0:
		return &Command{
			Dir:  "workspace",
			Path: "",
		}
	default:
		return &Command{
			Dir:  "workspace",
			Path: elements[0],
			Args: elements[1:],
		}
	}
}

func createMultiCommand(script []string) *Command {
	text := strings.Join(script, "\n")
	return &Command{
		Dir:  "workspace",
		Path: "sh",
		Args: []string{"-exc", text},
	}
}
