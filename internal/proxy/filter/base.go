package filter

import (
	"den-den-mushi-Go/pkg/types"
	"fmt"
	"regexp"
)

type CommandFilter interface {
	IsValid(str string, ouGroup string) (string, bool)
	load(patterns map[string][]regexp.Regexp)
}

var (
	whitelistFilter = &WhitelistFilter{
		ouGroupRegexFiltersMap: make(map[string][]regexp.Regexp),
	}

	blacklistFilter = &BlacklistFilter{
		ouGroupRegexFiltersMap: make(map[string][]regexp.Regexp),
	}
)

func GetFilter(filterType types.Filter) CommandFilter {
	switch filterType {
	case types.Whitelist:
		return whitelistFilter
	case types.Blacklist:
		fmt.Println("Debug: Setting up blacklist filter")
		return blacklistFilter
	default:
		return nil
	}
}

//// hard code for now
//func init() {
//	whitelist := []string{
//		"ll",
//		"whoami",
//		"pwd",
//		"ls",
//		"ls -a",
//		"ls -t",
//		"ls -l",
//		"ls -r",
//		"ls -rt",
//		"ls -lrt",
//		"ls -al",
//	}
//
//	blacklist := []string{
//		"rm -rf",
//		"sudo",
//		"sudo su",
//		"su",
//		"su -",
//		"shutdown",
//		"reboot",
//	}
//
//	whitelistFilter.UpdateCommands(whitelist)
//	blacklistFilter.UpdateCommands(blacklist)
//}
