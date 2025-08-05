package jti

import "time"

type Model struct {
	Id        string    `gorm:"primaryKey" json:"id"`
	ExpiresAt time.Time `gorm:"index" json:"expires_at"`
	Subject   string    `gorm:"index" json:"subject"`
	FromProxy string    `gorm:"index" json:"from_proxy"`
}

func (Model) TableName() string {
	return "ddm_jti"
}

type Record struct {
	Id        string
	ExpiresAt time.Time
	Subject   string
	FromProxy string
}

func ToModel(r *Record) *Model {
	return &Model{
		Id:        r.Id,
		ExpiresAt: r.ExpiresAt,
		Subject:   r.Subject,
		FromProxy: r.FromProxy,
	}
}

func ToRecord(id string, expiresAt time.Time, subject string, fromProxy string) *Record {
	return &Record{
		Id:        id,
		ExpiresAt: expiresAt,
		Subject:   subject,
		FromProxy: fromProxy,
	}
}
