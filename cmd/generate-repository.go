package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var generateRepositoryCmd = &cobra.Command{
	Use: "repository {-g url | -l [https|ssh]} -b branch",
	Example: `flo generate repository -g https://github.com/org/repo.git -b main
	flo generate repository -g git@github.com:org/repo.git -b develop -j all -i .drone.yml
	flo generate repository -l https -b develop`,
	Short: "Generate a Concourse pipeline for a specific repository",
	Long:  "Generate a Concourse pipeline for a specific repository and output the YAML document to stdout by default.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("generate branch")
	},
}

var generateRepositoryOptions = struct {
	gitURL           string
	localGitProtocol string
	branch           string
	jobs             string
	pathToInput      string
	pathToOutput     string
}{}

func init() {
	generateRepositoryCmd.Flags().StringVarP(&generateRepositoryOptions.gitURL, "git-url", "g", "",
		"URL to remote git repository")
	generateRepositoryCmd.Flags().StringVarP(&generateRepositoryOptions.localGitProtocol, "local-git-url", "l", "",
		"read URL from local \".git/config\", "+
			"optionally convert URL to \"https\" or \"ssh\" style, alternative to -g")
	generateRepositoryCmd.Flags().StringVarP(&generateRepositoryOptions.branch, "branch", "b", "",
		"git branch name")
	generateRepositoryCmd.Flags().StringVarP(&generateRepositoryOptions.jobs, "jobs", "j", "refresh",
		`Concourse jobs to generate:
"refresh": pipeline will contain job to auto-update pipeline only
"all": pipeline will contain refresh job and build job
`)
	generateRepositoryCmd.Flags().StringVarP(&generateRepositoryOptions.pathToInput, "input", "i", ".drone.yml",
		"path to Drone pipeline file")
	generateRepositoryCmd.Flags().StringVarP(&generateRepositoryOptions.pathToOutput, "output", "o", "",
		"path to Concourse pipeline file")
	generateRepositoryCmd.Flags().SortFlags = false
}
