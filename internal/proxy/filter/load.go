package filter

import (
	"den-den-mushi-Go/pkg/types"
	"go.uber.org/zap"
	"regexp"
	"time"
)

// todo: refactor, temporary solution
// todo pubsub
func (f *Service) startScheduler() {
	f.log.Info("Starting regex filters reload scheduler")
	go func() {
		ticker := time.NewTicker(f.cfg.Filters.DbPollIntervalSeconds * time.Second)
		for range ticker.C {
			if err := f.loadFilters(); err != nil {
				f.log.Error("Failed to reload filters", zap.Error(err))
			}
		}
	}()
}

func (f *Service) loadFilters() error {
	if err := f.loadFilter(types.Whitelist); err != nil {
		return err
	}
	if err := f.loadFilter(types.Blacklist); err != nil {
		return err
	}
	return nil
}

func (f *Service) loadFilter(filterType types.Filter) error {
	filters, err := f.regexFiltersSvc.FindAllEnabledByFilterType(filterType)
	if err != nil {
		f.log.Error("Failed to load filters", zap.Error(err))
		return err
	}
	if filters == nil {
		f.log.Info("No filters found for type", zap.String("filterType", string(filterType)))
		f.GetFilter(filterType, types.Healthcheck).load(nil) // load empty filters for healthcheck
		return nil
	}

	groupMap := make(map[string][]regexp.Regexp)
	for _, f := range *filters {
		groupMap[f.OuGroup] = append(groupMap[f.OuGroup], f.RegexPattern)
	}

	f.GetFilter(filterType, types.Healthcheck).load(groupMap)
	return nil
}
