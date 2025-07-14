package change_request

type Repository interface {
	FindById(id string) (*Entity, error)
}
