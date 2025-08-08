package policy

import (
	"den-den-mushi-Go/internal/control/app/iexpress"
	"den-den-mushi-Go/internal/control/app/pty_token/request"
	implementor_groups2 "den-den-mushi-Go/internal/control/core/implementor_groups"
	"den-den-mushi-Go/internal/control/core/os_adm_users"
	"den-den-mushi-Go/internal/control/policy/validators"
	"den-den-mushi-Go/pkg/types"
	"errors"
	"go.uber.org/zap"
)

type IexpressPolicy[T request.Ctx] struct {
	next Policy[T]

	impGroupService *implementor_groups2.Service
	osAdmUsersSvc   *os_adm_users.Service
	v               *validators.Validator

	log *zap.Logger
}

func NewIExpressPolicy[T request.Ctx](impGroupSvc *implementor_groups2.Service, osAdmUsersSvc *os_adm_users.Service,
	v *validators.Validator, log *zap.Logger) *IexpressPolicy[T] {
	return &IexpressPolicy[T]{
		impGroupService: impGroupSvc,
		osAdmUsersSvc:   osAdmUsersSvc,
		v:               v,
		log:             log,
	}
}

func (p *IexpressPolicy[T]) SetNext(n Policy[T]) {
	p.next = n
}

func (p *IexpressPolicy[T]) Check(r T) error {
	p.log.Debug("Checking iexpress policy...")
	if r.GetPurpose() != types.IExpress {
		p.log.Debug("Skipping non-iexpress request", zap.String("purpose", string(r.GetPurpose())))
		if p.next != nil {
			return p.next.Check(r)
		}
		return nil
	}

	// get iexpress Request from ctx
	exp := r.GetIexpress()
	if exp == nil {
		// should never happen cause service should verify before calling this policy
		p.log.Error("IExpress request not found")
		return errors.New("IExpress request not found")
	}

	// 3b. get user's implementor groups  todo: get from Request context and move this to implementor group policy
	impGroups, err := p.impGroupService.FindAllByUserId(r.GetUserId())
	if err != nil {
		p.log.Warn("Failed to find implementor groups for user", zap.String("userId", r.GetUserId()), zap.Error(err))
		return err
	}
	enrichedImpG := make([]string, 0)
	for _, group := range impGroups {
		enrichedImpG = append(enrichedImpG, group.GroupName)
	}
	enrichedImpG = *iexpress.EnrichImpGroupsWithGOV_(enrichedImpG)

	// 4. check if ticket valid
	if !p.v.IsValidWindow(*exp.StartTime, *exp.EndTime) {
		p.log.Warn("IExpress window invalid", zap.String("IExpress", r.GetChangeId()), zap.String("startTime", exp.StartTime.String()), zap.String("endTime", exp.EndTime.String()))
		return errors.New("IExpress window is invalid")
	}

	groups := nonEmptyStrings(exp.ApproverGroup1, exp.ApproverGroup2, exp.MDApproverGroup)
	if !validators.IsUsersGroupsInCRImplementerGroups(impGroups, groups) {
		p.log.Warn("User is not in iexpress implementer group", zap.String("user", r.GetUserId()))
		return errors.New("user is not in iexpress implementer group")
	}

	if !p.v.IsServerIpInObjects(r.GetServerInfo().IP, exp.CyberArkObjects) {
		p.log.Warn("Server IP is not in iexpress request", zap.String("ip", r.GetServerInfo().IP))
		return errors.New("server IP is not in iexpress request cyberark objects")
	}
	osAdmUsers := p.osAdmUsersSvc.GetNonCrOsUsers(r.GetUserId())

	if !p.v.IsOsUserInObjects(r.GetServerInfo().OSUser, exp.CyberArkObjects) &&
		!p.v.IsOsUserInOsAdmUsers(r.GetServerInfo().OSUser, osAdmUsers) {
		p.log.Warn("OS User is not in iexpress request", zap.String("osUser", r.GetServerInfo().OSUser))
		return errors.New("OS User is not in iexpress request cyberark objects")
	}

	if !p.v.IsApproved(exp.State) {
		p.log.Warn("iexpress request is not approved", zap.String("changeID", r.GetChangeId()), zap.String("state", exp.State))
		return errors.New("iexpress request is not approved")
	}

	if p.next != nil {
		return p.next.Check(r)
	}
	return nil
}

func nonEmptyStrings(values ...string) []string {
	var out []string
	for _, v := range values {
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}
