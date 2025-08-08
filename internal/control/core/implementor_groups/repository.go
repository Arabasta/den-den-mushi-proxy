package implementor_groups

import dto "den-den-mushi-Go/pkg/dto/implementor_groups"

type Repository interface {
	FindAllByUserId(userId string) ([]*dto.Record, error)
}
