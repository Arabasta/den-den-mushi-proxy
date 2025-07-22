package implementor_groups

type Model struct {
	ID               uint   `gorm:"primaryKey;column:ID;type:bigint"` // todo: tmp for dev
	MemberName       string `gorm:"column:MemberName"`
	GroupName        string `gorm:"column:Group"`
	MembershipStatus string `gorm:"column:GroupMembershipStatus"`
}

func (Model) TableName() string {
	return "ichamp_implementor_groups"
}
