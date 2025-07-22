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
