package repositories

type Repository[T any, ID any] interface {
	Create(entity *T) error
	DeleteByID(id ID) error
	FindByID(id ID) (*T, error)
	FindAll() ([]*T, error)
	Update(entity *T) error
}
