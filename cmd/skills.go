package cmd

import (
	"github.com/spf13/cobra"
)

var skillsCmd = &cobra.Command{
	Use:   "skills",
	Short: "Manage agent skills from SkillHub",
	Long:  "Search, install, and manage agent skills from Vtrix SkillHub",
}

func init() {
	skillsCmd.AddCommand(skillsFindCmd)
	skillsCmd.AddCommand(skillsAddCmd)
	skillsCmd.AddCommand(skillsListCmd)
	skillsCmd.AddCommand(skillsConfigCmd)
}
