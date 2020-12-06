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

var gitURL string
var localGitProtocol string
var branch string
var jobs string
var pathToInput string
var pathToOutput string

func init() {
	generateBranchCmd.Flags().StringVarP(&gitURL, "git-url", "g", "",
		"URL to remote git repository")
	generateBranchCmd.Flags().StringVarP(&branch, "branch", "b", "",
		"git branch name")
	generateBranchCmd.Flags().StringVarP(&jobs, "jobs", "j", "refresh",
		`Concourse jobs to generate:
"refresh": job that auto-updates the pipeline
"build": job that runs actual build steps
"all": refresh + build jobs
`)
	generateBranchCmd.Flags().StringVarP(&pathToInput, "input", "i", ".drone.yml",
		"path to Drone pipeline file")
	generateBranchCmd.Flags().StringVarP(&pathToOutput, "output", "o", "",
		"path to Concourse pipeline file")
	generateBranchCmd.Flags().SortFlags = false
	generateBranchCmd.MarkFlagRequired("git-url")
	generateBranchCmd.MarkFlagRequired("branch")
	generateBranchCmd.MarkFlagRequired("output")
}

func generateBranch() error {
	if gitURL == "" || branch == "" || jobs == "" || pathToInput == "" || pathToOutput == "" {
		return errors.New("Missing flags")
	}

	// open file
	inputFile, err := os.Open(pathToInput)
	if err != nil {
		return fmt.Errorf("cannot open '%s': %w", pathToInput, err)
	}
	defer inputFile.Close()

	// decode Drone pipeline
	reader := bufio.NewReader(inputFile)
	dronePipeline := &drone.Pipeline{}
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

	// create Concourse pipeline
	jobType := concourse.JobType(jobs)
	concoursePipeline, err := concourse.CreateBranchPipeline(dronePipeline, jobType)
	if err != nil {
		return err
	}

	// output file
	outputFile, err := os.Create(pathToOutput)
	if err != nil {
		return fmt.Errorf("cannot open '%s': %w", pathToOutput, err)
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
