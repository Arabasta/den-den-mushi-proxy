package policy

import (
	"den-den-mushi-Go/internal/control/implementor_groups"
	"den-den-mushi-Go/internal/control/pty_token/request"
	"go.uber.org/zap"
)

type ImplementorGroupPolicy[T request.Ctx] struct {
	next Policy[T]

	impGroupService *implementor_groups.Service
	log             *zap.Logger
}

func NewImplementorGroupPolicy[T request.Ctx](impGSvc *implementor_groups.Service, log *zap.Logger) *ChangePolicy[T] {
	return &ChangePolicy[T]{
		impGroupService: impGSvc,
		log:             log,
	}
}

func (p *ImplementorGroupPolicy[T]) SetNext(n Policy[T]) {
	p.next = n
}

func (p *ImplementorGroupPolicy[T]) Check(r T) error {
	// 1. get user's implementor groups  todo: get from Request context
	//impGroups, err := p.impGroupService.FindAllByUserId(r.GetUserId())
	//if err != nil {
	//	p.log.Warn("Failed to find implementor groups for user", zap.String("userId", r.GetUserId()), zap.Error(err))
	//	return err
	//}

	// todo load required imp groups in request context
	//if !implementor_groups.IsUsersGroupsInCRImplementerGroups(impGroups, cr.ImplementorGroups) {
	//	p.log.Warn("User is not in change implementer group", zap.String("user", r.GetUserId()))
	//	return errors.New("user is not in change implementer group")
	//}

	if p.next != nil {
		return p.next.Check(r)
	}
	return nil
}
