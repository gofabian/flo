package generate

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gofabian/flo/concourse"
	"github.com/gofabian/flo/drone"
	"github.com/gofabian/flo/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var GenerateCommand = &cobra.Command{
	Use: "generate-pipeline -s <style>",
	Example: util.Dedent(`
		` + "\x00" + `  Minimal pipelines:  flo generate-pipeline -s multibranch
		                      flo generate-pipeline -s branch
		` + "\x00" + `  Fully populated:    flo generate-pipeline -s multibranch -j self-update,build -b main,develop,feature/x
		                      flo generate-pipeline -s branch -j self-update -j build
	`),
	Short: "Generates a Concourse pipeline with a Drone pipeline as input.",
	Long:  "Generates a Concourse pipeline with a Drone pipeline as input. Supports multibranch pipelines.",
	Args:  cobra.NoArgs,
	RunE:  execute,
}

var (
	style    string
	branches = []string{}
	jobs     = []string{}
	input    string
	output   string

	selfUpdateJob = false
	buildJob      = false
	stdout        = false
)

func init() {
	GenerateCommand.Flags().SortFlags = false
	GenerateCommand.Flags().BoolP("help", "h", false, util.Dedent(`
		Print help text
	`))
	GenerateCommand.Flags().StringVarP(&style, "style", "s", "", util.Dedent(`
		Choose between multibranch and single branch `+"`style`:"+`

		multibranch  Generates a pipeline that setups pipelines for each branch of the repository.

		branch       Generates a pipeline for a single branch.
	`))
	GenerateCommand.Flags().StringSliceVarP(&jobs, "job", "j", nil, util.Dedent(`
		Select `+"`jobs`"+` to be part of the generated pipeline. Multiple flags are possible:
		"--job self-update --job build" or "--job self-update,build"

		self-update  Default option. Generates a job that updates the pipeline itself with a 
		             "set-pipeline: self" step. After self-update the pipeline will have both
		             of "self-update" job and "build" job.

		build        Job that setups a pipeline for each branch (in multibranch pipeline) or 
		             job that executes the actual build steps (in single branch pipeline).
					 Requires 1..n "-b" flags for multibranch pipeline.
					 Requires input file (default: "-i .drone.yml") to be present in working 
					 directory.
	`))
	GenerateCommand.Flags().StringSliceVarP(&branches, "branch", "b", []string{}, util.Dedent(`
		Git branch names. Requires 1..n `+"`branches`"+` combined with "-s multibranch -j build".
	`))
	GenerateCommand.Flags().StringVarP(&input, "input", "i", "", util.Dedent(`
		Path to input `+"`file`"+` (Drone pipeline), default: ".drone.yml"
	`))
	GenerateCommand.Flags().StringVarP(&output, "output", "o", "", util.Dedent(`
		Path to output `+"`file`"+` (Concourse pipeline), default: <stdout>
	`))
}

func execute(cmd *cobra.Command, args []string) error {
	if len(jobs) == 0 {
		jobs = []string{"self-update"}
		selfUpdateJob = true
	} else {
		for _, j := range jobs {
			switch j {
			case "self-update":
				selfUpdateJob = true
			case "build":
				buildJob = true
			default:
				return fmt.Errorf("Invalid job type '%s'", j)
			}
		}
	}

	if input == "" {
		input = ".drone.yml"
	}
	if output == "" {
		stdout = true
	}

	if style != "multibranch" && style != "branch" {
		return fmt.Errorf("'-s multibranch|branch' is required")
	}
	if style == "multibranch" && buildJob && len(branches) == 0 {
		return fmt.Errorf("'-s multibranch -j build' requires at least one '-b' flag")
	}
	if style == "multibranch" && !buildJob {
		branches = []string{}
	}

	var outputFile *os.File
	if stdout {
		outputFile = os.Stdout
	} else {
		var err error
		outputFile, err = os.Create(output)
		if err != nil {
			return fmt.Errorf("cannot open output file '%s': %w", output, err)
		}
		defer outputFile.Close()
	}

	fmt.Fprintf(os.Stderr, "\nGenerating %s pipeline...\n\n", style)
	fmt.Fprintf(os.Stderr, "style: %s\n", style)
	if len(branches) > 0 {
		fmt.Fprintf(os.Stderr, "branches: %s\n", strings.Join(branches, ", "))
	}
	fmt.Fprintf(os.Stderr, "jobs: %s\n", strings.Join(jobs, ", "))
	fmt.Fprintf(os.Stderr, "input: %s\n", input)
	if stdout {
		fmt.Fprintf(os.Stderr, "output: <stdout>\n\n")
	} else {
		fmt.Fprintf(os.Stderr, "output: %s\n\n", output)
	}

	cfg := &concourse.Config{
		SelfUpdateJob: selfUpdateJob,
		BuildJob:      buildJob,
		Branches:      branches,
		DronePipeline: nil,
	}

	switch style {
	case "branch":
		dronePipeline, err := getDronePipeline()
		if err != nil {
			return err
		}
		cfg.DronePipeline = dronePipeline
		return concourse.CreateBranchPipeline(cfg, outputFile)
	case "multibranch":
		return concourse.CreateRepositoryPipeline(cfg, outputFile)
	}

	return nil
}

func getDronePipeline() (*drone.Pipeline, error) {
	if !buildJob {
		return nil, nil
	}

	inputFile, err := os.Open(input)
	if err != nil {
		return nil, fmt.Errorf("cannot open input file '%s': %w", input, err)
	}
	defer inputFile.Close()

	reader := bufio.NewReader(inputFile)
	dronePipeline := &drone.Pipeline{}
	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)
	err = decoder.Decode(dronePipeline)
	if err != nil {
		return nil, fmt.Errorf("cannot decode drone pipeline: %w", err)
	}

	errs := drone.ValidatePipeline(dronePipeline)
	if len(errs) > 0 {
		msg := "Validation errors: "
		for _, e := range errs {
			msg += ", " + e.Error()
		}
		return nil, errors.New(msg)
	}
	return dronePipeline, nil
}
