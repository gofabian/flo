package concourse

import (
	"fmt"
	"strings"

	"github.com/gofabian/flo/drone"
)

type JobType string

const (
	Refresh JobType = "refresh"
	Build   JobType = "build"
)

func CreateBranchPipeline(dronePipeline *drone.Pipeline, jobTypes []JobType) (*Pipeline, error) {
	gitResource, err := createGitResource()
	if err != nil {
		return nil, err
	}

	jobs := []Job{}
	for _, jobType := range jobTypes {
		switch jobType {
		case Refresh:
			job := CreateRefreshJob(dronePipeline, gitResource)
			jobs = append(jobs, *job)
		case Build:
			job := CreateBuildJob(dronePipeline, gitResource)
			jobs = append(jobs, *job)
		}
	}

	pipeline := Pipeline{
		Resources: []Resource{*gitResource},
		Jobs:      jobs,
	}
	return &pipeline, nil
}

func createGitResource() (*Resource, error) {
	gitResource := Resource{
		Name: "branch-git",
		Type: "git",
		Source: map[string]string{
			"uri":    "((GIT_URL))",
			"branch": "((GIT_BRANCH))",
		},
	}
	return &gitResource, nil
}

func CreateRefreshJob(dronePipeline *drone.Pipeline, gitResource *Resource) *Job {
	checkoutStep := Step{
		Get:     gitResource.Name,
		Trigger: true,
	}
	floAdapterStep := Step{
		Task: "generate",
		Config: &Task{
			Platform:      Linux,
			ImageResource: *createImageResource("gofabian/flo:0"),
			Run: &Command{
				Dir:  "workspace",
				Path: "sh",
				Args: []string{
					"-exc",
					`flo generate branch -g "((GIT_URL))" -b "((GIT_BRANCH))" \
						-i .drone.yml -o ../flo/pipeline.yml \
						-j all`,
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
			"GIT_URL":    "((GIT_URL))",
			"GIT_BRANCH": "((GIT_BRANCH))",
		},
	}

	return &Job{
		Name: "refresh",
		Plan: []Step{checkoutStep, floAdapterStep, setPipelineStep},
	}
}

func CreateBuildJob(dronePipeline *drone.Pipeline, gitResource *Resource) *Job {
	checkoutStep := Step{
		Get:     gitResource.Name,
		Trigger: true,
		Passed:  []string{"refresh"},
	}
	taskSteps := createTaskSteps(gitResource.Name, dronePipeline)
	allSteps := append([]Step{checkoutStep}, taskSteps...)

	job := Job{
		Name: dronePipeline.Name,
		Plan: allSteps,
	}
	return &job
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
			InputMapping:  map[string]string{"workspace": previousWorkspace},
			OutputMapping: map[string]string{"workspace": nextWorkspace},
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
		return createSingleCommand(script[0])
	default:
		return createMultiCommand(script)
	}
}

func createSingleCommand(command string) *Command {
	elements := strings.SplitN(command, " ", 2)

	if len(elements) == 1 {
		return &Command{
			Dir:  "workspace",
			Path: elements[0],
		}
	}
	return &Command{
		Dir:  "workspace",
		Path: elements[0],
		Args: elements[1:],
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
