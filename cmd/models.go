package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/VtrixAI/vtrix-cli/internal/models"
)

var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "Browse available models",
}

var (
	modelsListType     string
	modelsListKeywords string
	modelsListPage     int
	modelsListPageSize int
	modelsListOutput   string
)

var modelsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available models",
	Long: `List models available on vtrix.

Output fields (--output json):
  id                 Model identifier, use this as <model_id> in "vtrix models spec <model_id>"
  name               Human-readable model name
  type               Model type: Video | Image | Audio
  description        What the model does
  input_modalities   Accepted input types: Text | Image | Video | Audio
  output_modalities  Output types produced by the model

Pagination fields:
  total              Total number of matching models
  page               Current page
  page_size          Results per page
  total_pages        Total number of pages`,
	Example: `  vtrix models list
  vtrix models list --type video
  vtrix models list --keywords kirin
  vtrix models list --output id
  vtrix models list --output json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := models.List(models.ListParams{
			Page:     modelsListPage,
			PageSize: modelsListPageSize,
			Type:     modelsListType,
			Keywords: modelsListKeywords,
		})
		if err != nil {
			return err
		}

		if modelsListOutput == "json" {
			b, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(b))
			return nil
		}

		if modelsListOutput == "id" {
			for _, m := range result.Models {
				fmt.Println(m.ID)
			}
			return nil
		}

		if len(result.Models) == 0 {
			fmt.Println("No models found.")
			return nil
		}

		fmt.Printf("Showing %d of %d models (page %d/%d)\n\n",
			len(result.Models), result.Total, result.Page, result.TotalPages)

		for _, m := range result.Models {
			fmt.Printf("%-30s  %-8s  %s\n", m.ID, m.Type, m.Name)
			if m.Description != "" {
				fmt.Printf("  %s\n", truncate(m.Description, 80))
			}
			fmt.Printf("  Input: %s  →  Output: %s\n\n",
				strings.Join(m.InputModalities, ", "),
				strings.Join(m.OutputModalities, ", "),
			)
		}
		return nil
	},
}

var modelsSpecOutput string

var modelsSpecCmd = &cobra.Command{
	Use:   "spec <model_id>",
	Short: "Get full parameter spec for a model",
	Long: `Get the complete parameter specification for a model.

Default output is agent_prompt: a formatted text containing:
  - API endpoint, method, and headers
  - Full request body template with all parameters
  - Parameter table (type, required, allowed values, default, description)
  - Async task polling instructions

Use --output json to get the raw structured spec including:
  model_id     Model identifier
  name         Model display name
  vendor       Provider (e.g. kling, vidu)
  type         Model type (video / image / audio)
  api          Endpoint, method, headers template
  parameters   Full parameter definitions with types, constraints, children
  agent_prompt Preformatted text ready to be injected into an LLM context`,
	Example: `  vtrix models spec kirin_v2_6_i2v
  vtrix models spec kirin_v2_6_i2v --output json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		modelID := args[0]

		spec, err := models.GetSpec(modelID)
		if err != nil {
			return err
		}

		if modelsSpecOutput == "json" {
			b, _ := json.MarshalIndent(spec, "", "  ")
			fmt.Println(string(b))
			return nil
		}

		// default: print agent_prompt (最适合 agent 阅读)
		fmt.Println(spec.AgentPrompt)
		return nil
	},
}

func buildModelsQuery() string {
	q := fmt.Sprintf("?page=%d&page_size=%d", modelsListPage, modelsListPageSize)
	if modelsListType != "" {
		q += "&type=" + modelsListType
	}
	if modelsListKeywords != "" {
		q += "&keywords=" + modelsListKeywords
	}
	return q
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func init() {
	modelsListCmd.Flags().StringVar(&modelsListType, "type", "", "Filter by type (video, image, audio)")
	modelsListCmd.Flags().StringVar(&modelsListKeywords, "keywords", "", "Search by keyword")
	modelsListCmd.Flags().IntVar(&modelsListPage, "page", 1, "Page number")
	modelsListCmd.Flags().IntVar(&modelsListPageSize, "page-size", 20, "Results per page")
	modelsListCmd.Flags().StringVar(&modelsListOutput, "output", "", "Output format: id (IDs only), json (full response)")

	modelsSpecCmd.Flags().StringVar(&modelsSpecOutput, "output", "", "Output format (json)")

	modelsCmd.AddCommand(modelsListCmd)
	modelsCmd.AddCommand(modelsSpecCmd)
}
