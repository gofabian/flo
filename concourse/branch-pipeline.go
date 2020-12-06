package concourse

import (
	"fmt"
	"strings"

	"github.com/gofabian/flo/drone"
)

func CreateBranchPipeline(dronePipeline *drone.Pipeline, jobType JobType) (*Pipeline, error) {
	gitResource := &Resource{
		Name: "checkout",
		Type: "git",
		Source: map[string]string{
			"uri":    "((GIT_URL))",
			"branch": "((GIT_BRANCH))",
		},
	}

	var jobs []Job
	switch jobType {
	case Refresh:
		jobs = []Job{*CreateBranchRefreshJob(gitResource)}
	case Build:
		jobs = []Job{*CreateBranchBuildJob(dronePipeline, gitResource)}
	case All:
		refreshJob := CreateBranchRefreshJob(gitResource)
		buildJob := CreateBranchBuildJob(dronePipeline, gitResource)
		buildJob.Plan[0].Passed = []string{refreshJob.Name}
		jobs = []Job{*refreshJob, *buildJob}
	default:
		return nil, fmt.Errorf("Unknown job type: %s", jobType)
	}

	pipeline := Pipeline{
		Resources: []Resource{*gitResource},
		Jobs:      jobs,
	}
	return &pipeline, nil
}

func CreateBranchRefreshJob(gitResource *Resource) *Job {
	checkoutStep := Step{
		Get:     gitResource.Name,
		Trigger: true,
	}
	generateStep := Step{
		Task: "generate",
		Config: &Task{
			Platform:      Linux,
			ImageResource: *createImageResource("gofabian/flo:0"),
			Run: &Command{
				Dir:  "workspace",
				Path: "sh",
				Args: []string{
					"-exc",
					`#
flo generate branch -g "((GIT_URL))" -b "((GIT_BRANCH))" \
-i .drone.yml -o ../flo/pipeline.yml \
-j all
cat ../flo/pipeline.yml`,
				},
			},
			Inputs:  []Input{{Name: "workspace"}},
			Outputs: []Output{{Name: "workspace"}, {Name: "flo"}},
		},
		InputMapping: map[string]string{"workspace": gitResource.Name},
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
		Plan: []Step{checkoutStep, generateStep, setPipelineStep},
	}
}

func CreateBranchBuildJob(dronePipeline *drone.Pipeline, gitResource *Resource) *Job {
	checkoutStep := Step{
		Get:     gitResource.Name,
		Trigger: true,
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

	for i, droneStep := range dronePipeline.Steps {
		taskSteps[i] = Step{
			Task: droneStep.Name,
			Config: &Task{
				Platform:      Linux,
				ImageResource: *createImageResource(droneStep.Image),
				Run:           createCommand(droneStep.Commands),
				Inputs:        []Input{{Name: "workspace"}},
				Outputs:       []Output{{Name: "workspace"}},
			},
		}
	}

	taskSteps[0].InputMapping = map[string]string{"workspace": gitWorkspace}
	return taskSteps
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
