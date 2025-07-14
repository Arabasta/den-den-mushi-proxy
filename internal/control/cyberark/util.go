package cyberark

import (
	"net"
	"strings"
)

func IsIpInCyberarkObject(ip string, objects []string) bool {
	for _, o := range objects {
		if o == "" {
			continue
		}

		extractedIP := extractIPFromCyberarkObject(o)
		if extractedIP == "" {
			continue
		}

		if extractedIP == ip {
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
	if len(parts) == 0 {
		return ""
	}

	ip := parts[0]
	if net.ParseIP(ip) == nil {
		return ""
	}

	return ip
}

func IsOsUserInCyberarkObject(osUser string, objects []string) bool {
	for _, o := range objects {
		if o == "" {
			continue
		}

		extractedOsUser := extractOsUsersFromObjects(o)
		if extractedOsUser == "" {
			continue
		}

		if extractedOsUser == osUser {
			return true
		}
	}
	return false
}

func extractOsUsersFromObjects(o string) string {
	return "" // todo
}
