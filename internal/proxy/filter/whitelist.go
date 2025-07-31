package filter

import (
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"regexp"
	"strings"
	"sync"
)

type WhitelistFilter struct {
	mu                     sync.RWMutex
	ouGroupRegexFiltersMap map[string][]regexp.Regexp
}

func (w *WhitelistFilter) IsValid(cmd string, ouGroup string) (string, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	ouGroup = "default"

	allowedCmds, ok := w.ouGroupRegexFiltersMap[ouGroup]
	if !ok {
		return cmd, false // no filters for this OU group
	}

	cmd = strings.TrimSpace(cmd)

	for _, allowed := range allowedCmds {
		log.Debug("Checking cmd against blacklist", zap.String("cmd", cmd), zap.String("ouGroup", ouGroup))
		if allowed.MatchString(cmd) {
			return cmd, true // allowed
		}
	}

	return cmd, false
}

func (w *WhitelistFilter) load(patterns map[string][]regexp.Regexp) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.ouGroupRegexFiltersMap = make(map[string][]regexp.Regexp)

	for ouGroup, regexList := range patterns {
		w.ouGroupRegexFiltersMap[ouGroup] = regexList
	}
}

//
//func (w *WhitelistFilter) UpdateCommands(newAllowed []string) {
//	w.mu.Lock()
//	defer w.mu.Unlock()
//
//	w.filteredCommands = make(map[string]struct{})
//	for _, cmd := range newAllowed {
//		w.filteredCommands[cmd] = struct{}{}
//	}
//}
