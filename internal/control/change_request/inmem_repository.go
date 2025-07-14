package change_request

type InMemRepository struct {
}

func NewInMemRepository() *InMemRepository {
	return &InMemRepository{}
}

func (r *InMemRepository) FindById(id string) (*Entity, error) {
	return nil, nil
}
