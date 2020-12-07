package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/gofabian/flo/concourse"
	"github.com/gofabian/flo/drone"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var generateBranchCmd = &cobra.Command{
	Use: "branch -g url -b branch -o pipeline.yml",
	Example: `flo generate branch -g https://github.com/org/repo.git -b main
flo generate branch -g git@github.com:org/repo.git -b develop -j all -i .drone.yml`,
	Short: "Generate a Concourse pipeline for a specific branch",
	Long: "Generate a Concourse pipeline for a specific branch and output the YAML " +
		"document to stdout by default.",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateBranch()
	},
}

var generateBranchOptions = &struct {
	gitURL       string
	branch       string
	jobs         string
	pathToInput  string
	pathToOutput string
}{}

func init() {
	cmd := generateBranchCmd
	options := generateBranchOptions
	cmd.Flags().StringVarP(&options.gitURL, "git-url", "g", "",
		"URL to remote git repository")
	cmd.Flags().StringVarP(&options.branch, "branch", "b", "",
		"git branch name")
	cmd.Flags().StringVarP(&options.jobs, "jobs", "j", "refresh",
		`Concourse jobs to generate:
"refresh": job that auto-updates the pipeline
"build": job that runs actual build steps
"all": refresh + build jobs
`)
	cmd.Flags().StringVarP(&options.pathToInput, "input", "i", ".drone.yml",
		"path to Drone pipeline file")
	cmd.Flags().StringVarP(&options.pathToOutput, "output", "o", "",
		"path to Concourse pipeline file")
	cmd.Flags().SortFlags = false
	cmd.MarkFlagRequired("git-url")
	cmd.MarkFlagRequired("branch")
	cmd.MarkFlagRequired("output")
}

func generateBranch() error {
	options := generateBranchOptions

	if options.gitURL == "" || options.branch == "" || options.jobs == "" || options.pathToInput == "" || options.pathToOutput == "" {
		return errors.New("Missing flags")
	}

	jobType := concourse.JobType(options.jobs)
	var dronePipeline *drone.Pipeline

	if jobType != concourse.Refresh {
		// open file
		inputFile, err := os.Open(options.pathToInput)
		if err != nil {
			return fmt.Errorf("cannot open '%s': %w", options.pathToInput, err)
		}
		defer inputFile.Close()

		// decode Drone pipeline
		reader := bufio.NewReader(inputFile)
		dronePipeline = &drone.Pipeline{}
		decoder := yaml.NewDecoder(reader)
		decoder.KnownFields(true)
		err = decoder.Decode(dronePipeline)
		if err != nil {
			return fmt.Errorf("cannot decode drone pipeline: %w", err)
		}

		// validate Drone pipeline
		errs := drone.ValidatePipeline(dronePipeline)
		if len(errs) > 0 {
			msg := "Validation errors: "
			for _, e := range errs {
				msg += ", " + e.Error()
			}
			return errors.New(msg)
		}
	}

	// create Concourse pipeline
	concoursePipeline, err := concourse.CreateBranchPipeline(dronePipeline, jobType)
	if err != nil {
		return err
	}

	// output file
	outputFile, err := os.Create(options.pathToOutput)
	if err != nil {
		return fmt.Errorf("cannot open '%s': %w", options.pathToOutput, err)
	}
	defer outputFile.Close()

	// write to file
	encoder := yaml.NewEncoder(outputFile)
	err = encoder.Encode(concoursePipeline)
	if err != nil {
		return fmt.Errorf("cannot encode concourse pipeline: %w", err)
	}

	return nil
}
