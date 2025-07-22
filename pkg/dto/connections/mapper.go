package connections

func FromModel(m *Model) *Record {
	return &Record{
		ID:           m.ID,
		UserID:       m.UserID,
		PtySessionID: m.PtySessionID,
		Status:       m.Status,
		StartRole:    m.StartRole,
		JoinTime:     m.JoinTime,
		LeaveTime:    m.LeaveTime,
	}
}

func FromModels(models []*Model) []*Record {
	if models == nil {
		return nil
	}
	records := make([]*Record, len(models))
	for i, m := range models {
		records[i] = FromModel(m)
	}
	return records
}

func ToModel(r *Record) *Model {
	if r == nil {
		return nil
	}
	return &Model{
		ID:           r.ID,
		UserID:       r.UserID,
		PtySessionID: r.PtySessionID,
		Status:       r.Status,
		StartRole:    r.StartRole,
		JoinTime:     r.JoinTime,
		LeaveTime:    r.LeaveTime,
	}
}

func ToModels(records []*Record) []*Model {
	if records == nil {
		return nil
	}
	models := make([]*Model, len(records))
	for i, r := range records {
		models[i] = ToModel(r)
	}
	return models
}
