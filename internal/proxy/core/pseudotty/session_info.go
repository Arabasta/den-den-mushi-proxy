package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/core/client"
	"den-den-mushi-Go/pkg/dto"
	"den-den-mushi-Go/pkg/types"
)

// todo: move to redis and split admin service and proxy service

type ConnectionInfo struct {
	UserSessionID string          `json:"user_session_id"`
	UserID        string          `json:"user_id"`
	StartRole     types.StartRole `json:"start_role,required"`
	JoinTime      string          `json:"join_time,omitempty"`
	LeaveTime     string          `json:"leave_time,omitempty"`
}

func extractConnectionInfo(c *client.Connection) ConnectionInfo {
	return ConnectionInfo{
		UserSessionID: c.Claims.Connection.UserSession.Id,
		UserID:        c.Claims.Subject,
		StartRole:     c.Claims.Connection.UserSession.StartRole,
		JoinTime:      c.JoinTime,
		LeaveTime:     c.LeaveTime,
	}
}

type SessionInfo struct {
	SessionID              string           `json:"session_id"`
	ProxyDetails           ProxyDetails     `json:"proxy_details"`
	StartConnectionDetails dto.Connection   `json:"start_connection_details"`
	CreatedBy              string           `json:"created_by"`
	StartTime              string           `json:"start_time"`
	EndTime                string           `json:"end_time,omitempty"`
	State                  string           `json:"state,omitempty"`         // todo: use enum
	LastActivity           string           `json:"last_activity,omitempty"` // ISO 8601 format
	ActiveConnections      []ConnectionInfo `json:"active_connections"`
	LivetimeConnections    []ConnectionInfo `json:"livetime_connections"`
}

type ProxyDetails struct {
	Hostname    string `json:"hostname"`
	IP          string `json:"ip"`
	Type        string `json:"type"`
	Region      string `json:"region"`
	Environment string `json:"environment"`
}

func (s *Session) GetDetails() SessionInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var activeParticipants []ConnectionInfo
	if s.activePrimary != nil && s.activePrimary.Claims != nil {
		activeParticipants = append(activeParticipants, extractConnectionInfo(s.activePrimary))
	}

	for o := range s.activeObservers {
		if o.Claims != nil {
			activeParticipants = append(activeParticipants, extractConnectionInfo(o))
		}
	}

	var livetimeConnections []ConnectionInfo
	for c, _ := range s.livetimeConnections {
		if c.Claims != nil {
			livetimeConnections = append(livetimeConnections, extractConnectionInfo(c))
		}
	}

	// todo:
	proxyDetails := ProxyDetails{}

	return SessionInfo{
		SessionID:              s.id,
		ProxyDetails:           proxyDetails,
		StartConnectionDetails: s.startClaims.Connection,
		CreatedBy:              s.startClaims.Subject,
		StartTime:              s.startTime,
		EndTime:                s.endTime,
		State: func() string {
			if s.Closed {
				return "closed"
			}
			return "active"
		}(),
		LastActivity:        "", // TODO: add last activity tracking
		ActiveConnections:   activeParticipants,
		LivetimeConnections: livetimeConnections,
	}
}
