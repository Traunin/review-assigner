package dto

import (
	"time"

	"github.com/Traunin/review-assigner/internal/domain/entities"
)

type ReviewerDTO struct {
	UserID     entities.UserID
	AssignedAt time.Time
}

type ReassignReviewerCmd struct {
	OldUserID     entities.UserID
	PullRequestID entities.PullRequestID
}

type ReassignedDTO struct {
	PullRequestID *PullRequestDTO
	Assigned      entities.UserID
}
