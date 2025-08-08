package validators

import dto "den-den-mushi-Go/pkg/dto/implementor_groups"

func IsUsersGroupsInCRImplementerGroups(userGroups []*dto.Record, crGroups []string) bool {
	if userGroups == nil || crGroups == nil {
		return false
	}

	userGroupSet := make(map[string]struct{}, len(userGroups))
	for _, g := range userGroups {
		userGroupSet[g.GroupName] = struct{}{}
	}

	for _, g := range crGroups {
		if _, exists := userGroupSet[g]; exists {
			return true
		}
	}

	return false
}
