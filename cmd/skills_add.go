package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/VtrixAI/vtrix-cli/internal/skillhub"
)

var (
	addVersion string
	addGlobal  bool
	addYes     bool
)

func dryRunSkillAdd(slug, version string, global bool) error {
	fmt.Fprintf(os.Stderr, "[dry-run] Would install skill: %s\n", slug)
	if version != "" {
		fmt.Fprintf(os.Stderr, "[dry-run]   Version: %s\n", version)
	} else {
		fmt.Fprintf(os.Stderr, "[dry-run]   Version: latest\n")
	}

	home, _ := os.UserHomeDir()
	globalRepo := fmt.Sprintf("%s/.agents/skills/%s", home, slug)
	fmt.Fprintf(os.Stderr, "[dry-run]   Target: %s\n", globalRepo)

	agents := skillhub.DetectAllInstalledAgents()
	if len(agents) > 0 {
		fmt.Fprintf(os.Stderr, "[dry-run]   Detected agents: ")
		for i, a := range agents {
			if i > 0 {
				fmt.Fprintf(os.Stderr, ", ")
			}
			fmt.Fprintf(os.Stderr, "%s", a.DisplayName)
		}
		fmt.Fprintf(os.Stderr, "\n")

		linkAction := "symlinks"
		if runtime.GOOS == "windows" {
			linkAction = "copies"
		}
		fmt.Fprintf(os.Stderr, "[dry-run]   Would create %s to %d agent(s)\n", linkAction, len(agents))
	}

	fmt.Fprintf(os.Stderr, "[dry-run] No changes made. Remove --dry-run to install.\n")
	return nil
}

var skillsAddCmd = &cobra.Command{
	Use:   "add <slug>",
	Short: "Install a skill",
	Long:  "Install a skill from SkillHub to your local or global skills directory",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		slug := args[0]

		if IsDryRun() {
			return dryRunSkillAdd(slug, addVersion, addGlobal)
		}

		client := skillhub.NewClient()
		return client.Add(slug, addVersion, addGlobal, addYes)
	},
}

func init() {
	skillsAddCmd.Flags().StringVarP(&addVersion, "version", "v", "", "Specific version to install (default: latest)")
	skillsAddCmd.Flags().BoolVarP(&addGlobal, "global", "g", false, "Install globally (auto-detects agent)")
	skillsAddCmd.Flags().BoolVarP(&addYes, "yes", "y", false, "Skip confirmation prompts")
}
