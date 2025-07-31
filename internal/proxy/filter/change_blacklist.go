package filter

import (
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"regexp"
	"sync"
)

type ChangeBlacklistFilter struct {
	mu                     sync.RWMutex
	ouGroupRegexFiltersMap map[string][]regexp.Regexp
}

func (b *ChangeBlacklistFilter) IsValid(cmd string, ouGroup string) (string, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	ouGroup = "default"

	log.Debug("Debug: Checking command for su", zap.String("cmd", cmd))
	blockedCmds, ok := b.ouGroupRegexFiltersMap[ouGroup]
	if !ok {
		log.Debug("No su found")
		return cmd, true
	}

	for _, blocked := range blockedCmds {
		log.Debug("Checking cmd against blacklist", zap.String("cmd", cmd))
		if blocked.MatchString(cmd) {
			return cmd, false // blocked
		}
	}

	return cmd, true
}

func (b *ChangeBlacklistFilter) load(_ map[string][]regexp.Regexp) {
	// gg ISP fail, rush me more
}
