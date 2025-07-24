package change_request

type Model struct {
	CRNumber          string `gorm:"column:TicketNumber"`
	Country           string `gorm:"column:CountryImpacted"` // Comma-separated. e.g. "SG,HK"
	Lob               string `gorm:"column:LOB"`
	Summary           string `gorm:"column:ChangeSummary"`
	Description       string `gorm:"column:ChangeDescription"`
	ChangeStartTime   string `gorm:"column:ChangeSchedStartTime"` // 2025-06-24 23:00:00
	ChangeEndTime     string `gorm:"column:ChangeSchedEndTime"`   // 2025-06-25 01:00:00
	ImplementorGroups string `gorm:"column:ImplementorGroup"`     // Comma-separated. e.g. "Group1,Group2,Group3"
	State             string `gorm:"column:State"`                // Approved, Reopen, etc
	CyberArkObjects   string `gorm:"column:CyberarkObjects"`      // Comma-separated. e.g. "Object1,Object2,Object3"
}

func (Model) TableName() string {
	return "changerequests"
}
