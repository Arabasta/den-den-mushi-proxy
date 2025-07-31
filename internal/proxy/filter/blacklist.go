package filter

import (
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"regexp"
	"sync"
)

type BlacklistFilter struct {
	mu                     sync.RWMutex
	ouGroupRegexFiltersMap map[string][]regexp.Regexp
}

func (b *BlacklistFilter) IsValid(cmd string, ouGroup string) (string, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	ouGroup = "default"

	log.Debug("Debug: Checking command against blacklist", zap.String("cmd", cmd), zap.String("ouGroup", ouGroup))
	blockedCmds, ok := b.ouGroupRegexFiltersMap[ouGroup]
	if !ok {
		log.Debug("No blacklist found for OU group", zap.String("ouGroup", ouGroup))
		return cmd, true // no filters for this OU group
	}

	for _, blocked := range blockedCmds {
		log.Debug("Checking cmd against blacklist", zap.String("cmd", cmd), zap.String("ouGroup", ouGroup))
		if blocked.MatchString(cmd) {
			return cmd, false // blocked
		}
	}

	return cmd, true
}

func (b *BlacklistFilter) load(patterns map[string][]regexp.Regexp) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// clear old filters
	b.ouGroupRegexFiltersMap = make(map[string][]regexp.Regexp)

	for ouGroup, regexList := range patterns {
		b.ouGroupRegexFiltersMap[ouGroup] = regexList
	}
}
