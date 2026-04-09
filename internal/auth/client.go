package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/VtrixAI/vtrix-cli/internal/buildinfo"
	"github.com/VtrixAI/vtrix-cli/internal/clierrors"
)

// BaseURL can be overridden at build time via ldflags:
//
//	go build -ldflags "-X github.com/VtrixAI/vtrix-cli/internal/auth.BaseURL=https://vtrix.ai"
//
// Or at runtime via the VTRIX_BASE_URL environment variable.
var BaseURL = ""

const AppID = "@seacloud/web"

type Client struct {
	httpClient *http.Client
	token      string
	baseURL    string
}

func NewClient(token string) *Client {
	base := BaseURL
	if env := os.Getenv("VTRIX_BASE_URL"); env != "" {
		base = env
	}
	return &Client{
		httpClient: &http.Client{Timeout: 15 * time.Second},
		token:      token,
		baseURL:    base,
	}
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", buildinfo.UserAgent())
	req.Header.Set("X-Source", "cli")
	req.Header.Set("X-App-Id", AppID)
	req.Header.Set("X-Version", buildinfo.Version)
	req.Header.Set("X-Plat", "cli")
	req.Header.Set("X-Device-Type", "cli")
	req.Header.Set("X-Skip-Nextauth", "true")

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
		req.Header.Set("X-Auth-Priority", "auth_token")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil, clierrors.ErrNetworkTimeout(err)
		}
		return nil, clierrors.ErrNetwork(err)
	}
	return resp, nil
}

// apiResponse is the common envelope: {"data": ..., "status": {"code": ..., "message": ...}}
type apiResponse struct {
	Data   json.RawMessage `json:"data"`
	Status struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"status"`
}

func (c *Client) post(path string, body []byte, out any) error {
	req, err := http.NewRequest(http.MethodPost, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var envelope apiResponse
	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return fmt.Errorf("unexpected response: %s", string(respBody))
	}

	if envelope.Status.Code != 0 && envelope.Status.Code != 200 {
		switch envelope.Status.Code {
		case 401:
			return clierrors.ErrTokenExpired()
		case 403:
			return clierrors.ErrTokenInvalid()
		}
		return fmt.Errorf("%s", envelope.Status.Message)
	}

	if envelope.Status.Code == 0 && envelope.Data == nil {
		return fmt.Errorf("unexpected response: %s", string(respBody))
	}

	if out != nil && envelope.Data != nil {
		return json.Unmarshal(envelope.Data, out)
	}
	return nil
}

func (c *Client) get(path string, out any) error {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var envelope apiResponse
	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return json.Unmarshal(respBody, out)
	}

	if envelope.Status.Code != 0 && envelope.Status.Code != 200 {
		switch envelope.Status.Code {
		case 401:
			return clierrors.ErrTokenExpired()
		case 403:
			return clierrors.ErrTokenInvalid()
		}
		return fmt.Errorf("%s", envelope.Status.Message)
	}

	if envelope.Data != nil {
		return json.Unmarshal(envelope.Data, out)
	}
	return json.Unmarshal(respBody, out)
}

// MeResponse mirrors the data.user field from /api/v1/auth/me
type MeResponse struct {
	UserID  string `json:"id"`
	Email   string `json:"email"`
	Account string `json:"account"`
	Name    string `json:"name"`
	Role    string `json:"role"`
}

type meData struct {
	User *MeResponse `json:"user"`
}

func (c *Client) Me() (*MeResponse, error) {
	var data meData
	if err := c.get("/api/v1/auth/me", &data); err != nil {
		return nil, err
	}
	if data.User == nil {
		return nil, fmt.Errorf("user not found in response")
	}
	return data.User, nil
}

// DeviceCodeRequest is the request body for POST /api/v1/cli/device/code
type DeviceCodeRequest struct {
	ClientID        string `json:"client_id"`
	ClientPublicKey string `json:"client_public_key"`
}

// DeviceCodeResponse is the response from POST /api/v1/cli/device/code
type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

func (c *Client) RequestDeviceCode(req DeviceCodeRequest) (*DeviceCodeResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	var resp DeviceCodeResponse
	if err := c.post("/api/v1/cli/device/code", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// TokenRequest is the request body for POST /api/v1/cli/token
type TokenRequest struct {
	DeviceCode string `json:"device_code"`
	Timestamp  string `json:"timestamp"`
	Nonce      string `json:"nonce"`
	Proof      string `json:"proof"`
}

// TokenResponse is the response from POST /api/v1/cli/token
type TokenResponse struct {
	Status       string `json:"status"`        // "pending" | "expired"
	AccessToken  string `json:"access_token"`  // set when authorized
	RefreshToken string `json:"refresh_token"` // set when authorized
	APIKey       string `json:"api_key"`       // set when authorized
}

func (c *Client) PollToken(req TokenRequest) (*TokenResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	var resp TokenResponse
	if err := c.post("/api/v1/cli/token", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
