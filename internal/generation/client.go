package generation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/VtrixAI/vtrix-cli/internal/buildinfo"
)

// BaseURL is the generation gateway base URL used by GetTask.
// Set at build time via ldflags:
//
//	go build -ldflags "-X github.com/VtrixAI/vtrix-cli/internal/generation.BaseURL=https://gateway.vtrix.ai"
//
// Or at runtime via the VTRIX_GENERATION_URL environment variable.
var BaseURL = ""

// Client is a lightweight HTTP client for the vtrix gateway.
// It uses the API key (not the auth JWT) for authorization.
type Client struct {
	httpClient *http.Client
	apiKey     string
}

func NewClient(apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		apiKey:     apiKey,
	}
}

func (c *Client) do(method, endpoint string, body []byte, out any) error {
	var reqBody *bytes.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	} else {
		reqBody = bytes.NewReader([]byte{})
	}
	req, err := http.NewRequest(method, endpoint, reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("User-Agent", buildinfo.UserAgent())
	req.Header.Set("X-Source", "cli")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		var errBody struct {
			Message string `json:"message"`
			Error   string `json:"error"`
		}
		if json.Unmarshal(respBody, &errBody) == nil {
			msg := errBody.Message
			if msg == "" {
				msg = errBody.Error
			}
			if msg != "" {
				return fmt.Errorf("HTTP %d: %s", resp.StatusCode, msg)
			}
		}
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	return json.Unmarshal(respBody, out)
}

type Request struct {
	Model string       `json:"model"`
	Input []InputBlock `json:"input"`
}

type InputBlock struct {
	Params map[string]any `json:"params"`
}

// TaskStatus mirrors the gateway's task response.
type TaskStatus struct {
	ID        string        `json:"id"`
	Status    string        `json:"status"` // in_progress | completed | failed
	Model     string        `json:"model"`
	Output    []OutputGroup `json:"output,omitempty"`
	Error     *string       `json:"error,omitempty"`
	CreatedAt int64         `json:"created_at"`
	Progress  float64       `json:"progress,omitempty"`
	Usage     *UsageInfo    `json:"usage,omitempty"`
}

type OutputGroup struct {
	Content []OutputContent `json:"content"`
}

type OutputContent struct {
	Type     string `json:"type"` // video | image | audio
	URL      string `json:"url"`
	Duration string `json:"duration,omitempty"`
	ID       string `json:"id,omitempty"`
	ImgID    int64  `json:"img_id,omitempty"`
}

type UsageInfo struct {
	Cost      string               `json:"cost"`
	Quantity  int                  `json:"quantity"`
	UnitPrice string               `json:"unit_price"`
	ExtraInfo map[string]any       `json:"extra_info,omitempty"`
}

func (c *Client) Submit(endpoint, modelID string, params map[string]any) (*TaskStatus, error) {
	body, err := json.Marshal(Request{
		Model: modelID,
		Input: []InputBlock{{Params: params}},
	})
	if err != nil {
		return nil, err
	}

	var resp TaskStatus
	if err := c.do(http.MethodPost, endpoint, body, &resp); err != nil {
		return nil, err
	}
	if resp.ID == "" {
		if resp.Error != nil && *resp.Error != "" {
			return nil, fmt.Errorf("%s", *resp.Error)
		}
		return nil, fmt.Errorf("task creation failed: no id in response")
	}
	return &resp, nil
}

func (c *Client) PollTask(generationEndpoint, taskID string, pollInterval, timeout time.Duration, onProgress func(float64)) (*TaskStatus, error) {
	taskEndpoint := taskEndpointFrom(generationEndpoint, taskID)
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		var status TaskStatus
		if err := c.do(http.MethodGet, taskEndpoint, nil, &status); err != nil {
			time.Sleep(pollInterval)
			continue
		}

		if onProgress != nil {
			onProgress(status.Progress)
		}

		switch status.Status {
		case "completed":
			return &status, nil
		case "failed":
			reason := "unknown error"
			if status.Error != nil && *status.Error != "" {
				reason = *status.Error
			}
			return &status, fmt.Errorf("%s", reason)
		}

		time.Sleep(pollInterval)
	}

	return nil, fmt.Errorf("timed out after %s", timeout)
}

func (c *Client) GetTask(taskID string) (*TaskStatus, error) {
	base := BaseURL
	if env := os.Getenv("VTRIX_GENERATION_URL"); env != "" {
		base = env
	}
	if base == "" {
		return nil, fmt.Errorf("generation base URL not configured: set VTRIX_GENERATION_URL or rebuild with -ldflags")
	}
	endpoint := strings.TrimRight(base, "/") + "/model/v1/generation/task/" + taskID
	var status TaskStatus
	if err := c.do(http.MethodGet, endpoint, nil, &status); err != nil {
		return nil, err
	}
	return &status, nil
}

func (t *TaskStatus) URLs() []string {
	var urls []string
	for _, g := range t.Output {
		for _, c := range g.Content {
			if c.URL != "" {
				urls = append(urls, c.URL)
			}
		}
	}
	return urls
}

func taskEndpointFrom(generationEndpoint, taskID string) string {
	u, err := url.Parse(generationEndpoint)
	if err != nil {
		return generationEndpoint + "/task/" + taskID
	}
	u.Path = strings.TrimRight(u.Path, "/") + "/task/" + taskID
	return u.String()
}
