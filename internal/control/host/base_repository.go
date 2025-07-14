package host

type Repository interface {
	FindByIp(ip string) (*Entity, error)
}
