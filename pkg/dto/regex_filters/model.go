package regex_filters

import (
	"den-den-mushi-Go/pkg/types"
	"gorm.io/gorm"
	"time"
)

type Model struct {
	Id           uint           `gorm:"column:ID;primaryKey;autoIncrement"`
	FilterType   types.Filter   `gorm:"column:Type;size:50"`
	RegexPattern string         `gorm:"column:Pattern;type:text"`
	OuGroup      string         `gorm:"column:OUGroup;size:191"`
	IsEnabled    bool           `gorm:"column:IsEnabled;default:true"`
	CreatedAt    time.Time      `gorm:"column:CreatedAt;autoCreateTime"`
	CreatedBy    string         `gorm:"column:CreatedBy;size:191"`
	UpdatedAt    time.Time      `gorm:"column:UpdatedAt;autoUpdateTime"`
	UpdatedBy    string         `gorm:"column:UpdatedBy;size:191"`
	DeletedAt    gorm.DeletedAt `gorm:"column:DeletedAt;index;default:null"`
	DeletedBy    string         `gorm:"column:DeletedBy;size:191"`
}

func (Model) TableName() string {
	return "ddm_regex_filters"
}
