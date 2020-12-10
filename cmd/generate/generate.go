package generate

import (
	"github.com/gofabian/flo/util"
	"github.com/spf13/cobra"
)

var GenerateCommand = &cobra.Command{
	Use: "generate-pipeline -g <url> -s <style>",
	Example: util.Dedent(`
		` + "\x00" + `  Minimal pipelines:  flo generate-pipeline -g <url> -s multibranch
		                      flo generate-pipeline -g <url> -s branch -b main
		` + "\x00" + `  Fully populated:    flo generate-pipeline -g <url> -s multibranch -j self-update,branch -b main,develop,feature/x
		                      flo generate-pipeline -g <url> -s branch -j self-update -jbranch -b main
	`),
	Short: "Generates a Concourse pipeline with a Drone pipeline as input.",
	Long:  "Generates a Concourse pipeline with a Drone pipeline as input. Supports multibranch pipelines.",
	Args:  cobra.NoArgs,
	Run:   execute,
}

var (
	style    string
	jobs     []string
	gitURL   string
	branches []string
	input    string
	output   string
)

func init() {
	GenerateCommand.Flags().SortFlags = false
	GenerateCommand.Flags().BoolP("help", "h", false, util.Dedent(`
		Print help text
	`))
	GenerateCommand.Flags().StringVarP(&style, "style", "s", "", util.Dedent(`
		Choose between multibranch and single branch `+"`style`:"+`

		multibranch  Generates a pipeline that setups pipelines for each branch of the repository.

		branch       Generates a pipeline for a single branch. Requires one "-b" flag.
	`))
	GenerateCommand.Flags().StringVarP(&gitURL, "git-url", "g", "", util.Dedent(`
		Git remote `+"`url`"+`, e. g. "https://github.com/org/repo.git", "git@github.com:org/repo.git"
	`))
	GenerateCommand.Flags().StringSliceVarP(&branches, "branch", "b", nil, util.Dedent(`
		Git branch name. Requires 1 branch combined with "-s branch" or 1..n `+"`branches`"+` combined with 
		"-s multibranch -j build".
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
	`))
	GenerateCommand.Flags().StringVarP(&input, "input", "i", "", util.Dedent(`
		Path to input `+"`file`"+` (Drone pipeline), default: ".drone.yml"
	`))
	GenerateCommand.Flags().StringVarP(&output, "output", "o", "", util.Dedent(`
		Path to output `+"`file`"+` (Concourse pipeline), default: <stdout>
	`))
}

func execute(cmd *cobra.Command, args []string) {
	cmd.HelpFunc()(cmd, args)
}
