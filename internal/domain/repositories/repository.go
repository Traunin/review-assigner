package repositories

import "context"

type Repository[T any, ID any] interface {
	Create(ctx context.Context, entity *T) error
	DeleteByID(ctx context.Context, id ID) error
	FindByID(ctx context.Context, id ID) (*T, error)
	FindAll(ctx context.Context) ([]*T, error)
	Update(ctx context.Context, entity *T) error
}
