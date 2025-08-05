package policy

import (
	"den-den-mushi-Go/internal/control/implementor_groups"
	"den-den-mushi-Go/internal/control/os_adm_users"
	"den-den-mushi-Go/internal/control/policy/validators"
	"den-den-mushi-Go/internal/control/pty_token/request"
	"den-den-mushi-Go/pkg/types"
	"errors"
	"go.uber.org/zap"
)

type ChangePolicy[T request.Ctx] struct {
	next Policy[T]

	impGroupService *implementor_groups.Service
	osAdmUsersSvc   *os_adm_users.Service

	log *zap.Logger
}

func NewChangePolicy[T request.Ctx](impGroupSvc *implementor_groups.Service, osAdmUsersSvc *os_adm_users.Service, log *zap.Logger) *ChangePolicy[T] {
	return &ChangePolicy[T]{
		impGroupService: impGroupSvc,
		osAdmUsersSvc:   osAdmUsersSvc,
		log:             log,
	}
}

func (p *ChangePolicy[T]) SetNext(n Policy[T]) {
	p.next = n
}

func (p *ChangePolicy[T]) Check(r T) error {
	p.log.Debug("Checking change policy...")
	// 1. skip non-change requests
	if r.GetPurpose() != types.Change {
		p.log.Debug("Skipping non-change request", zap.String("purpose", string(r.GetPurpose())))
		if p.next != nil {
			return p.next.Check(r)
		}
		return nil
	}

	// 3a. get Change Request from ctx
	cr := r.GetChangeRequest()
	if cr == nil {
		p.log.Warn("Change request not found", zap.String("changeID", r.GetChangeId()))
		return errors.New("change request not found")
	}

	// 3b. get user's implementor groups  todo: get from Request context and move this to implementor group policy
	impGroups, err := p.impGroupService.FindAllByUserId(r.GetUserId())
	if err != nil {
		p.log.Warn("Failed to find implementor groups for user", zap.String("userId", r.GetUserId()), zap.Error(err))
		return err
	}

	// 4. check if CR valid
	if !validators.IsValidWindow(*cr.ChangeStartTime, *cr.ChangeEndTime) {
		p.log.Warn("Change request time invalid", zap.String("changeID", r.GetChangeId()))
		return errors.New("change request window is invalid")
	}

	if !implementor_groups.IsUsersGroupsInCRImplementerGroups(impGroups, cr.ImplementorGroups) {
		p.log.Warn("User is not in change implementer group", zap.String("user", r.GetUserId()))
		return errors.New("user is not in change implementer group")
	}

	if !validators.IsServerIpInObjects(r.GetServerInfo().IP, cr.CyberArkObjects) {
		p.log.Warn("Server IP is not in change request", zap.String("ip", r.GetServerInfo().IP))
		return errors.New("server IP is not in change request cyberark objects")
	}
	osAdmUsers := p.osAdmUsersSvc.GetNonCrOsUsers(r.GetUserId())

	if !validators.IsOsUserInObjects(r.GetServerInfo().OSUser, cr.CyberArkObjects) &&
		!validators.IsOsUserInOsAdmUsers(r.GetServerInfo().OSUser, osAdmUsers) {
		p.log.Warn("OS User is not in change request", zap.String("osUser", r.GetServerInfo().OSUser))
		return errors.New("OS User is not in change request cyberark objects")
	}

	if !validators.IsApproved(cr.State) {
		p.log.Warn("Change request is not approved", zap.String("changeID", r.GetChangeId()), zap.String("state", cr.State))
		return errors.New("change request is not approved")
	}

	if p.next != nil {
		return p.next.Check(r)
	}
	return nil
}
