package main

import (
	"fmt"
	"os"

	"github.com/gofabian/flo/cmd/generate"
	"github.com/gofabian/flo/cmd/setup"
	"github.com/spf13/cobra"
)

var floCmd = &cobra.Command{
	Use:   "flo",
	Short: "Flo manages Drone pipelines within a Concourse Server",
	Long: `A Fast and Flexible Static Site Generator built with
				  love by spf13 and friends in Go.
				  Complete documentation is available at http://hugo.spf13.com`,
}

func main() {
	cobra.EnableCommandSorting = false
	floCmd.AddCommand(setup.SetupCommand)
	floCmd.AddCommand(generate.GenerateCommand)

	if err := floCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
