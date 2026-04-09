package models

import "github.com/VtrixAI/vtrix-cli/internal/clierrors"

func List(params ListParams) (*ModelsListResponse, error) {
	result, err := NewClient().List(params)
	if err != nil {
		return nil, clierrors.ErrFetchModels(err)
	}
	return result, nil
}

func GetSpec(modelID string) (*ModelSpec, error) {
	spec, err := NewClient().GetSpec(modelID)
	if err != nil {
		if isNotFound(err) {
			return nil, clierrors.ErrModelNotFound(modelID)
		}
		return nil, clierrors.ErrFetchModelSpec(modelID, err)
	}
	return spec, nil
}

func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return len(s) >= 10 && s[:10] == "status 404"
}
