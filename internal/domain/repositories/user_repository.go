package repositories

import "github.com/Traunin/review-assigner/internal/domain/entities"

type UserRepository interface {
	Repository[entities.User, entities.UserID]
	GetActiveUsers() ([]*entities.User, error)
}
