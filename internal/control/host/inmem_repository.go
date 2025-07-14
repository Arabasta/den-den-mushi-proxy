package host

type InMemRepository struct {
}

func NewInMemRepository() *InMemRepository {
	return &InMemRepository{}
}

func (r *InMemRepository) FindByIp(ip string) (*Entity, error) {
	return nil, nil
}
