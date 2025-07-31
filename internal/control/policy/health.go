package policy

import (
	"den-den-mushi-Go/internal/control/config"
	"den-den-mushi-Go/internal/control/host"
	"den-den-mushi-Go/internal/control/implementor_groups"
	"den-den-mushi-Go/internal/control/pty_token/request"
	"den-den-mushi-Go/pkg/types"
	"errors"
	"go.uber.org/zap"
	"strings"
)

type HealthcheckPolicy[T request.Ctx] struct {
	next Policy[T]

	impGroupService *implementor_groups.Service
	hostService     *host.Service
	log             *zap.Logger
	l1OuGroupPrefix string
	cfg             *config.Config
}

func NewHealthcheckPolicy[T request.Ctx](hostService *host.Service, impGroupService *implementor_groups.Service,
	log *zap.Logger, cfg *config.Config) *HealthcheckPolicy[T] {
	return &HealthcheckPolicy[T]{
		hostService:     hostService,
		impGroupService: impGroupService,
		log:             log,
		l1OuGroupPrefix: cfg.OuGroup.Prefix.L1,
		cfg:             cfg,
	}
}

func (p *HealthcheckPolicy[T]) SetNext(n Policy[T]) {
	p.next = n
}

func (p *HealthcheckPolicy[T]) Check(r T) error {
	// 1. skip healthcheck requests
	if r.GetPurpose() != types.Healthcheck {
		if p.next != nil {
			return p.next.Check(r)
		}
		return nil
	}

	// user must be in L1 OU group
	if p.cfg.OuGroup.IsValidationEnabled {
		ouGroup := r.GetUserOuGroup()
		if ouGroup == "" {
			p.log.Warn("User OU group not found", zap.String("userId", r.GetUserId()))
			return errors.New("OU group not found")
		}

		p.log.Debug("Checking user OU group for healthcheck", zap.String("userId", r.GetUserId()), zap.String("ouGroup", ouGroup))
		if strings.HasPrefix(ouGroup, p.l1OuGroupPrefix) == false {
			p.log.Warn("User OU group is not L1", zap.String("userId", r.GetUserId()), zap.String("ouGroup", ouGroup), zap.String("expectedPrefix", p.l1OuGroupPrefix))
			return errors.New("OU group is not L1")
		}

		p.log.Debug("User OU group is L1", zap.String("userId", r.GetUserId()), zap.String("ouGroup", ouGroup))
	}

	// 2. get host type

	// 3. check if user is in host healthcheck group

	// 4. check is os user is a valid readonly user

	if p.next != nil {
		return p.next.Check(r)
	}
	return nil
}
