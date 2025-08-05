package validators

import "den-den-mushi-Go/pkg/util/cyberark"

func IsOsUserInObjects(osUser string, objects []string) bool {
	if osUser == "" {
		return false
	}

	return cyberark.IsOsUserInObjects(osUser, objects)
}

func IsOsUserInOsAdmUsers(osUser string, osAdmUsers []string) bool {
	if osUser == "" {
		return false
	}

	for _, o := range osAdmUsers {
		if o == osUser {
			return true
		}
	}
	return false
}
