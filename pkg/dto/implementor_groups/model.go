package implementor_groups

type Model struct {
	MemberName       string `gorm:"column:MemberName"`
	GroupName        string `gorm:"column:Group"`
	MembershipStatus string `gorm:"column:GroupMembershipStatus"`
}

func (Model) TableName() string {
	return "ichamp_implementer_groups"
}
