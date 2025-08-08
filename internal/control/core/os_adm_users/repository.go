package os_adm_users

import dto "den-den-mushi-Go/pkg/dto/os_adm_users"

type Repository interface {
	FindAllByUserId(userId string) ([]*dto.Record, error)
}
