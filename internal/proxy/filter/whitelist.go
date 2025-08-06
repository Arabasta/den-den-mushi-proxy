package filter

import (
	"den-den-mushi-Go/internal/proxy/config"
	"go.uber.org/zap"
	"regexp"
	"sync"
)

type WhitelistFilter struct {
	mu                     sync.RWMutex
	ouGroupRegexFiltersMap map[string][]regexp.Regexp
	log                    *zap.Logger
	cfg                    *config.Config
}

func (w *WhitelistFilter) IsValid(cmd string, ouGroup string) (string, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if !w.cfg.Filters.IsHealthcheckWhitelistEnabled {
		w.log.Info("Healthcheck whitelist is disabled, allowing all commands")
		return cmd, true
	}

	cmd = preprocessCommand(cmd)

	if isValidForDefault(cmd, false, w.ouGroupRegexFiltersMap) {
		w.log.Debug("Command is valid for default group", zap.String("cmd", cmd))
		return cmd, true
	}

	if isValidForOuGroup(cmd, false, ouGroup, w.ouGroupRegexFiltersMap) {
		w.log.Debug("Command is valid for OU group", zap.String("cmd", cmd), zap.String("ouGroup", ouGroup))
		return cmd, true
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
