package repositories

import (
	"context"

	"github.com/Traunin/review-assigner/internal/domain/entities"
)

type PullRequestRepository interface {
	Repository[entities.PullRequest, entities.PullRequestID]
	FindPullRequestByUserID(
		ctx context.Context,
		id entities.UserID,
	) ([]*entities.PullRequest, error)
	FindOpenPullRequests(ctx context.Context) ([]*entities.PullRequest, error)
}
