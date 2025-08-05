package os_adm_users

import (
	dto "den-den-mushi-Go/pkg/dto/os_adm_users"
	"go.uber.org/zap"
)

type Service struct {
	repo Repository
	log  *zap.Logger
}

func NewService(r Repository, log *zap.Logger) *Service {
	log.Info("Initialising OS Admin Users Service")
	return &Service{
		repo: r,
		log:  log,
	}
}

func (s *Service) FindAllByUserId(userId string) ([]*dto.Record, error) {
	return s.repo.FindAllByUserId(userId)
}

func (s *Service) GetNonCrOsUsers(userId string) []string {
	records, err := s.FindAllByUserId(userId)
	if err != nil {
		return nil
	}
	osUsers := make([]string, 0, len(records)+1)
	for _, r := range records {
		if r.OsUser != "" {
			osUsers = append(osUsers, r.OsUser)
		}
	}

	osUsers = append(osUsers, userId)

	return osUsers
}
