package util

import (
	"den-den-mushi-Go/pkg/types"
	"errors"
)

func GetChangeRequestIDOrError(p types.ConnectionPurpose, id string) (string, error) {
	switch p {
	case types.Change:
		if id == "" {
			return "", errors.New("missing change request ID for change purpose")
		}
		return id, nil
	case types.Healthcheck:
		return "", nil
	default:
		return "", errors.New("invalid connection purpose: " + string(p))
	}
}
