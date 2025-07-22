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
		Status:      m.Status,
		Environment: m.Environment,
		Country:     m.Country,
	}
}
