package policy

import (
	"den-den-mushi-Go/internal/control/cyberark"
	"go.uber.org/zap"
	"time"
)

func (p *ChangePolicy[T]) isUserInChangeImplementerGroup(user string, implGroup []string) bool {
	if user == "" || implGroup == nil {
		return false
	}

	//groups, err := p.impGroupService.FindAllByUserId(user)
	//if err != nil {
	//	p.log.Error("Failed to find implementer groups for user", zap.String("user", user), zap.Error(err))
	//	return false
	//}

	// todo: optimise
	//for _, g := range groups {
	//	if arrayContains(implGroup, g.Name) {
	//		return true
	//	}
	//}
	return false
}

func (p *ChangePolicy[T]) isServerIpInChangeRequest(ip string, objects []string) bool {
	if ip == "" {
		p.log.Error("Server IP missing", zap.String("ip", ip))
		return false
	}

	return cyberark.IsIpInCyberarkObject(ip, objects)
}

func (p *ChangePolicy[T]) isOsUserInChangeRequest(osUser string, objects []string) bool {
	if osUser == "" {
		p.log.Error("OS User missing", zap.String("osUser", osUser))
		return false
	}

	return cyberark.IsOsUserInCyberarkObject(osUser, objects)
}

func (p *ChangePolicy[T]) isCRApproved(crState string) bool {
	if crState == "" {
		p.log.Error("Change request state missing", zap.String("state", crState))
		return false
	}

	return crState == "approved"
}

func isValidWindow(start, end string) bool {
	if start == "" || end == "" {
		return false
	}
	startTime, err := time.Parse(time.RFC3339, start)
	if err != nil {
		return false
	}
	endTime, err := time.Parse(time.RFC3339, end)
	if err != nil {
		return false
	}

	return time.Now().After(startTime) && time.Now().Before(endTime)
}
