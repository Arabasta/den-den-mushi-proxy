package users

import (
	"den-den-mushi-Go/pkg/dto/roles"
	"gorm.io/gorm"
	"time"
)

type Model struct {
	Id        uint           `gorm:"column:ID;primaryKey;autoIncrement"`
	Username  string         `gorm:"column:Username;not null;uniqueIndex"`
	Roles     []roles.Model  `gorm:"many2many:ddm_user_roles;constraint:OnDelete:CASCADE;"`
	CreatedAt time.Time      `gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:UpdatedAt;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:DeletedAt;index;default:null"`
}

func (Model) TableName() string {
	return "ddm_users"
}
