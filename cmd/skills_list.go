package cmd

import (
	"github.com/spf13/cobra"
	"github.com/VtrixAI/vtrix-cli/internal/skillhub"
)

var (
	listCategory string
	listSort     string
)

var skillsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List skills",
	Long:  "List all available skills from SkillHub",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := skillhub.NewClient()
		return client.List(listCategory, listSort)
	},
}

func init() {
	skillsListCmd.Flags().StringVarP(&listCategory, "category", "c", "", "Filter by category")
	skillsListCmd.Flags().StringVarP(&listSort, "sort", "s", "", "Sort by (stars, downloads, updated)")
}
