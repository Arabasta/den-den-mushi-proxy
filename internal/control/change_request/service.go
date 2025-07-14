package change_request

import "go.uber.org/zap"

type Service struct {
	repository Repository
	log        *zap.Logger
}

func NewService(repository Repository, log *zap.Logger) *Service {
	return &Service{
		repository: repository,
		log:        log,
	}
}

func (s *Service) FindById(id string) (*Entity, error) {
	return s.repository.FindById(id)
}
