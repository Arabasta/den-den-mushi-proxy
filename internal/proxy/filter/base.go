package filter

import (
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/regex_filters"
	"den-den-mushi-Go/pkg/types"
	"go.uber.org/zap"
	"regexp"
)

type CommandFilter interface {
	IsValid(str string, ouGroup string) (string, bool)
	load(patterns map[string][]regexp.Regexp)
}

// todo: refactor
type Service struct {
	Whitelist       CommandFilter
	Blacklist       CommandFilter
	ChangeBlacklist CommandFilter
	regexFiltersSvc *regex_filters.Service
	log             *zap.Logger
	cfg             *config.Config
}

func NewService(regexFiltersSvc *regex_filters.Service, log *zap.Logger, cfg *config.Config) *Service {
	whitelist := &WhitelistFilter{ouGroupRegexFiltersMap: make(map[string][]regexp.Regexp), log: log, cfg: cfg}
	blacklist := &BlacklistFilter{ouGroupRegexFiltersMap: make(map[string][]regexp.Regexp), log: log, cfg: cfg}
	changeBlacklist := &ChangeBlacklistFilter{ouGroupRegexFiltersMap: make(map[string][]regexp.Regexp), log: log, cfg: cfg}

	filters := &Service{
		regexFiltersSvc: regexFiltersSvc,
		Whitelist:       whitelist,
		Blacklist:       blacklist,
		ChangeBlacklist: changeBlacklist,
		log:             log,
		cfg:             cfg,
	}

	changeBlacklist.load(nil)
	if err := filters.loadFilters(); err != nil {
		log.Error("Failed to load initial filters", zap.Error(err))
	}

	filters.startScheduler()

	return filters
}

func (f *Service) GetFilter(filterType types.Filter, purpose types.ConnectionPurpose) CommandFilter {
	if purpose == types.Change {
		return f.ChangeBlacklist
	}

	if purpose == types.Healthcheck {
		switch filterType {
		case types.Whitelist:
			return f.Whitelist
		case types.Blacklist:
			return f.Blacklist
		default:
			return nil
		}
	}

	if purpose == types.IExpress {
		switch filterType {
		case types.Whitelist:
			return f.Whitelist
		case types.Blacklist:
			return f.Blacklist
		default:
			return nil
		}
	}
	return nil
}
