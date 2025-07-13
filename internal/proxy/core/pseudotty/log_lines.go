package pseudotty

import (
	"den-den-mushi-Go/pkg/types"
	"time"
)

func getLogHeader(s *Session) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	claims := s.primary.Claims

	header :=
		"# Session Start Time: " + time.Now().UTC().Format(time.RFC3339) + "\n" +
			"# Created By: " + claims.Subject + "\n\n"

	// todo: add ou group

	header += "# Connection Details:\n" +
		"#\t- Server IP: " + claims.Connection.Server.IP + "\n" +
		"#\t- OS User: " + claims.Connection.Server.OSUser + "\n" +
		"#\t- Port: " + claims.Connection.Server.Port + "\n\n"

	header += "# Purpose: " + string(claims.Connection.Purpose) + "\n\n"

	if claims.Connection.Purpose == types.Change {
		header +=
			"# Change Request Details:\n" +
				"#\t- Change Request ID: " + claims.Connection.ChangeRequest.Id + "\n" +
				"#\t- Implementor Group: " + claims.Connection.ChangeRequest.ImplementorGroup + "\n" +
				"#\t- End Time: " + claims.Connection.ChangeRequest.EndTime
	} else if claims.Connection.Purpose == types.Healthcheck {
		header +=
			"# Health Check Details:\n" +
				"#\t- Filter: " + string(claims.Connection.FilterType)
		// todo: add more stuff
	} else {
		header += "# No additional details for this session purpose\n"
	}

	header += "\n\n\n"

	return header
}

func getLogFooter(s *Session) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	footer := "\n# Session End Time: " + s.endTime
	// todo: add list of all users
	return footer
}
