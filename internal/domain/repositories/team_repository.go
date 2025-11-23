package repositories

import (
	"context"

	"github.com/Traunin/review-assigner/internal/domain/entities"
)

type TeamRepository interface {
	Repository[entities.Team, entities.TeamID]
	FindByName(ctx context.Context, name string) (*entities.Team, error)
	FindTeamsByUserID(
		ctx context.Context,
		id entities.UserID,
	) (*entities.Team, error)
	FindActiveReviewersByTeamID(
		ctx context.Context,
		id entities.TeamID,
	) (*entities.User, error)
	TeamExists(ctx context.Context, name string) (bool, error)
}
