package concourse

import (
	"fmt"
	"regexp"
)

func CreateRepositoryPipeline(jobType JobType, branches []string) (*Pipeline, error) {
	branchesResourceType := &ResourceType{
		Name:   "git-branches",
		Type:   "registry-image",
		Source: *createSourceFromImage("vito/git-branches-resource"),
	}
	branchesResource := &Resource{
		Name: "branches",
		Type: branchesResourceType.Name,
		Source: map[string]string{
			"uri": "((GIT_URL))",
		},
	}

	refreshJob := createRepositoryRefreshJob(branchesResource)
	buildJob := CreateRepostoryBuildJob(branchesResource, branches)

	var jobs []Job
	switch jobType {
	case Refresh:
		jobs = []Job{*refreshJob}
	case Build:
		jobs = []Job{*buildJob}
	case All:
		buildJob.Plan[0].Passed = []string{refreshJob.Name}
		jobs = []Job{*refreshJob, *buildJob}
	default:
		return nil, fmt.Errorf("Unknown job type: %s", jobType)
	}

	pipeline := Pipeline{
		ResourceTypes: []ResourceType{*branchesResourceType},
		Resources:     []Resource{*branchesResource},
		Jobs:          jobs,
	}
	return &pipeline, nil
}

func createRepositoryRefreshJob(branchesResource *Resource) *Job {
	checkoutStep := Step{
		Get:     branchesResource.Name,
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
bs=$(tr '\n' ',' < branches | sed -e 's/,*$//' | sed -e 's/,/ -b /g')
flo generate repository -g "((GIT_URL))" -b $bs \
	-i .drone.yml -o ../flo/pipeline.yml \
	-j all
cat ../flo/pipeline.yml`,
				},
			},
			Inputs:  []Input{{Name: "workspace"}},
			Outputs: []Output{{Name: "workspace"}, {Name: "flo"}},
		},
		InputMapping: map[string]string{"workspace": branchesResource.Name},
	}

	setPipelineStep := Step{
		SetPipeline: "self",
		File:        "flo/pipeline.yml",
		Vars: map[string]string{
			"GIT_URL": "((GIT_URL))",
		},
	}

	return &Job{
		Name: "refresh",
		Plan: []Step{checkoutStep, generateStep, setPipelineStep},
	}
}

func CreateRepostoryBuildJob(branchesResource *Resource, branches []string) *Job {
	checkoutStep := Step{
		Get:     branchesResource.Name,
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
flo generate branch -g "((GIT_URL))" -b dummy \
	-i .drone.yml -o ../flo/pipeline.yml \
	-j refresh
cat ../flo/pipeline.yml`,
				},
			},
			Inputs:  []Input{{Name: "workspace"}},
			Outputs: []Output{{Name: "workspace"}, {Name: "flo"}},
		},
		InputMapping: map[string]string{"workspace": branchesResource.Name},
	}

	setPipelineSteps := make([]Step, len(branches))

	// todo: improve
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)

	for i, branch := range branches {
		setPipelineSteps[i] = Step{
			SetPipeline: re.ReplaceAllString(branch, "-"),
			File:        "flo/pipeline.yml",
			Vars: map[string]string{
				"GIT_URL":    "((GIT_URL))",
				"GIT_BRANCH": branch,
			},
		}
	}

	job := Job{
		Name: "pipelines",
		Plan: append([]Step{checkoutStep, generateStep}, setPipelineSteps...),
	}
	return &job
}
