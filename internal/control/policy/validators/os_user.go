package validators

import "den-den-mushi-Go/pkg/util/cyberark"

func (v *Validator) IsOsUserInObjects(osUser string, objects []string) bool {
	if osUser == "" {
		return false
	}

	return cyberark.IsOsUserInObjects(osUser, objects)
}

func (v *Validator) IsOsUserInOsAdmUsers(osUser string, osAdmUsers []string) bool {
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
