package filter

import (
	"den-den-mushi-Go/internal/proxy/regex_filters"
	"den-den-mushi-Go/pkg/types"
	"fmt"
	"go.uber.org/zap"
	"regexp"
	"time"
)

// todo: refactor, temporary solution

type LoadService struct {
	svc *regex_filters.Service
	log *zap.Logger
}

func NewLoadService(svc *regex_filters.Service, log *zap.Logger) *LoadService {
	log.Info("Initializing LoadService for regex filters")
	return &LoadService{
		svc: svc,
		log: log,
	}
}

func (s *LoadService) StartScheduler() {
	s.log.Info("Starting regex filters reload scheduler")
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		for range ticker.C {
			if err := s.loadFilters(); err != nil {
				s.log.Error("Failed to reload filters", zap.Error(err))
			}
		}
	}()
}

func (s *LoadService) loadFilters() error {
	if err := s.loadFilter(types.Whitelist); err != nil {
		return err
	}
	if err := s.loadFilter(types.Blacklist); err != nil {
		return err
	}
	return nil
}

func (s *LoadService) loadFilter(filterType types.Filter) error {
	filters, err := s.svc.FindAllEnabledByFilterType(filterType)
	if err != nil {
		return err
	}
	if filters == nil {
		return nil
	}

	groupMap := make(map[string][]regexp.Regexp)
	for _, f := range *filters {
		groupMap[f.OuGroup] = append(groupMap[f.OuGroup], f.RegexPattern)
	}

	filter := GetFilter(filterType)
	if filter == nil {
		return fmt.Errorf("unknown filter type: %v", filterType)
	}

	filter.load(groupMap)
	//s.log.Debug("Loaded regex filters", zap.String("filterType", string(filterType)),
	//	zap.Int("groupCount", len(groupMap)))
	//s.log.Debug("Loaded regex filters", zap.Any("groupMap", groupMap))
	return nil
}
