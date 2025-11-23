package repositories

import "github.com/Traunin/review-assigner/internal/domain/entities"

type TeamRepository interface {
	Repository[entities.Team, entities.TeamID]
	FindByName(name string) (*entities.Team, error)
	FindTeamsByUserID(id entities.UserID) ([]*entities.Team, error)
	FindActiveReviewersByTeamID(id entities.TeamID) ([]*entities.User, error)
}
