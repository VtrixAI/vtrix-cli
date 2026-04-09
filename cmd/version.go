package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/VtrixAI/vtrix-cli/internal/buildinfo"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(buildinfo.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
