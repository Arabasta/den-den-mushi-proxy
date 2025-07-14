package policy

import (
	"den-den-mushi-Go/internal/control/change_request"
	"den-den-mushi-Go/internal/control/dto"
	"den-den-mushi-Go/internal/control/host"
	"den-den-mushi-Go/internal/control/implementor_groups"
	"den-den-mushi-Go/pkg/types"
	"errors"
	"go.uber.org/zap"
)

type ChangePolicy[T dto.RequestCtx] struct {
	next Policy[T]

	crService       *change_request.Service
	impGroupService *implementor_groups.Service
	hostService     *host.Service
	log             *zap.Logger
}

func NewChangePolicy[T dto.RequestCtx](
	crService *change_request.Service,
	impGroupService *implementor_groups.Service,
	hostService *host.Service, log *zap.Logger) *ChangePolicy[T] {
	return &ChangePolicy[T]{
		crService:       crService,
		impGroupService: impGroupService,
		hostService:     hostService,
		log:             log,
	}
}

func (p *ChangePolicy[T]) SetNext(n Policy[T]) {
	p.next = n
}

func (p *ChangePolicy[T]) Check(req T) error {
	// 1. skip non-change requests
	if req.GetPurpose() != types.Change {
		if p.next != nil {
			return p.next.Check(req)
		}
		return nil
	}

	// 2. check for Change ID
	if req.GetChangeId() == "" {
		return errors.New("changeID is empty")
	}

	// 3. get Change Request
	cr, err := p.crService.FindById(req.GetChangeId())
	if err != nil {
		p.log.Error("Failed to find change request", zap.String("changeID", req.GetChangeId()), zap.Error(err))
		return err
	}
	if cr == nil {
		p.log.Error("Change request not found", zap.String("changeID", req.GetChangeId()))
		return errors.New("change request not found")
	}

	// 4. check if CR valid
	if !isValidWindow(cr.ChangeStartTime, cr.ChangeEndTime) {
		p.log.Error("Change request is not valid", zap.String("changeID", req.GetChangeId()),
			zap.String("start", cr.ChangeStartTime), zap.String("end", cr.ChangeEndTime))
		return errors.New("change request is not valid")
	}

	if !p.isUserInChangeImplementerGroup(req.GetUserId(), cr.ImplementorGroups) {
		p.log.Error("User is not in change implementer group", zap.String("user", req.GetUserId()))
		return errors.New("user is not in change implementer group")
	}

	if !p.isServerIpInChangeRequest(req.GetServerInfo().IP, cr.CyberArkObjects) {
		p.log.Error("Server IP is not in change request", zap.String("ip", req.GetServerInfo().IP))
		return errors.New("server IP is not in change request")
	}

	if !p.isOsUserInChangeRequest(req.GetServerInfo().OSUser, cr.CyberArkObjects) {
		p.log.Error("OS User is not in change request", zap.String("osUser", req.GetServerInfo().OSUser))
		return errors.New("OS User is not in change request")
	}

	if !p.isCRApproved(cr.State) {
		p.log.Error("Change request is not approved", zap.String("changeID", req.GetChangeId()), zap.String("state", cr.State))
		return errors.New("change request is not approved")
	}

	return nil
}
