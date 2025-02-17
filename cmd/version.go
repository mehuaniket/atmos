package cmd

import (
	"github.com/spf13/cobra"

	u "github.com/cloudposse/atmos/pkg/utils"
)

var Version = "0.0.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the CLI version",
	Long:  `This command prints the CLI version`,
	Run: func(cmd *cobra.Command, args []string) {
		u.PrintMessage(Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
