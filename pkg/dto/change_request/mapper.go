package change_request

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

	fmt.Println("Parsing time" + m.ChangeStartTime + " to " + timeFmt)
	fmt.Println("Parsing time" + m.ChangeEndTime + " to " + timeFmt)
	start, ok := convert.ParseTime(m.ChangeStartTime, timeFmt)
	end, ok := convert.ParseTime(m.ChangeEndTime, timeFmt)

	if !ok {
		return nil, errors.New("failed to parse change time")
	}

	return &Record{
		ChangeRequestId:   m.CRNumber,
		Country:           convert.CsvToSlice(m.Country),
		Lob:               m.Lob,
		Summary:           m.Summary,
		Description:       m.Description,
		ChangeStartTime:   &start,
		ChangeEndTime:     &end,
		ImplementorGroups: convert.CsvToSlice(m.ImplementorGroups),
		State:             m.State,
		CyberArkObjects:   convert.CsvToSlice(m.CyberArkObjects),
	}, nil
}

func FromModels(models []Model) ([]*Record, error) {
	if models == nil {
		return nil, nil
	}

	records := make([]*Record, 0, len(models))
	for _, m := range models {
		record, err := FromModel(&m)
		fmt.Println("Converted model to record:", m.CRNumber, "->", record)
		if err != nil {
			fmt.Println("Failed to convert model to record:", m.CRNumber, err)
			continue
		}
		records = append(records, record)
	}

	return records, nil
}
