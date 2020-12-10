package setup

import (
	"github.com/gofabian/flo/util"
	"github.com/spf13/cobra"
)

var SetupCommand = &cobra.Command{
	Use: "setup-pipeline -g <url> -s <style>",
	Example: util.Dedent(`
		` + "\x00" + `  multibranch:  flo generate-pipeline -g <url> -s multibranch
		` + "\x00" + `  branch:       flo generate-pipeline -g <url> -s branch -b main
	`),
	Short: "Setups a pipeline at a Concourse server with a Drone pipeline as input.",
	Long:  "Setups a pipeline at a Concourse server with a Drone pipeline as input. Supports multibranch pipelines.",
	Args:  cobra.NoArgs,
	Run:   execute,
}

var (
	style  string
	gitURL string
	branch string
	input  string
)

func init() {
	SetupCommand.Flags().SortFlags = false
	SetupCommand.Flags().BoolP("help", "h", false, util.Dedent(`
		Print help text
	`))
	SetupCommand.Flags().StringVarP(&gitURL, "git-url", "g", "", util.Dedent(`
		Git remote `+"`url`"+`, e. g. "https://github.com/org/repo.git", "git@github.com:org/repo.git"
	`))
	SetupCommand.Flags().StringVarP(&style, "style", "s", "", util.Dedent(`
		Choose between multibranch and single branch `+"`style`:"+`

		multibranch  Generates a pipeline that setups pipelines for each branch of the repository.

		branch       Generates a pipeline for a single branch. Requires one "-b" flag.
	`))
	SetupCommand.Flags().StringVarP(&branch, "branch", "b", "", util.Dedent(`
		Git branch name. Combined with "-s branch".
	`))
	SetupCommand.Flags().StringVarP(&input, "input", "i", "", util.Dedent(`
		Path to input `+"`file`"+` (Drone pipeline), default: ".drone.yml"
	`))
}

func execute(cmd *cobra.Command, args []string) {
	cmd.HelpFunc()(cmd, args)
}
