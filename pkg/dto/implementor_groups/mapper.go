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
	if len(models) == 0 {
		return nil
	}
	records := make([]*Record, len(models))
	for i := range models {
		records[i] = FromModel(&models[i])
	}
	return records
}
