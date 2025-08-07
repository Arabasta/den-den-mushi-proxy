package iexpress

import "time"

type Record struct {
	RequestId       string
	OriginCountry   string
	Lob             string
	Requestor       string
	AppImpacted     []string
	Action          string
	RelatedTicket   string
	StartTime       *time.Time
	EndTime         *time.Time
	State           string
	CyberArkObjects []string
	ApproverGroup1  string
	ApproverGroup2  string
	MDApproverGroup string
}
