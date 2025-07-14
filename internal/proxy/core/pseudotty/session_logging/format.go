package session_logging

import (
	"den-den-mushi-Go/pkg/token"
	"den-den-mushi-Go/pkg/types"
	"fmt"
	"time"
)

func FormatLogLine(header, data string) string {
	return fmt.Sprintf("%s [%s] %s", time.Now().Format(time.TimeOnly), header, data)
}

func FormatHeader(claims *token.Claims) string {
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

func FormatFooter(endTime string) string {
	footer := "\n# Session End Time: " + endTime
	// todo: add list of all users
	return footer
}
