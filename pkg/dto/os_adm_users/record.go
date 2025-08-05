package os_adm_users

import "time"

type Record struct {
	UserId               string
	OsUser               string
	Platform             string
	ServerClassification string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
