package implementor_groups

type Model struct {
	MemberName       string `gorm:"column:MemberName"`
	MemberEmail      string `gorm:"column:MemberEmail"`
	GroupName        string `gorm:"column:Group"`
	MembershipStatus string `gorm:"column:GroupMembershipStatus"`
}

func (Model) TableName() string {
	return "ICHAMP_IMPLEMENTER_GROUPS"
}
