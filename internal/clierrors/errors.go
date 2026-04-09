package clierrors

import "fmt"

// CLIError is an error with a human-readable hint for what to do next.
// Cobra prints err.Error() directly, so the hint appears inline.
type CLIError struct {
	Message string
	Hint    string
}

func (e *CLIError) Error() string {
	if e.Hint != "" {
		return e.Message + "\n  Hint: " + e.Hint
	}
	return e.Message
}

func ErrNotLoggedIn() error {
	return &CLIError{
		Message: "not logged in",
		Hint:    "Run: vtrix auth login",
	}
}

func ErrTokenExpired() error {
	return &CLIError{
		Message: "session expired",
		Hint:    "Run: vtrix auth login",
	}
}

func ErrTokenInvalid() error {
	return &CLIError{
		Message: "invalid token",
		Hint:    "Run: vtrix auth login",
	}
}

func ErrTokenVerification(err error) error {
	return &CLIError{
		Message: fmt.Sprintf("token verification failed: %v", err),
		Hint:    "Run: vtrix auth login",
	}
}

func ErrSaveConfig(err error) error {
	return &CLIError{
		Message: fmt.Sprintf("failed to save config: %v", err),
		Hint:    "Check write permissions for ~/.config/vtrix/",
	}
}

func ErrLogout(err error) error {
	return &CLIError{
		Message: fmt.Sprintf("failed to clear credentials: %v", err),
		Hint:    "Try deleting ~/.config/vtrix/config.yml manually",
	}
}

func ErrNetwork(err error) error {
	return &CLIError{
		Message: fmt.Sprintf("network error: %v", err),
		Hint:    "Check your network connection and that the vtrix API is reachable",
	}
}

func ErrNetworkTimeout(err error) error {
	return &CLIError{
		Message: fmt.Sprintf("request timed out: %v", err),
		Hint:    "Check your network connection or try again",
	}
}

func ErrModelNotFound(id string) error {
	return &CLIError{
		Message: fmt.Sprintf("model %q not found", id),
		Hint:    "Run: vtrix models list to see available models",
	}
}

func ErrFetchModels(err error) error {
	return &CLIError{
		Message: fmt.Sprintf("failed to fetch models: %v", err),
		Hint:    "Check your network connection and try again",
	}
}

func ErrFetchModelSpec(id string, err error) error {
	return &CLIError{
		Message: fmt.Sprintf("failed to fetch spec for %q: %v", id, err),
		Hint:    "Run: vtrix models list to see available models",
	}
}

func ErrNoAPIKey() error {
	return &CLIError{
		Message: "API key not set",
		Hint:    "Run: vtrix auth login to obtain an API key",
	}
}

func ErrInvalidParam(modelID, name, reason string) error {
	return &CLIError{
		Message: fmt.Sprintf("invalid value for parameter %q: %s", name, reason),
		Hint:    fmt.Sprintf("Run: vtrix models spec %s to see allowed values", modelID),
	}
}

func ErrMissingParam(modelID, name string) error {
	return &CLIError{
		Message: fmt.Sprintf("missing required parameter: %q", name),
		Hint:    fmt.Sprintf("Run: vtrix models spec %s to see required parameters", modelID),
	}
}

func ErrSubmitFailed(err error) error {
	return &CLIError{
		Message: fmt.Sprintf("generation request failed: %v", err),
		Hint:    "Check your API key with: vtrix auth status",
	}
}

func ErrTaskFailed(taskID, reason string) error {
	return &CLIError{
		Message: fmt.Sprintf("task %s failed: %s", taskID, reason),
	}
}

func ErrTaskTimeout(taskID string) error {
	return &CLIError{
		Message: fmt.Sprintf("task %s timed out waiting for result", taskID),
		Hint:    fmt.Sprintf("Run: vtrix task status %s to check later", taskID),
	}
}
