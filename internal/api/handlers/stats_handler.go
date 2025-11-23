package handlers

import (
	"net/http"

	"github.com/Traunin/review-assigner/internal/domain/entities"
	"github.com/labstack/echo/v4"
)

func (s *Server) GetStatsReviewers(ctx echo.Context) error {
	userID := ctx.QueryParam("user_id")

	// If user_id is provided, return stats for specific user
	if userID != "" {
		user, err := s.userRepo.FindByID(ctx.Request().Context(), entities.UserID(userID))
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

		openCount := 0
		mergedCount := 0
		for _, pr := range prs {
			if pr.Status().String() == "OPEN" {
				openCount++
			} else if pr.Status().String() == "MERGED" {
				mergedCount++
			}
		}

		return ctx.JSON(http.StatusOK, map[string]any{
			"user_id":            userID,
			"username":           user.Username(),
			"total_assignments":  len(prs),
			"open_assignments":   openCount,
			"merged_assignments": mergedCount,
		})
	}

	// If no user_id, return stats for all users
	allPRs, err := s.prRepo.FindAll(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]any{
			"error": map[string]string{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
	}

	// Count assignments per user
	userStats := make(map[entities.UserID]*struct {
		Username string
		Total    int
		Open     int
		Merged   int
	})

	for _, pr := range allPRs {
		reviewers := pr.Reviewers()
		for _, reviewer := range reviewers {
			if _, exists := userStats[reviewer.UserID]; !exists {
				user, err := s.userRepo.FindByID(ctx.Request().Context(), reviewer.UserID)
				if err != nil || user == nil {
					continue
				}
				userStats[reviewer.UserID] = &struct {
					Username string
					Total    int
					Open     int
					Merged   int
				}{
					Username: user.Username(),
				}
			}

			userStats[reviewer.UserID].Total++
			if pr.Status().String() == "OPEN" {
				userStats[reviewer.UserID].Open++
			} else if pr.Status().String() == "MERGED" {
				userStats[reviewer.UserID].Merged++
			}
		}
	}

	users := make([]map[string]any, 0, len(userStats))
	for userID, stats := range userStats {
		users = append(users, map[string]any{
			"user_id":            string(userID),
			"username":           stats.Username,
			"total_assignments":  stats.Total,
			"open_assignments":   stats.Open,
			"merged_assignments": stats.Merged,
		})
	}

	return ctx.JSON(http.StatusOK, map[string]any{
		"users":               users,
		"total_pull_requests": len(allPRs),
	})
}

// GetStatsPullRequests returns statistics about pull requests
func (s *Server) GetStatsPullRequests(ctx echo.Context) error {
	allPRs, err := s.prRepo.FindAll(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]any{
			"error": map[string]string{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
	}

	openCount := 0
	mergedCount := 0
	totalReviewers := 0

	for _, pr := range allPRs {
		if pr.Status().String() == "OPEN" {
			openCount++
		} else if pr.Status().String() == "MERGED" {
			mergedCount++
		}
		totalReviewers += len(pr.Reviewers())
	}

	avgReviewers := 0.0
	if len(allPRs) > 0 {
		avgReviewers = float64(totalReviewers) / float64(len(allPRs))
	}

	return ctx.JSON(http.StatusOK, map[string]any{
		"total_pull_requests":  len(allPRs),
		"open_pull_requests":   openCount,
		"merged_pull_requests": mergedCount,
		"total_reviewers":      totalReviewers,
		"avg_reviewers_per_pr": avgReviewers,
	})
}
