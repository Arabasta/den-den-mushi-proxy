package implementor_groups

func FromModel(m *Model) *Record {
	if m == nil {
		return nil
	}
	return &Record{
		MemberName:       m.MemberName,
		GroupName:        m.GroupName,
		MembershipStatus: m.MembershipStatus,
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
