package session_logging

import (
	"den-den-mushi-Go/pkg/token"
	"den-den-mushi-Go/pkg/types"
	"fmt"
	"time"
)

func FormatLogLine(header, data string) string {
	return fmt.Sprintf("%s [%s] %s", time.Now().Format(time.DateTime), header, data)
}

func FormatInputOnlyLogLine(userId, data string) string {
	return fmt.Sprintf("%s %s: %s", time.Now().Format(time.DateTime), userId, data)
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
				"#\t- End Time: " + claims.Connection.ChangeRequest.EndTime.Format(time.RFC3339)
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

func FormatFooter(endTime time.Time) string {
	footer := "\n# Session End Time: " + endTime.Format(time.RFC3339)
	// todo: add list of all users
	return footer
}
