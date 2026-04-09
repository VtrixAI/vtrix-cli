package generation

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/VtrixAI/vtrix-cli/internal/clierrors"
	"github.com/VtrixAI/vtrix-cli/internal/models"
)

func Submit(apiKey, endpoint, modelID string, params map[string]any) (*TaskStatus, error) {
	return NewClient(apiKey).Submit(endpoint, modelID, params)
}

func PollTask(apiKey, generationEndpoint, taskID string, pollInterval, timeout time.Duration, onProgress func(float64)) (*TaskStatus, error) {
	return NewClient(apiKey).PollTask(generationEndpoint, taskID, pollInterval, timeout, onProgress)
}

func GetTask(apiKey, taskID string) (*TaskStatus, error) {
	return NewClient(apiKey).GetTask(taskID)
}

// ParseParams splits "key=value" pairs into a map.
// Supports dot notation for nested fields: "camera_control.type=simple"
func ParseParams(pairs []string) (map[string]string, error) {
	out := make(map[string]string, len(pairs))
	for _, p := range pairs {
		idx := strings.Index(p, "=")
		if idx < 0 {
			return nil, fmt.Errorf("invalid --param %q: expected key=value format", p)
		}
		out[p[:idx]] = p[idx+1:]
	}
	return out, nil
}

// ValidateAndCoerce checks raw string params against the spec and returns
// a typed map ready to send. Nested object fields (dot notation) are expanded
// into nested maps.
func ValidateAndCoerce(modelID string, raw map[string]string, specParams []models.ModelParam) (map[string]any, error) {
	out := make(map[string]any)

	for _, p := range specParams {
		if p.Type == "array" || strings.HasPrefix(p.Type, "array[") || strings.HasPrefix(p.Type, "array\\[") {
			value, hasVal := raw[p.Name]
			if !hasVal {
				if p.Required {
					return nil, clierrors.ErrMissingParam(modelID, p.Name)
				}
				continue
			}
			var arr []any
			if err := json.Unmarshal([]byte(value), &arr); err != nil {
				return nil, clierrors.ErrInvalidParam(modelID, p.Name,
					fmt.Sprintf("expected a JSON array (e.g. '[\"url1\",\"url2\"]' or '[{\"key\":\"value\"}]'), got: %s", value))
			}
			out[p.Name] = arr
			continue
		}

		if p.Type == "object" && len(p.Children) > 0 {
			childRaw := make(map[string]string)
			for k, v := range raw {
				prefix := p.Name + "."
				if strings.HasPrefix(k, prefix) {
					childRaw[strings.TrimPrefix(k, prefix)] = v
				}
			}
			if len(childRaw) == 0 {
				if p.Required {
					return nil, clierrors.ErrMissingParam(modelID, p.Name)
				}
				continue
			}
			nested, err := ValidateAndCoerce(modelID, childRaw, p.Children)
			if err != nil {
				return nil, err
			}
			out[p.Name] = nested
			continue
		}

		value, hasVal := raw[p.Name]
		if !hasVal {
			if p.Required {
				return nil, clierrors.ErrMissingParam(modelID, p.Name)
			}
			continue
		}

		coerced, err := coerceValue(modelID, p, value)
		if err != nil {
			return nil, err
		}
		out[p.Name] = coerced
	}

	return out, nil
}

func coerceValue(modelID string, p models.ModelParam, raw string) (any, error) {
	c := p.Constraints

	if c != nil && len(c.Enum) > 0 {
		found := false
		for _, allowed := range c.Enum {
			if raw == allowed {
				found = true
				break
			}
		}
		if !found {
			return nil, clierrors.ErrInvalidParam(modelID, p.Name,
				fmt.Sprintf("%q is not allowed. Allowed values: %s", raw, strings.Join(c.Enum, ", ")))
		}
	}

	switch p.Type {
	case "int", "integer":
		n, err := strconv.Atoi(raw)
		if err != nil {
			return nil, clierrors.ErrInvalidParam(modelID, p.Name, fmt.Sprintf("%q is not a valid integer", raw))
		}
		if c != nil {
			if c.Min != nil && float64(n) < *c.Min {
				return nil, clierrors.ErrInvalidParam(modelID, p.Name, fmt.Sprintf("%d is below minimum %g", n, *c.Min))
			}
			if c.Max != nil && float64(n) > *c.Max {
				return nil, clierrors.ErrInvalidParam(modelID, p.Name, fmt.Sprintf("%d exceeds maximum %g", n, *c.Max))
			}
		}
		return n, nil

	case "float", "number":
		f, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return nil, clierrors.ErrInvalidParam(modelID, p.Name, fmt.Sprintf("%q is not a valid number", raw))
		}
		if c != nil {
			if c.Min != nil && f < *c.Min {
				return nil, clierrors.ErrInvalidParam(modelID, p.Name, fmt.Sprintf("%g is below minimum %g", f, *c.Min))
			}
			if c.Max != nil && f > *c.Max {
				return nil, clierrors.ErrInvalidParam(modelID, p.Name, fmt.Sprintf("%g exceeds maximum %g", f, *c.Max))
			}
		}
		return f, nil

	case "boolean", "bool":
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return nil, clierrors.ErrInvalidParam(modelID, p.Name, fmt.Sprintf("%q is not a valid boolean (use true/false)", raw))
		}
		return b, nil

	default:
		if c != nil {
			if c.MinLength != nil && len(raw) < *c.MinLength {
				return nil, clierrors.ErrInvalidParam(modelID, p.Name,
					fmt.Sprintf("value too short (min %d chars)", *c.MinLength))
			}
			if c.MaxLength != nil && len(raw) > *c.MaxLength {
				return nil, clierrors.ErrInvalidParam(modelID, p.Name,
					fmt.Sprintf("value too long (max %d chars)", *c.MaxLength))
			}
		}
		return raw, nil
	}
}
