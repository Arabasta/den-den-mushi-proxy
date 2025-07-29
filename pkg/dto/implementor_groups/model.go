package implementor_groups

type Model struct {
	MemberName       string `gorm:"column:MemberName"`
	MemberEmail      string `gorm:"column:MemberEmail"`
	GroupName        string `gorm:"column:Group"`
	MembershipStatus string `gorm:"column:GroupMembershipStatus"`
}

var tableName = "ichamp_implementer_groups"

func (Model) TableName() string {
	return tableName
}

func SetTableName(name string) {
	tableName = name
}
