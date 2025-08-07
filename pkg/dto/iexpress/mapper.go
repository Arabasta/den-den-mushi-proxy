package iexpress

import (
	"den-den-mushi-Go/pkg/util/convert"
	"errors"
	"fmt"
)

const timeFmt = "2006-01-02 15:04:05"

func FromModel(m *Model) (*Record, error) {
	if m == nil {
		return nil, nil
	}

	start, ok := convert.ParseTime(m.StartTime, timeFmt)
	end, ok := convert.ParseTime(m.EndTime, timeFmt)

	if !ok {
		return nil, errors.New("failed to parse time")
	}

	return &Record{
		RequestId:       m.RequestId,
		OriginCountry:   m.OriginCountry,
		Lob:             m.Lob,
		Requestor:       m.Requestor,
		AppImpacted:     convert.CsvToSlice(m.AppImpacted),
		Action:          m.Action,
		RelatedTicket:   m.RelatedTicket,
		StartTime:       &start,
		EndTime:         &end,
		State:           m.State,
		CyberArkObjects: convert.CsvToSlice(m.CyberArkObjects),
		ApproverGroup1:  m.ApproverGroup1,
		ApproverGroup2:  m.ApproverGroup2,
		MDApproverGroup: m.MDApproverGroup,
	}, nil
}

func FromModels(models []Model) ([]*Record, error) {
	if models == nil {
		return nil, nil
	}

	records := make([]*Record, 0, len(models))
	for _, m := range models {
		record, err := FromModel(&m)
		if err != nil {
			fmt.Println("Failed to convert model to record:", m.RequestId, err)
			continue
		}
		records = append(records, record)
	}

	return records, nil
}
