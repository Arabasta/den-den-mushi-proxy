package validators

import (
	"errors"
	"go.uber.org/zap"
	"strings"
)

func (v *Validator) HasOuGroup(ouGroup string) error {
	if !v.cfg.OuGroup.IsValidationEnabled {
		v.log.Debug("OU group validation is disabled")
		return nil
	}

	if ouGroup == "" {
		v.log.Warn("OU group is empty")
		return errors.New("not part of any OU group")
	}

	return nil
}

func (v *Validator) IsOuGroupTypeMatch(ouGroup string, ouType string) error {
	return v.isValidOuGroupType(ouGroup, ouType)
}

func (v *Validator) isValidOuGroupType(ouGroup, expectedOuGroupSuffix string) error {
	if !v.cfg.OuGroup.IsValidationEnabled {
		v.log.Debug("OU group validation is disabled")
		return nil
	}

	if ouGroup == "" {
		v.log.Warn("OU group is empty")
		return ErrMissingOuGroup
	}

	if !strings.HasSuffix(ouGroup, expectedOuGroupSuffix) {
		v.log.Warn("OU group type does not match expected suffix",
			zap.String("expected", expectedOuGroupSuffix),
			zap.String("actual", ouGroup))
		return ErrInvalidOuGroup
	}
	return nil
}
