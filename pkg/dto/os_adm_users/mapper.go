package os_adm_users

func FromModel(m *Model) *Record {
	return &Record{
		UserId:               m.UserID,
		OsUser:               m.OsUser,
		Platform:             m.Platform,
		ServerClassification: m.ServerClassification,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}
}

func FromModels(models []Model) []*Record {
	if models == nil {
		return nil
	}
	records := make([]*Record, len(models))
	for i := range models {
		records[i] = FromModel(&models[i])
	}
	return records
}
