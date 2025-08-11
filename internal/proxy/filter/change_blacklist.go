package filter

import (
	"den-den-mushi-Go/internal/proxy/config"
	"go.uber.org/zap"
	"regexp"
	"sync"
)

type ChangeBlacklistFilter struct {
	mu                     sync.RWMutex
	ouGroupRegexFiltersMap map[string][]regexp.Regexp
	log                    *zap.Logger
	cfg                    *config.Config
}

func (c *ChangeBlacklistFilter) IsValid(cmd string, ouGroup string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.cfg.Filters.IsChangeBlacklistEnabled {
		c.log.Info("Change blacklist is disabled, allowing all commands")
		return cmd, true
	}

	// todo this is repeated everywhere, refactor once stable
	cmd = preprocessCommand(cmd)

	if !isValidForDefault(cmd, true, c.ouGroupRegexFiltersMap) {
		c.log.Debug("Command is not valid for default group", zap.String("cmd", cmd))
		return cmd, false
	}

	if !isValidForOuGroup(cmd, true, ouGroup, c.ouGroupRegexFiltersMap) {
		c.log.Debug("Command is not valid for OU group", zap.String("cmd", cmd), zap.String("ouGroup", ouGroup))
		return cmd, false
	}

	return cmd, true
}

// todo load from db
func (c *ChangeBlacklistFilter) load(_ map[string][]regexp.Regexp) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ouGroupRegexFiltersMap["default"] = []regexp.Regexp{
		*regexp.MustCompile(`(?i)^\s*su\s*$`),         // su
		*regexp.MustCompile(`(?i)^\s*su\s*-$`),        // su -
		*regexp.MustCompile(`(?i)^\s*su\s+-\s*root$`), // su - root
		*regexp.MustCompile(`(?i)^\s*su\s+root$`),     // su root

		*regexp.MustCompile(`(?i)^\s*sudo\s*$`),         // sudo
		*regexp.MustCompile(`(?i)^\s*sudo\s*-$`),        // sudo -
		*regexp.MustCompile(`(?i)^\s*sudo\s+-\s*root$`), // sudo - root
		*regexp.MustCompile(`(?i)^\s*sudo\s+root$`),     // sudo root
	}
}
