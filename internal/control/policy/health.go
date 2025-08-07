package policy

import (
	"den-den-mushi-Go/internal/control/ep/pty_token/request"
	"den-den-mushi-Go/internal/control/host"
	"den-den-mushi-Go/internal/control/implementor_groups"
	"den-den-mushi-Go/internal/control/os_adm_users"
	"den-den-mushi-Go/internal/control/policy/validators"
	"den-den-mushi-Go/pkg/types"
	"errors"
	"go.uber.org/zap"
)

type HealthcheckPolicy[T request.Ctx] struct {
	next Policy[T]

	impGroupService *implementor_groups.Service
	hostService     *host.Service
	osAdmUsersSvc   *os_adm_users.Service

	log *zap.Logger
	v   *validators.Validator
}

func NewHealthcheckPolicy[T request.Ctx](hostService *host.Service, impGroupService *implementor_groups.Service, osAdmUsersSvc *os_adm_users.Service,
	log *zap.Logger, v *validators.Validator) *HealthcheckPolicy[T] {
	return &HealthcheckPolicy[T]{
		hostService:     hostService,
		impGroupService: impGroupService,
		osAdmUsersSvc:   osAdmUsersSvc,
		log:             log,
		v:               v,
	}
}

func (p *HealthcheckPolicy[T]) SetNext(n Policy[T]) {
	p.next = n
}

func (p *HealthcheckPolicy[T]) Check(r T) error {
	p.log.Debug("Checking healthcheck policy...")

	// skip healthcheck requests
	if r.GetPurpose() != types.Healthcheck {
		if p.next != nil {
			return p.next.Check(r)
		}
		return nil
	}

	// 2. get host type todo how to map?
	// hostType, err := p.hostService.GetHostType(r.GetServerInfo().IP)

	// 3. check if user is in host healthcheck group todo uncomment this when host type mapping is up
	//if err := p.v.IsOuGroupTypeMatch(r.GetUserOuGroup(), hostType); err != nil {
	//	p.log.Warn("User OU group type validation failed", zap.String("userId", r.GetUserId()), zap.Error(err))
	//	return err
	//}

	// 4. check is os user is a valid readonly user for that user? idk lol
	osAdmUsers := p.osAdmUsersSvc.GetNonCrOsUsers(r.GetUserId())

	if !p.v.IsOsUserInOsAdmUsers(r.GetServerInfo().OSUser, osAdmUsers) {
		p.log.Warn("OS User is not valid", zap.String("osUser", r.GetServerInfo().OSUser), zap.Strings("expected os users", osAdmUsers))
		return errors.New("OS User is not in OS Admin Users or UserId user")
	}

	if p.next != nil {
		return p.next.Check(r)
	}
	return nil
}
