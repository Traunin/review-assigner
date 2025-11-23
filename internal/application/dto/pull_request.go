package dto

import (
	"time"

	"github.com/Traunin/review-assigner/internal/domain/entities"
)

type CreatePRCmd struct {
	PullRequestID   entities.PullRequestID
	PullRequestName string
	AuthorID        entities.UserID
}

type PullRequestDTO struct {
	PullRequestID   entities.PullRequestID
	PullRequestName string
	AuthorID        entities.UserID
	Status          string
	CreatedAt       time.Time
	MergedAt        *time.Time
	Reviewers       []ReviewerDTO
}
