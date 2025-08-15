package meta

import (
	"den-den-mushi-Go/internal/control/core/implementor_groups"
	"den-den-mushi-Go/pkg/middleware"
	"go.uber.org/zap"
)

type Service struct {
	impGroupSvc *implementor_groups.Service
	log         *zap.Logger
}

func NewService(impGroupSvc *implementor_groups.Service, log *zap.Logger) *Service {
	log.Info("Initializing Meta Service...")
	return &Service{
		impGroupSvc: impGroupSvc,
		log:         log,
	}
}

func (s *Service) getUserImplementorGroups(authCtx *middleware.AuthContext) ([]string, error) {
	impGroups, err := s.impGroupSvc.FindAllByUserId(authCtx.UserID)
	if err != nil {
		s.log.Error("Failed to fetch implementor groups", zap.Error(err))
		return nil, err
	}

	userImplGroups := make([]string, 0, len(impGroups))
	for _, group := range impGroups {
		userImplGroups = append(userImplGroups, group.GroupName)
	}

	return userImplGroups, nil
}
