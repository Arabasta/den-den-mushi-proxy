package roles

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        uint           `gorm:"column:ID;primaryKey;autoIncrement"`
	Name      string         `gorm:"column:Name;not null;uniqueIndex"`
	CreatedAt time.Time      `gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:UpdatedAt;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:DeletedAt;index"`
}

func (Model) TableName() string {
	return "ddm_roles"
}
