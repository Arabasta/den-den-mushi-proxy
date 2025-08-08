package whiteblacklist

import (
	"den-den-mushi-Go/internal/control/core/regex_filters"
	dto "den-den-mushi-Go/pkg/dto/regex_filters"
	"den-den-mushi-Go/pkg/middleware"
	"den-den-mushi-Go/pkg/types"
	"den-den-mushi-Go/pkg/util/regex"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Service struct {
	svc *regex_filters.Service
	log *zap.Logger
}

func NewService(svc *regex_filters.Service, log *zap.Logger) *Service {
	log.Info("Initializing White/Blacklist Service")
	return &Service{
		svc: svc,
		log: log,
	}
}

// todo: verify user has permission to CRUD ou group filters

func (s *Service) GetRegexFilters(t types.Filter, c *gin.Context) (*[]dto.Record, error) {
	authCtx, ok := middleware.GetAuthContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth context missing"})
		return nil, errors.New("auth context missing")
	}

	return s.svc.FindAllByFilterTypeAndOuGroup(t, authCtx.OuGroup)
}

func (s *Service) CreateRegex(pattern string, t types.Filter, c *gin.Context) (*dto.Record, error) {
	authCtx, ok := middleware.GetAuthContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth context missing"})
		return nil, errors.New("auth context missing")
	}

	re, err := regex.CompilePattern(pattern)
	if err != nil {
		return nil, err
	}

	filter := &dto.Record{
		RegexPattern: *re,
		FilterType:   t,
		IsEnabled:    true,
		CreatedBy:    authCtx.UserID,
		UpdatedBy:    authCtx.UserID,
		OuGroup:      authCtx.OuGroup,
	}

	f, err := s.svc.Save(filter)
	if err != nil {
		s.log.Error("failed to create filter", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return nil, err
	}

	return f, nil
}

func (s *Service) UpdateRegex(id int, pattern string, isEnabled bool, c *gin.Context) (*dto.Record, error) {
	authCtx, ok := middleware.GetAuthContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth context missing"})
		return nil, errors.New("auth context missing")
	}

	re, err := regex.CompilePattern(pattern)
	if err != nil {
		return nil, err
	}

	filter := &dto.Record{
		Id:           uint(id),
		RegexPattern: *re,
		IsEnabled:    isEnabled,
		UpdatedBy:    authCtx.UserID,
	}

	f, err := s.svc.Update(filter)
	if err != nil {
		s.log.Error("failed to update filter", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return nil, err
	}

	return f, nil
}

func (s *Service) SoftDeleteRegex(id int, c *gin.Context) (*dto.Record, error) {
	authCtx, ok := middleware.GetAuthContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth context missing"})
		return nil, errors.New("auth context missing")
	}

	filter := &dto.Record{
		Id:        uint(id),
		DeletedBy: authCtx.UserID,
	}

	f, err := s.svc.SoftDelete(filter)
	if err != nil {
		s.log.Error("failed to delete filter", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return nil, err
	}

	return f, nil
}
