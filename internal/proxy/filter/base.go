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

	changeFilter = &ChangeBlacklistFilter{
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

func GetChangeFilter() CommandFilter {
	return changeFilter
}
