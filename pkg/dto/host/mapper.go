package host

import (
	"fmt"
)

func FromModel(m *Model) *Record {
	if m == nil {
		return nil
	}
	return &Record{
		InventoryID: fmt.Sprintf("%d", m.ID),
		IpAddress:   m.IpAddress,
		HostName:    m.HostName,
		OSType:      m.OSType,
		Appcode:     m.Appcode,
		Status:      m.Status,
		Environment: m.Environment,
		Country:     m.Country,
	}
}

func FromModels(models []Model) []*Record {
	if len(models) == 0 {
		return nil
	}
	records := make([]*Record, len(models))
	for i, m := range models {
		records[i] = FromModel(&m)
	}
	return records
}
