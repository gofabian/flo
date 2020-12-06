package cmd

import (
	"github.com/spf13/cobra"
)

var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a Concourse pipeline",
	Long:  "Generate a Concourse pipeline and output the YAML document to stdout by default.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {
	GenerateCmd.AddCommand(generateBranchCmd)
	GenerateCmd.AddCommand(generateRepositoryCmd)
}
