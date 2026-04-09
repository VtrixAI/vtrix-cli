package cmd

import (
	"github.com/spf13/cobra"
	"github.com/VtrixAI/vtrix-cli/internal/skillhub"
)

var (
	findCategory    string
	findInteractive bool
	findCursor      string
)

var skillsFindCmd = &cobra.Command{
	Use:   "find [query]",
	Short: "Search for skills",
	Long:  "Search for skills by keyword or browse by category interactively",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := ""
		if len(args) > 0 {
			query = args[0]
		}

		client := skillhub.NewClient()
		return client.Find(query, findCategory, findInteractive, findCursor)
	},
}

func init() {
	skillsFindCmd.Flags().StringVarP(&findCategory, "category", "c", "", "Filter by category")
	skillsFindCmd.Flags().BoolVarP(&findInteractive, "interactive", "i", false, "Interactive mode (browse by category)")
	skillsFindCmd.Flags().StringVar(&findCursor, "cursor", "", "Page cursor for pagination")
}
