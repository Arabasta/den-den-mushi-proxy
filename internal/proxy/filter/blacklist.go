package filter

import (
	"den-den-mushi-Go/internal/proxy/config"
	"go.uber.org/zap"
	"regexp"
	"sync"
)

type BlacklistFilter struct {
	mu                     sync.RWMutex
	ouGroupRegexFiltersMap map[string][]regexp.Regexp
	log                    *zap.Logger
	cfg                    *config.Config
}

func (b *BlacklistFilter) IsValid(cmd string, ouGroup string) (string, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if !b.cfg.Filters.IsHealthcheckBlacklistEnabled {
		b.log.Info("Healthcheck blacklist is disabled, allowing all commands")
		return cmd, true
	}

	cmd = preprocessCommand(cmd)

	if !isValidForDefault(cmd, true, b.ouGroupRegexFiltersMap) {
		b.log.Debug("Command is not valid for default group", zap.String("cmd", cmd))
		return cmd, false
	}

	if !isValidForOuGroup(cmd, true, ouGroup, b.ouGroupRegexFiltersMap) {
		b.log.Debug("Command is not valid for OU group", zap.String("cmd", cmd), zap.String("ouGroup", ouGroup))
		return cmd, false
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
