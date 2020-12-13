package setup

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gofabian/flo/concourse"
	"github.com/gofabian/flo/util"
	"github.com/spf13/cobra"
)

var SetupCommand = &cobra.Command{
	Use: "setup-pipeline -g <url> -s <style>",
	Example: util.Dedent(`
		` + "\x00" + `  multibranch:  flo setup-pipeline -t <target> -g <url> -s multibranch
		` + "\x00" + `  branch:       flo setup-pipeline -t <target> -g <url> -s branch -b main
	`),
	Short: "Setups a pipeline at a Concourse server with a Drone pipeline as input.",
	Long:  "Setups a pipeline at a Concourse server with a Drone pipeline as input. Supports multibranch pipelines.",
	Args:  cobra.NoArgs,
	RunE:  execute,
}

var (
	target  string
	gitURL  string
	style   string
	branch  string
	input   string
	verbose bool
)

func init() {
	SetupCommand.Flags().SortFlags = false
	SetupCommand.Flags().BoolP("help", "h", false, util.Dedent(`
		Print help text
	`))
	SetupCommand.Flags().StringVarP(&target, "fly-target", "t", "", util.Dedent(`
		Fly `+"`target`"+` used to create pipelines
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
		Git `+"`branch`"+` name. Combined with "-s branch".
	`))
	SetupCommand.Flags().StringVarP(&input, "input", "i", "", util.Dedent(`
		Path to input `+"`file`"+` (Drone pipeline), default: ".drone.yml"
	`))
	SetupCommand.Flags().BoolVarP(&verbose, "verbose", "v", false, util.Dedent(`
		Print stdout + stderr of fly
	`))
}

func execute(cmd *cobra.Command, args []string) error {
	if target == "" {
		return fmt.Errorf("'-t' is required")
	}
	if gitURL == "" {
		return fmt.Errorf("'-g' is required")
	}
	if style != "multibranch" && style != "branch" {
		return fmt.Errorf("'-s multibranch|branch' is required")
	}
	if style == "branch" && branch == "" {
		return fmt.Errorf("'-s branch' requires a single '-b' flag")
	}
	if input == "" {
		input = ".drone.yml"
	}

	fmt.Fprintf(os.Stdout, "\nSetup %s pipeline...\n\n", style)
	fmt.Fprintf(os.Stdout, "fly target: %s\n", target)
	fmt.Fprintf(os.Stdout, "Git URL: %s\n", gitURL)
	fmt.Fprintf(os.Stdout, "style: %s\n", style)
	if branch != "" {
		fmt.Fprintf(os.Stdout, "branch: %s\n", branch)
	}
	fmt.Fprintf(os.Stdout, "input: %s\n", input)
	if verbose {
		fmt.Fprintf(os.Stdout, "verbose: true\n")
	}
	fmt.Fprintln(os.Stdout)

	if !isFlyAvailable() {
		fmt.Fprintf(os.Stderr, "'fly' executable cannot be found in PATH!\n")
		os.Exit(1)
	}
	if !isTargetLoggedIn() {
		fmt.Fprintf(os.Stderr, "fly target '%s' is not logged in\n", target)
		os.Exit(1)
	}

	cfg := &concourse.Config{
		SelfUpdateJob: true,
		BuildJob:      false,
		GitURL:        gitURL,
		Branches:      []string{},
		DronePipeline: nil,
	}

	var err error
	pipelineBuffer := &bytes.Buffer{}
	if style == "multibranch" {
		err = concourse.CreateRepositoryPipeline(cfg, pipelineBuffer)
	} else {
		cfg.Branches = []string{branch}
		err = concourse.CreateBranchPipeline(cfg, pipelineBuffer)
	}
	if err != nil {
		return err
	}

	var pipelineName string
	if style == "multibranch" {
		pipelineName = "multibranch"
	} else {
		pipelineName = concourse.HarmonizeGitURL(gitURL + "-" + branch)
	}

	command := newCmd("fly", "-t", target, "get-pipeline", "-p", pipelineName)
	err = command.Run()
	existedPipelineBefore := (err == nil)

	command = newCmd("fly", "-t", target, "set-pipeline", "-p", pipelineName, "-v",
		"GIT_URL=https://github.com/gofabian/flo.git", "-v", "GIT_BRANCH="+branch, "-c=-")
	command.Stdin = pipelineBuffer
	err = command.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fly failed with %s", err)
		os.Exit(1)
	}

	if !existedPipelineBefore {
		command = newCmd("fly", "-t", target, "unpause-pipeline", "-p", pipelineName)
		err = command.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fly failed with %s", err)
			os.Exit(1)
		}
	}

	return nil
}

func isFlyAvailable() bool {
	_, err := exec.LookPath("fly")
	return (err == nil)
}

func isTargetLoggedIn() bool {
	cmd := newCmd("fly", "-t", target, "status")
	err := cmd.Run()
	return (err == nil)
}

func newCmd(path string, args ...string) *exec.Cmd {
	if verbose {
		fmt.Fprintln(os.Stdout)
	}
	fmt.Fprintf(os.Stdout, "RUN: %s %s\n\n", path, strings.Join(args, " "))
	cmd := exec.Command(path, args...)
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdout
	}
	return cmd
}
