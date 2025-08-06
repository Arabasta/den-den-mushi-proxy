package validators

import (
	"den-den-mushi-Go/pkg/util/cyberark"
	"time"
)

func (v *Validator) IsServerIpInObjects(ip string, objects []string) bool {
	if ip == "" {
		return false
	}

	return cyberark.IsIpInObjects(ip, objects)
}

func (v *Validator) IsApproved(crState string) bool {
	return crState == "Approved"
}

func (v *Validator) IsValidWindow(start, end time.Time) bool {
	now := time.Now()
	return now.After(start) && now.Before(end)
}
