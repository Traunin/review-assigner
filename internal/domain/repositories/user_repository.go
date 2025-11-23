package repositories

import (
	"context"

	"github.com/Traunin/review-assigner/internal/domain/entities"
)

type UserRepository interface {
	Repository[entities.User, entities.UserID]
	GetActiveUsers(ctx context.Context) ([]*entities.User, error)
}
