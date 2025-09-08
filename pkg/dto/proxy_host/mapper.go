package proxy_host

func FromModel(m *Model) *Record {
	if m == nil {
		return nil
	}
	return &Record{
		ProxyType:            m.ProxyType,
		IpAddress:            m.IpAddress,
		HostName:             m.HostName,
		OSType:               m.OSType,
		Status:               m.Status,
		Environment:          m.Environment,
		Country:              m.Country,
		LoadBalancerEndpoint: m.LoadBalancerEndpoint,
	}
}

func ToModel(r *Record) *Model {
	if r == nil {
		return nil
	}
	return &Model{
		ProxyType:            r.ProxyType,
		IpAddress:            r.IpAddress,
		HostName:             r.HostName,
		OSType:               r.OSType,
		Status:               r.Status,
		Environment:          r.Environment,
		Country:              r.Country,
		LoadBalancerEndpoint: r.LoadBalancerEndpoint,
	}
}

func FromModel2(m *Model2) *Record2 {
	if m == nil {
		return nil
	}
	return &Record2{
		HostName:        m.HostName,
		ProxyType:       m.ProxyType,
		Url:             m.Url,
		IpAddress:       m.IpAddress,
		Environment:     m.Environment,
		Country:         m.Country,
		LastHeartbeatAt: m.LastHeartbeatAt,
		DrainingStartAt: m.DrainingStartAt,
		DrainingEndAt:   m.DrainingEndAt,
		DeploymentColor: m.DeploymentColor,
		MaxSessions:     m.MaxSessions,
		ActiveSessions:  m.ActiveSessions,
	}
}

func FromModels2(ms []*Model2) []*Record2 {
	rs := make([]*Record2, 0, len(ms))
	for _, m := range ms {
		rs = append(rs, FromModel2(m))
	}
	return rs
}
