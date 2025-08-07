package iexpress

type Model struct {
	RequestId       string `gorm:"column:Ticket"`
	OriginCountry   string `gorm:"column:Country_of_Origin"` // eg SG
	Lob             string `gorm:"column:Unit"`              // eg IBGT
	Requestor       string `gorm:"column:Requestor"`
	AppImpacted     string `gorm:"column:Application_Impacted"` // comma separated. e.g. "App1,App2,App3"
	Action          string `gorm:"column:Action_for_Tickets"`
	RelatedTicket   string `gorm:"column:Related_Incident_Ticket"`
	StartTime       string `gorm:"column:Schedule_Start"`     // 2025-06-24 23:00:00
	EndTime         string `gorm:"column:Schedule_End"`       // 2025-06-25 01:00:00
	State           string `gorm:"column:State"`              // Approved, Reopen, etc
	CyberArkObjects string `gorm:"column:CyberArk_Objects_1"` // Comma-separated. e.g. "Object1,Object2,Object3"
	ApproverGroup1  string `gorm:"column:Approver_Group_1"`
	ApproverGroup2  string `gorm:"column:Approver_Group_2"`
	MDApproverGroup string `gorm:"column:MD_Approver_Group"`
}

func (Model) TableName() string {
	return "IEXPRESSREQUESTS"
}
