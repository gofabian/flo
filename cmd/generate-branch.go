package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var generateBranchCmd = &cobra.Command{
	Use: "branch {-g url | -l [https|ssh]} -b branch",
	Example: `flo generate branch -g https://github.com/org/repo.git -b main
	flo generate branch -g git@github.com:org/repo.git -b develop -j all -i .drone.yml
	flo generate branch -l https -b develop`,
	Short: "Generate a Concourse pipeline for a specific branch",
	Long:  "Generate a Concourse pipeline for a specific branch and output the YAML document to stdout by default.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("generate branch")
	},
}

var generateBranchOptions = struct {
	gitURL           string
	localGitProtocol string
	branch           string
	jobs             string
	pathToInput      string
	pathToOutput     string
}{}

func init() {
	generateBranchCmd.Flags().StringVarP(&generateBranchOptions.gitURL, "git-url", "g", "",
		"URL to remote git repository")
	generateBranchCmd.Flags().StringVarP(&generateBranchOptions.localGitProtocol, "local-git-url", "l", "",
		"read URL from local \".git/config\", "+
			"optionally convert URL to \"https\" or \"ssh\" style, alternative to -g")
	generateBranchCmd.Flags().StringVarP(&generateBranchOptions.branch, "branch", "b", "",
		"git branch name")
	generateBranchCmd.Flags().StringVarP(&generateBranchOptions.jobs, "jobs", "j", "refresh",
		`Concourse jobs to generate:
"refresh": pipeline will contain job to auto-update pipeline only
"all": pipeline will contain refresh job and build job
`)
	generateBranchCmd.Flags().StringVarP(&generateBranchOptions.pathToInput, "input", "i", ".drone.yml",
		"path to Drone pipeline file")
	generateBranchCmd.Flags().StringVarP(&generateBranchOptions.pathToOutput, "output", "o", "",
		"path to Concourse pipeline file")
	generateBranchCmd.Flags().SortFlags = false
}
