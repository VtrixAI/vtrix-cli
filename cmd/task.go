package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/VtrixAI/vtrix-cli/internal/clierrors"
	"github.com/VtrixAI/vtrix-cli/internal/config"
	"github.com/VtrixAI/vtrix-cli/internal/generation"
)

var taskStatusOutput string

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage generation tasks",
}

var taskStatusCmd = &cobra.Command{
	Use:   "status <task_id>",
	Short: "Get the current status of a generation task",
	Long: `Fetch the current status of a generation task by its ID.

Exit codes:
  0   request succeeded (task may still be in_progress)
  1   error (network, API, missing credentials)`,
	Example: `  vtrix task status d758n65e878c73cmdg20
  vtrix task status d758n65e878c73cmdg20 --output url
  vtrix task status d758n65e878c73cmdg20 --output json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID := args[0]

		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if cfg.APIKey == "" {
			return clierrors.ErrNoAPIKey()
		}

		task, err := generation.GetTask(cfg.APIKey, taskID)
		if err != nil {
			return fmt.Errorf("failed to fetch task %s: %w", taskID, err)
		}

		if taskStatusOutput == "url" {
			for _, u := range task.URLs() {
				fmt.Println(u)
			}
			return nil
		}

		if taskStatusOutput == "json" {
			b, _ := json.MarshalIndent(task, "", "  ")
			fmt.Println(string(b))
			return nil
		}

		// Human-readable
		fmt.Printf("Task:   %s\n", task.ID)
		fmt.Printf("Status: %s\n", task.Status)
		if task.Status == "failed" && task.Error != nil {
			fmt.Printf("Error:  %s\n", *task.Error)
		}
		for _, u := range task.URLs() {
			fmt.Printf("URL:    %s\n", u)
		}
		return nil
	},
}

func init() {
	taskStatusCmd.Flags().StringVar(&taskStatusOutput, "output", "", "Output format: url (URLs only), json (full response)")
	taskCmd.AddCommand(taskStatusCmd)
}
