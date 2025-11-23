package postgres

import (
	"context"
	"time"

	"github.com/Traunin/review-assigner/internal/domain/entities"
	"github.com/Traunin/review-assigner/internal/infrastructure/db/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type PullRequestRepository struct {
	db *DB
}

func NewPullRequestRepository(db *DB) *PullRequestRepository {
	return &PullRequestRepository{
		db: db,
	}
}

func prStatusToDomain(dbStatus string) entities.PRStatus {
	switch dbStatus {
	case "OPEN":
		return entities.StatusOpen
	case "MERGED":
		return entities.StatusMerged
	default:
		return entities.StatusOpen
	}
}

func prStatusToDB(status entities.PRStatus) string {
	switch status {
	case entities.StatusOpen:
		return "OPEN"
	case entities.StatusMerged:
		return "MERGED"
	default:
		return "OPEN"
	}
}

func timeToPgTimestamptz(t time.Time) pgtype.Timestamptz {
    if t.IsZero() {
        return pgtype.Timestamptz{Valid: false}
    }
    return pgtype.Timestamptz{
        Time:  t,
        Valid: true,
    }
}

func timePtrToPgTimestamptz(t *time.Time) pgtype.Timestamptz {
    if t == nil || t.IsZero() {
        return pgtype.Timestamptz{Valid: false}
    }
    return pgtype.Timestamptz{
        Time:  *t,
        Valid: true,
    }
}

func pgTimestamptzToTime(ts pgtype.Timestamptz) time.Time {
	if !ts.Valid {
		return time.Time{}
	}
	return ts.Time
}

func (r *PullRequestRepository) buildPullRequestWithReviewers(
	ctx context.Context,
	prRow sqlc.PullRequest,
) (*entities.PullRequest, error) {
	reviewerRows, err := r.db.Queries.GetReviewersByPR(ctx, prRow.PullRequestID)
	if err != nil {
		return nil, err
	}

	reviewers := make([]entities.Reviewer, len(reviewerRows))
	for i, reviewer := range reviewerRows {
		reviewers[i] = entities.Reviewer{
			UserID:     entities.UserID(reviewer.UserID),
			AssignedAt: pgTimestamptzToTime(reviewer.AssignedAt),
		}
	}

	var mergedAt *time.Time
	if prRow.MergedAt.Valid {
		t := prRow.MergedAt.Time
		mergedAt = &t
	}

	pr, err := entities.NewPullRequest(
		entities.PullRequestID(prRow.PullRequestID),
		prRow.PullRequestName,
		entities.UserID(prRow.AuthorID),
		prStatusToDomain(prRow.Status),
		reviewers,
		pgTimestamptzToTime(prRow.CreatedAt),
		mergedAt,
	)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (r *PullRequestRepository) Create(
	ctx context.Context,
	pr *entities.PullRequest,
) error {
	return r.db.execTx(ctx, func(q *sqlc.Queries) error {
		_, err := q.CreatePullRequest(ctx, sqlc.CreatePullRequestParams{
			PullRequestID:   pr.ID().String(),
			PullRequestName: pr.Name(),
			AuthorID:        pr.AuthorID().String(),
			Status:          prStatusToDB(pr.Status()),
			CreatedAt:       timeToPgTimestamptz(pr.CreatedAt()),
		})
		if err != nil {
			return err
		}

		for _, reviewer := range pr.Reviewers() {
			if err := q.AddReviewer(ctx, sqlc.AddReviewerParams{
				PullRequestID: pr.ID().String(),
				UserID:        reviewer.UserID.String(),
				AssignedAt:    timeToPgTimestamptz(reviewer.AssignedAt),
			}); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *PullRequestRepository) FindByID(
	ctx context.Context,
	id entities.PullRequestID,
) (*entities.PullRequest, error) {
	prRow, err := r.db.Queries.GetPullRequestByID(ctx, id.String())
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return r.buildPullRequestWithReviewers(ctx, prRow)
}

func (r *PullRequestRepository) FindAll(
	ctx context.Context,
) ([]*entities.PullRequest, error) {
	prRows, err := r.db.Queries.GetPullRequests(ctx)
	if err != nil {
		return nil, err
	}

	prs := make([]*entities.PullRequest, len(prRows))
	for i, prRow := range prRows {
		pr, err := r.buildPullRequestWithReviewers(ctx, prRow)
		if err != nil {
			return nil, err
		}
		prs[i] = pr
	}

	return prs, nil
}

func (r *PullRequestRepository) Update(
	ctx context.Context,
	pr *entities.PullRequest,
) error {
	return r.db.execTx(ctx, func(q *sqlc.Queries) error {
		_, err := q.UpdatePRStatus(ctx, sqlc.UpdatePRStatusParams{
			PullRequestID: pr.ID().String(),
			Status:        prStatusToDB(pr.Status()),
			MergedAt:      timePtrToPgTimestamptz(pr.MergedAtPtr()),
		})
		if err != nil {
			return err
		}

		currentReviewers, err := q.GetReviewersByPR(ctx, pr.ID().String())
		if err != nil {
			return err
		}

		currentReviewerMap := make(map[string]bool)
		for _, reviewer := range currentReviewers {
			currentReviewerMap[reviewer.UserID] = true
		}

		newReviewerMap := make(map[string]entities.Reviewer)
		for _, reviewer := range pr.Reviewers() {
			newReviewerMap[reviewer.UserID.String()] = reviewer
		}

		for _, reviewer := range currentReviewers {
			if _, exists := newReviewerMap[reviewer.UserID]; !exists {
				if err := q.RemoveReviewer(ctx, sqlc.RemoveReviewerParams{
					PullRequestID: pr.ID().String(),
					UserID:        reviewer.UserID,
				}); err != nil {
					return err
				}
			}
		}

		for _, reviewer := range pr.Reviewers() {
			if !currentReviewerMap[reviewer.UserID.String()] {
				if err := q.AddReviewer(ctx, sqlc.AddReviewerParams{
					PullRequestID: pr.ID().String(),
					UserID:        reviewer.UserID.String(),
					AssignedAt:    timeToPgTimestamptz(reviewer.AssignedAt),
				}); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *PullRequestRepository) DeleteByID(
	ctx context.Context,
	id entities.PullRequestID,
) error {
	return r.db.Queries.DeletePullRequest(ctx, id.String())
}

func (r *PullRequestRepository) FindPullRequestByUserID(
	ctx context.Context,
	id entities.UserID,
) ([]*entities.PullRequest, error) {
	prRows, err := r.db.Queries.GetPRsByReviewer(ctx, id.String())
	if err != nil {
		return nil, err
	}

	prs := make([]*entities.PullRequest, len(prRows))
	for i, prRow := range prRows {
		fullPR, err := r.db.Queries.GetPullRequestByID(ctx, prRow.PullRequestID)
		if err != nil {
			return nil, err
		}

		pr, err := r.buildPullRequestWithReviewers(ctx, fullPR)
		if err != nil {
			return nil, err
		}
		prs[i] = pr
	}

	return prs, nil
}

func (r *PullRequestRepository) FindOpenPullRequests(
	ctx context.Context,
) ([]*entities.PullRequest, error) {
	prRows, err := r.db.Queries.GetOpenPRs(ctx)
	if err != nil {
		return nil, err
	}

	prs := make([]*entities.PullRequest, len(prRows))
	for i, prRow := range prRows {
		pr, err := r.buildPullRequestWithReviewers(ctx, prRow)
		if err != nil {
			return nil, err
		}
		prs[i] = pr
	}

	return prs, nil
}
