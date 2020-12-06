package cmd

import (
	"github.com/spf13/cobra"
)

var SetupCmd = &cobra.Command{
	Use:     "setup",
	Example: "",
	Short:   "Setup multi-branch pipeline in Concourse",
	Long:    "Setup multi-branch pipeline in Concourse",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.HelpFunc()(cmd, args)
	},
}

var setupOptions = &struct {
	flyTarget        string
	gitURL           string
	localGitProtocol string
	pathToInput      string
}{}

func init() {
	SetupCmd.Flags().StringVarP(&setupOptions.flyTarget, "target", "t", "",
		"fly target (must be logged in)")
	SetupCmd.Flags().StringVarP(&setupOptions.gitURL, "git-url", "g", "",
		"URL to remote git repository")
	SetupCmd.Flags().StringVarP(&setupOptions.localGitProtocol, "local-git-url", "l", "",
		"read URL from local \".git/config\", "+
			"optionally convert URL to \"https\" or \"ssh\" style, alternative to -g")
	SetupCmd.Flags().StringVarP(&setupOptions.pathToInput, "input", "i", ".drone.yml",
		"path to Drone pipeline file")
	SetupCmd.Flags().SortFlags = false
}
