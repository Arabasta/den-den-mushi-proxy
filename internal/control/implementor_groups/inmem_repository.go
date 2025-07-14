package implementor_groups

type InMemRepository struct {
}

func NewInMemRepository() *InMemRepository {
	return &InMemRepository{}
}

func (r *InMemRepository) FindAllByUserId(userId string) ([]*Entity, error) {
	return nil, nil
}

func (r *InMemRepository) FindByGroupName(g string) (*Entity, error) {
	return nil, nil
}
