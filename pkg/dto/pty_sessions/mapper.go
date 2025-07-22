package pty_sessions

import (
	"den-den-mushi-Go/pkg/dto/connections"
)

func FromModel(m *Model) *Record {
	if m == nil {
		return nil
	}
	return &Record{
		ID:                       m.ID,
		CreatedBy:                m.CreatedBy,
		StartTime:                m.StartTime,
		EndTime:                  m.EndTime,
		State:                    m.State,
		LastActivity:             m.LastActivity,
		ProxyHostName:            m.ProxyHostName,
		ProxyDetails:             m.ProxyDetails,
		StartConnServerIP:        m.StartConnServerIP,
		StartConnServerPort:      m.StartConnServerPort,
		StartConnServerOSUser:    m.StartConnServerOSUser,
		StartConnType:            m.StartConnType,
		StartConnPurpose:         m.StartConnPurpose,
		StartConnUserSessionID:   m.StartConnUserSessionID,
		StartConnChangeRequestID: m.StartConnChangeRequestID,
		StartConnFilterType:      m.StartConnFilterType,
		Connections:              connections.FromModels(m.Connections),
	}
}

func FromModels(models []Model) []*Record {
	if models == nil {
		return nil
	}
	records := make([]*Record, len(models))
	for i, m := range models {
		records[i] = FromModel(&m)
	}
	return records
}

func ToModel(r *Record) *Model {
	if r == nil {
		return nil
	}
	return &Model{
		ID:                       r.ID,
		CreatedBy:                r.CreatedBy,
		StartTime:                r.StartTime,
		EndTime:                  r.EndTime,
		State:                    r.State,
		LastActivity:             r.LastActivity,
		ProxyHostName:            r.ProxyHostName,
		ProxyDetails:             r.ProxyDetails,
		StartConnServerIP:        r.StartConnServerIP,
		StartConnServerPort:      r.StartConnServerPort,
		StartConnServerOSUser:    r.StartConnServerOSUser,
		StartConnType:            r.StartConnType,
		StartConnPurpose:         r.StartConnPurpose,
		StartConnUserSessionID:   r.StartConnUserSessionID,
		StartConnChangeRequestID: r.StartConnChangeRequestID,
		StartConnFilterType:      r.StartConnFilterType,
		Connections:              connections.ToModels(r.Connections),
	}
}
