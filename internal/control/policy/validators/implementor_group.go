package validators

import (
	dto "den-den-mushi-Go/pkg/dto/implementor_groups"
	"strings"
)

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

func (v *Validator) IsUserInAnyIexpressGroup(userGroups, iexpressGroups []string) bool {
	if len(userGroups) == 0 || len(iexpressGroups) == 0 {
		return false
	}
	norm := func(s string) string {
		return strings.ToLower(strings.TrimSpace(s))
	}

	set := make(map[string]struct{}, len(userGroups))
	for _, g := range userGroups {
		ng := norm(g)
		if ng != "" {
			set[ng] = struct{}{}
		}
	}
	for _, g := range iexpressGroups {
		ng := norm(g)
		if ng == "" {
			continue
		}
		if _, ok := set[ng]; ok {
			return true
		}
	}
	return false
}
