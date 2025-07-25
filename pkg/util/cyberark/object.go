package cyberark

import (
	"fmt"
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
		if extractOsUserFromObject(o) == osUser {
			return true
		}
	}
	return false
}

func extractOsUserFromObject(o string) string {
	if o == "" {
		return ""
	}

	// special cases
	if strings.Contains(o, "ec2-user") {
		return "ec2-user"
	} else if strings.Contains(o, "ec2-read") {
		return "ec2-read"
	} else if strings.Contains(o, "ec2-app") {
		return "ec2-app"
	}

	// extract os user, assuming object format xxx-osuser-xxxxx where osuser is the second part
	parts := strings.Split(o, "-")
	if len(parts) < 2 {
		return ""
	}

	return parts[1] // second element
}

func MapIPToOSUsers(objects []string) map[string][]string {
	ipToUsers := make(map[string]map[string]struct{})

	for _, o := range objects {
		ip := extractIPFromCyberarkObject(o)
		user := extractOsUserFromObject(o)
		fmt.Printf("Parsed IP=%s, User=%s, Raw=%s\n", ip, user, o)

		if ip == "" || user == "" {
			continue
		}

		if _, ok := ipToUsers[ip]; !ok {
			ipToUsers[ip] = make(map[string]struct{})
		}
		ipToUsers[ip][user] = struct{}{}
	}

	result := make(map[string][]string)
	for ip, userSet := range ipToUsers {
		for user := range userSet {
			result[ip] = append(result[ip], user)
		}
	}

	return result
}
