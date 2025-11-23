package handlers

import (
	"net/http"

	"github.com/Traunin/review-assigner/internal/domain/entities"
	"github.com/labstack/echo/v4"
)

func (s *Server) PostUsersSetIsActive(ctx echo.Context) error {
	var req struct {
		UserID   string `json:"user_id"`
		IsActive bool   `json:"is_active"`
	}

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]any{
			"error": map[string]string{
				"code":    "INVALID_REQUEST",
				"message": "invalid request body",
			},
		})
	}

	user, err := s.userRepo.FindByID(ctx.Request().Context(), entities.UserID(req.UserID))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]any{
			"error": map[string]string{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
	}

	if user == nil {
		return ctx.JSON(http.StatusNotFound, map[string]any{
			"error": map[string]string{
				"code":    "NOT_FOUND",
				"message": "user not found",
			},
		})
	}

	user.SetActive(req.IsActive)
	if err := s.userRepo.Update(ctx.Request().Context(), user); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]any{
			"error": map[string]string{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
	}

	var teamName string
	if user.TeamID() != nil {
		team, err := s.teamRepo.FindByID(ctx.Request().Context(), *user.TeamID())
		if err == nil && team != nil {
			teamName = team.Name()
		}
	}

	return ctx.JSON(http.StatusOK, map[string]any{
		"user": map[string]any{
			"user_id":   string(user.ID()),
			"username":  user.Username(),
			"team_name": teamName,
			"is_active": user.IsActive(),
		},
	})
}

func (s *Server) GetUsersGetReview(ctx echo.Context) error {
	userID := ctx.QueryParam("user_id")
	if userID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]any{
			"error": map[string]string{
				"code":    "INVALID_REQUEST",
				"message": "user_id is required",
			},
		})
	}

	prs, err := s.prRepo.FindPullRequestByUserID(
		ctx.Request().Context(),
		entities.UserID(userID),
	)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]any{
			"error": map[string]string{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
	}

	pullRequests := make([]map[string]any, len(prs))
	for i, pr := range prs {
		pullRequests[i] = map[string]any{
			"pull_request_id":   string(pr.ID()),
			"pull_request_name": pr.Name(),
			"author_id":         string(pr.AuthorID()),
			"status":            pr.Status().String(),
		}
	}

	return ctx.JSON(http.StatusOK, map[string]any{
		"user_id":       userID,
		"pull_requests": pullRequests,
	})
}
