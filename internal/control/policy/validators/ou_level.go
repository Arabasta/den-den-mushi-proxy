package validators

import (
	"errors"
	"go.uber.org/zap"
	"strings"
)

var (
	ErrMissingOuGroup = errors.New("missing OU group")
	ErrInvalidOuGroup = errors.New("invalid OU group")
)

func (v *Validator) IsL1OuGroup(ouGroup string) error {
	return v.isValidOuGroupLevel(ouGroup, v.cfg.OuGroup.Prefix.L1)
}

func (v *Validator) IsL2L3OuGroup(ouGroup string) error {
	return v.isValidOuGroupLevel(ouGroup, v.cfg.OuGroup.Prefix.L2L3)
}

func (v *Validator) isValidOuGroupLevel(ouGroup, expectedOuGroupPrefix string) error {
	if !v.cfg.OuGroup.IsValidationEnabled {
		v.log.Debug("OU group validation is disabled")
		return nil
	}

	if ouGroup == "" {
		v.log.Warn("OU group is empty")
		return ErrMissingOuGroup
	}

	if !strings.HasPrefix(ouGroup, expectedOuGroupPrefix) {
		v.log.Warn("OU group level does not match expected prefix",
			zap.String("expected", expectedOuGroupPrefix),
			zap.String("actual", ouGroup))
		return ErrInvalidOuGroup
	}
	return nil
}
