package cyberark

import (
	"net"
	"strings"
)

func IsIpInCyberarkObjects(ip string, objects []string) bool {
	for _, o := range objects {
		if extractIPFromCyberarkObject(o) == ip {
			return true
		}
	}
	return false
}

// extractIPFromCyberarkObject is required as cyberark team doesn't want to give us api endpoint
func extractIPFromCyberarkObject(o string) string {
	if o == "" {
		return ""
	}

	// extract ip, assuming object format 127.0.1-xxxxx
	parts := strings.SplitN(o, "-", 2)
	ip := parts[0]
	if net.ParseIP(ip) == nil {
		return ""
	}

	return ip
}

func IsOsUserInCyberarkObjects(osUser string, objects []string) bool {
	for _, o := range objects {
		if extractOsUsersFromObject(o) == osUser {
			return true
		}
	}
	return false
}

func extractOsUsersFromObject(o string) string {
	if o == "" {
		return ""
	}

	// special case for ec2-user
	if strings.Contains(o, "ec2-user") {
		return "ec2-user"
	}

	// extract os user, assuming object format xxx-osuser-xxxxx where osuser is the second last part
	parts := strings.Split(o, "-")
	if len(parts) < 3 {
		return ""
	}

	return parts[len(parts)-2] // second last element
}
