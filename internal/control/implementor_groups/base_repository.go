package implementor_groups

type Repository interface {
	FindAllByUserId(userId string) ([]*Entity, error)
	FindByGroupName(g string) (*Entity, error)
}
