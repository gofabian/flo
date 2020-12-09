package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/gofabian/flo/concourse"
	"github.com/spf13/cobra"
)

var generateRepoCmd = &cobra.Command{
	Use: "repository -g url -b branch -o pipeline.yml",
	Example: `flo generate branch -g https://github.com/org/repo.git -b main
flo generate branch -g git@github.com:org/repo.git -b develop -j all -i .drone.yml`,
	Short: "Generate a Concourse pipeline for a specific branch",
	Long: "Generate a Concourse pipeline for a specific branch and output the YAML " +
		"document to stdout by default.",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateRepository()
	},
}

var generateRepoOptions = &struct {
	gitURL       string
	branches     []string
	jobs         string
	pathToInput  string
	pathToOutput string
}{}

func init() {
	cmd := generateRepoCmd
	options := generateRepoOptions
	cmd.Flags().StringVarP(&options.gitURL, "git-url", "g", "",
		"URL to remote git repository")
	cmd.Flags().StringSliceVarP(&options.branches, "branch", "b", []string{},
		"git branch names")
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

func generateRepository() error {
	options := generateRepoOptions

	if options.gitURL == "" || len(options.branches) == 0 || options.jobs == "" || options.pathToInput == "" || options.pathToOutput == "" {
		return errors.New("Missing flags")
	}

	jobType := JobType(options.jobs)

	var templateName string
	switch jobType {
	case All:
		templateName = "full-pipeline"
	case Refresh:
		templateName = "self-update-pipeline"
	case Build:
		templateName = "build-pipeline"
	default:
		return fmt.Errorf("invalid jobs value: %s", jobType)
	}

	// output file
	outputFile, err := os.Create(options.pathToOutput)
	if err != nil {
		return fmt.Errorf("cannot open '%s': %w", options.pathToOutput, err)
	}
	defer outputFile.Close()

	// create Concourse pipeline
	concourse.CreateRepositoryPipeline(templateName, options.branches, outputFile)
	if err != nil {
		return err
	}

	return nil
}
