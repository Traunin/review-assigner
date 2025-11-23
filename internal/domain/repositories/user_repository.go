package repositories

import (
	"context"

	"github.com/Traunin/review-assigner/internal/domain/entities"
)

type UserRepository interface {
	Repository[entities.User, entities.UserID]
	GetActiveUsers(ctx context.Context) ([]*entities.User, error)
	GetByTeamID(ctx context.Context, id entities.TeamID) ([]*entities.User, error)
}
