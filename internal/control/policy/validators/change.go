package validators

import (
	"den-den-mushi-Go/pkg/util/cyberark"
	"time"
)

// todo: refactor absolute garbage

func IsServerIpInObjects(ip string, objects []string) bool {
	if ip == "" {
		return false
	}

	return cyberark.IsIpInObjects(ip, objects)
}

func IsApproved(crState string) bool {
	return crState == "Approved"
}

func IsValidWindow(start, end time.Time) bool {
	now := time.Now()
	return now.After(start) && now.Before(end)
}
