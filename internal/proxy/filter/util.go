package filter

import (
	"regexp"
	"strings"
)

var allWhiteSpace = regexp.MustCompile(`\s+`)

func preprocessCommand(cmd string) string {
	cmd = strings.TrimSpace(cmd)

	// normalise whitespace
	cmd = allWhiteSpace.ReplaceAllString(cmd, " ")

	// todo use mvdan here
	return cmd
}

func isValidForDefault(cmd string, isBlacklist bool, filterMap map[string][]regexp.Regexp) bool {
	return isValid(cmd, isBlacklist, filterMap["default"])
}

func isValidForOuGroup(cmd string, isBlacklist bool, ouGroup string, filterMap map[string][]regexp.Regexp) bool {
	return isValid(cmd, isBlacklist, filterMap[ouGroup])
}

func isValid(cmd string, isBlacklist bool, filters []regexp.Regexp) bool {
	matched := false
	for _, pattern := range filters {
		if pattern.MatchString(cmd) {
			matched = true
			break
		}
	}
	if isBlacklist {
		return !matched // blacklist, true if not matched
	}
	return matched // whitelist, true if matched
}
