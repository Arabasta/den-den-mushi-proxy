package os_adm_users

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID                   uint           `gorm:"column:ID;primaryKey;autoIncrement"`
	UserID               string         `gorm:"column:UserID;not null"`
	OsUser               string         `gorm:"column:OsUser;not null"`
	Platform             string         `gorm:"column:Platform"`             // e.g., "Linux", "Windows", "AIX"
	ServerClassification string         `gorm:"column:ServerClassification"` // e.g., "OS", "DB", "Storage"
	CreatedAt            time.Time      `gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt            time.Time      `gorm:"column:UpdatedAt;autoUpdateTime"`
	DeletedAt            gorm.DeletedAt `gorm:"column:DeletedAt"`

	// todo: use ddm_users table, not doing now cause rush
	// User      users.Model    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

func (Model) TableName() string {
	return "ddm_os_adm_users"
}
