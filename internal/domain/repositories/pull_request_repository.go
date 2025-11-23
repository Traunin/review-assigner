package repositories

import "github.com/Traunin/review-assigner/internal/domain/entities"

type PullRequestRepository interface {
	Repository[entities.User, entities.UserID]
	FindPullRequestByUserID(id entities.UserID) ([]*entities.PullRequest, error)
	FindOpenPullRequests() ([]*entities.PullRequest, error)
}
