package cyberark

func FromModel(m *Model) *Record {
	if m == nil {
		return nil
	}
	return &Record{
		Object:   m.Object,
		Hostname: m.Hostname,
		Ip:       m.Ip,
	}
}
