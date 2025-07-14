package implementor_groups

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

func (s *Service) FindAllByUserId(userId string) ([]*Entity, error) {
	groups, err := s.repository.FindAllByUserId(userId)
	if err != nil {
		s.log.Error("Failed to find groups by user ID", zap.String("userId", userId), zap.Error(err))
		return nil, err
	}
	return groups, nil
}

func (s *Service) IsMemberOfGroup(groupName, userId string) (bool, error) {
	g, err := s.repository.FindByGroupName(groupName)
	if err != nil {
		return false, err
	}
	if g == nil {
		return false, nil
	}

	for _, m := range g.Members {
		if m == userId {
			return true, nil
		}
	}
	return false, nil
}
