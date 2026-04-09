package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/VtrixAI/vtrix-cli/internal/clierrors"
	"github.com/VtrixAI/vtrix-cli/internal/config"
	"github.com/VtrixAI/vtrix-cli/internal/generation"
	"github.com/VtrixAI/vtrix-cli/internal/models"
)

var (
	runParams  []string
	runOutput  string
	runTimeout int
)

var runCmd = &cobra.Command{
	Use:   "run <model_id>",
	Short: "Run a model and wait for the result",
	Long: `Submit a generation request and poll until the output is ready.

Parameters are passed as --param key=value pairs (repeatable).
Nested object fields use dot notation: --param camera_control.type=simple
Array fields use a JSON string: --param content='[{"type":"text","text":"hello"}]'

Values are coerced to the type declared in the model spec
(string / int / float / boolean / array). Enum and range constraints
are validated before the request is sent.

Exit codes:
  0   task succeeded
  1   error (validation, network, API, timeout)`,
	Example: `  vtrix run kirin_v2_6_i2v --param image=https://example.com/cat.jpg
  vtrix run kirin_v2_6_i2v --param prompt="a cat running" --param duration=5
  vtrix run kirin_v2_6_i2v --param mode=pro --output url
  vtrix run kirin_v2_6_i2v --param mode=pro --output json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		modelID := args[0]

		// dry-run first — no credentials needed
		if IsDryRun() {
			fmt.Fprintf(os.Stderr, "[dry-run] Would execute: POST <spec.api.endpoint>\n")
			fmt.Fprintf(os.Stderr, "[dry-run] model=%s params=%v\n", modelID, runParams)
			fmt.Fprintf(os.Stderr, "[dry-run] Fetch spec first with: vtrix models spec %s\n", modelID)
			return nil
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if cfg.APIKey == "" {
			return clierrors.ErrNoAPIKey()
		}

		// Fetch spec for param validation and endpoint
		spec, err := models.GetSpec(modelID)
		if err != nil {
			return err
		}

		// Parse --param key=value flags
		raw, err := generation.ParseParams(runParams)
		if err != nil {
			return err
		}

		// Validate and coerce params against spec
		params, err := generation.ValidateAndCoerce(modelID, raw, spec.Parameters)
		if err != nil {
			return err
		}

		// Submit generation request
		resp, err := generation.Submit(cfg.APIKey, spec.API.Endpoint, modelID, params)
		if err != nil {
			return clierrors.ErrSubmitFailed(err)
		}

		fmt.Fprintf(os.Stderr, "Task submitted: %s\nWaiting for result...\n", resp.ID)

		// Always poll — POST status reflects submission, not generation completion.
		timeout := time.Duration(runTimeout) * time.Second
		lastProgress := -1.0
		task, pollErr := generation.PollTask(cfg.APIKey, spec.API.Endpoint, resp.ID, 5*time.Second, timeout,
			func(progress float64) {
				pct := int(progress * 100)
				// Only print when progress advances by at least 5%
				if float64(pct)-lastProgress >= 5 || (lastProgress < 0 && pct == 0) {
					fmt.Fprintf(os.Stderr, "Progress: %d%%\n", pct)
					lastProgress = float64(pct)
				}
			},
		)
		if pollErr != nil {
			if task != nil && task.Status == "failed" {
				errMsg := "unknown error"
				if task.Error != nil {
					errMsg = *task.Error
				}
				return clierrors.ErrTaskFailed(resp.ID, errMsg)
			}
			return clierrors.ErrTaskTimeout(resp.ID)
		}

		if task.Status == "failed" {
			errMsg := "unknown error"
			if task.Error != nil {
				errMsg = *task.Error
			}
			return clierrors.ErrTaskFailed(resp.ID, errMsg)
		}

		if runOutput == "url" {
			for _, u := range task.URLs() {
				fmt.Println(u)
			}
			return nil
		}

		if runOutput == "json" {
			b, _ := json.MarshalIndent(task, "", "  ")
			fmt.Println(string(b))
			return nil
		}

		// Human-readable output
		fmt.Printf("Status: %s\n", task.Status)
		for _, g := range task.Output {
			for _, c := range g.Content {
				if c.URL != "" {
					fmt.Printf("URL: %s\n", c.URL)
				}
				if c.ImgID != 0 {
					fmt.Printf("ImgID: %d\n", c.ImgID)
				}
			}
		}
		return nil
	},
}

func init() {
	runCmd.Flags().StringArrayVar(&runParams, "param", nil, "Parameter as key=value (repeatable)")
	runCmd.Flags().StringVar(&runOutput, "output", "", "Output format: url (URLs only), json (full response)")
	runCmd.Flags().IntVar(&runTimeout, "timeout", 600, "Maximum seconds to wait for result (default 10 minutes)")
}
