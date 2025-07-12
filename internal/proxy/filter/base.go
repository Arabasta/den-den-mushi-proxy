package filter

import (
	"den-den-mushi-Go/pkg/types"
)

type CommandFilter interface {
	IsValid(str string) (string, bool)
}

var (
	whitelistFilter = &WhitelistFilter{
		filteredCommands: make(map[string]struct{}),
	}
	blacklistFilter = &BlacklistFilter{
		filteredCommands: make(map[string]struct{}),
	}
)

// hard code for now
func init() {
	whitelist := []string{
		"ll",
		"whoami",
		"pwd",
		"ls",
		"ls -a",
		"ls -t",
		"ls -l",
		"ls -r",
		"ls -rt",
		"ls -lrt",
		"ls -al",
	}

	blacklist := []string{
		"rm -rf",
		"sudo",
		"sudo su",
		"su",
		"su -",
		"shutdown",
		"reboot",
	}

	whitelistFilter.UpdateCommands(whitelist)
	blacklistFilter.UpdateCommands(blacklist)
}

func GetFilter(filterType types.Filter) CommandFilter {
	switch filterType {
	case types.Whitelist:
		return whitelistFilter
	case types.Blacklist:
		return blacklistFilter
	default:
		return nil
	}
}
