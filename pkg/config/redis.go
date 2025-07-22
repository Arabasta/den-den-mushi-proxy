package config

type Redis struct {
	Addrs    []string
	Password string
	PoolSize int
}
