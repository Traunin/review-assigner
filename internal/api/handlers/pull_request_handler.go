package handlers

import (
	"errors"
	"net/http"

	"github.com/Traunin/review-assigner/internal/application/dto"
	"github.com/Traunin/review-assigner/internal/application/services"
	"github.com/Traunin/review-assigner/internal/domain/entities"
	domainservices "github.com/Traunin/review-assigner/internal/domain/services"
	"github.com/labstack/echo/v4"
)

// PostPullRequestCreate handles POST /pullRequest/create
func (s *Server) PostPullRequestCreate(ctx echo.Context) error {
	var req struct {
		PullRequestID   string `json:"pull_request_id"`
		PullRequestName string `json:"pull_request_name"`
		AuthorID        string `json:"author_id"`
	}

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]any{
			"error": map[string]string{
				"code":    "INVALID_REQUEST",
				"message": "invalid request body",
			},
		})
	}

	cmd := dto.CreatePRCmd{
		PullRequestID:   entities.PullRequestID(req.PullRequestID),
		PullRequestName: req.PullRequestName,
		AuthorID:        entities.UserID(req.AuthorID),
	}

	pr, err := s.prService.Create(ctx.Request().Context(), cmd)
	if err != nil {
		// Check for domain-specific errors
		if errors.Is(err, domainservices.ErrAuthorNotFound) || 
		   errors.Is(err, domainservices.ErrTeamNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]any{
				"error": map[string]string{
					"code":    "NOT_FOUND",
					"message": err.Error(),
				},
			})
		}
		if errors.Is(err, domainservices.ErrPRAlreadyExists) {
			return ctx.JSON(http.StatusConflict, map[string]any{
				"error": map[string]string{
					"code":    "PR_EXISTS",
					"message": "PR id already exists",
				},
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]any{
			"error": map[string]string{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
	}

	return ctx.JSON(http.StatusCreated, map[string]any{
		"pr": formatPullRequest(pr),
	})
}

func (s *Server) PostPullRequestMerge(ctx echo.Context) error {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
	}

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]any{
			"error": map[string]string{
				"code":    "INVALID_REQUEST",
				"message": "invalid request body",
			},
		})
	}

	pr, err := s.prService.Merge(
		ctx.Request().Context(),
		entities.PullRequestID(req.PullRequestID),
	)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]any{
				"error": map[string]string{
					"code":    "NOT_FOUND",
					"message": "PR not found",
				},
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]any{
			"error": map[string]string{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"pr": formatPullRequest(pr),
	})
}

func (s *Server) PostPullRequestReassign(ctx echo.Context) error {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
		OldUserID     string `json:"old_user_id"`
	}

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]any{
			"error": map[string]string{
				"code":    "INVALID_REQUEST",
				"message": "invalid request body",
			},
		})
	}

	cmd := dto.ReassignReviewerCmd{
		OldUserID:     entities.UserID(req.OldUserID),
		PullRequestID: entities.PullRequestID(req.PullRequestID),
	}

	result, err := s.prService.ReassignReviewer(ctx.Request().Context(), cmd)
	if err != nil {
		if errors.Is(err, domainservices.ErrPRNotFound) || 
		   errors.Is(err, domainservices.ErrUserNotReviewer) {
			return ctx.JSON(http.StatusNotFound, map[string]any{
				"error": map[string]string{
					"code":    "NOT_FOUND",
					"message": err.Error(),
				},
			})
		}
		if errors.Is(err, domainservices.ErrPRAlreadyMerged) {
			return ctx.JSON(http.StatusConflict, map[string]any{
				"error": map[string]string{
					"code":    "PR_MERGED",
					"message": "cannot reassign on merged PR",
				},
			})
		}
		if errors.Is(err, domainservices.ErrUserNotReviewer) {
			return ctx.JSON(http.StatusConflict, map[string]any{
				"error": map[string]string{
					"code":    "NOT_ASSIGNED",
					"message": "reviewer is not assigned to this PR",
				},
			})
		}
		if errors.Is(err, domainservices.ErrNoCandidate) {
			return ctx.JSON(http.StatusConflict, map[string]any{
				"error": map[string]string{
					"code":    "NO_CANDIDATE",
					"message": "no active replacement candidate in team",
				},
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]any{
			"error": map[string]string{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
	}

	return ctx.JSON(http.StatusOK, map[string]any{
		"pr":          formatPullRequest(*result.PullRequestID),
		"replaced_by": string(result.Assigned),
	})
}

func formatPullRequest(pr dto.PullRequestDTO) map[string]any {
	assignedReviewers := make([]string, len(pr.Reviewers))
	for i, r := range pr.Reviewers {
		assignedReviewers[i] = string(r.UserID)
	}

	response := map[string]any{
		"pull_request_id":    string(pr.PullRequestID),
		"pull_request_name":  pr.PullRequestName,
		"author_id":          string(pr.AuthorID),
		"status":             pr.Status,
		"assigned_reviewers": assignedReviewers,
		"createdAt":          pr.CreatedAt,
	}

	if pr.MergedAt != nil {
		response["mergedAt"] = pr.MergedAt
	}

	return response
}
