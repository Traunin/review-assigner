package handlers

import (
	"net/http"

	"github.com/Traunin/review-assigner/internal/application/dto"
	"github.com/Traunin/review-assigner/internal/application/services"
	"github.com/Traunin/review-assigner/internal/domain/entities"
	"github.com/labstack/echo/v4"
)

func (s *Server) PostTeamAdd(ctx echo.Context) error {
	var req struct {
		TeamName string `json:"team_name"`
		Members  []struct {
			UserID   string `json:"user_id"`
			Username string `json:"username"`
			IsActive bool   `json:"is_active"`
		} `json:"members"`
	}

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": map[string]string{
				"code":    "INVALID_REQUEST",
				"message": "invalid request body",
			},
		})
	}

	cmd := dto.CreateTeamCmd{
		TeamName: req.TeamName,
		Members:  make([]dto.TeamMemberCmd, len(req.Members)),
	}

	for i, m := range req.Members {
		cmd.Members[i] = dto.TeamMemberCmd{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		}
	}

	team, err := s.teamService.CreateTeam(ctx.Request().Context(), cmd)
	if err != nil {
		if err == services.ErrTeamExists {
			return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": map[string]string{
					"code":    "TEAM_EXISTS",
					"message": "team_name already exists",
				},
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": map[string]string{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
	}

	users, err := s.userRepo.GetByTeamID(
		ctx.Request().Context(),
		entities.TeamID(team.ID),
	)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": map[string]string{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
	}

	members := make([]map[string]interface{}, len(users))
	for i, u := range users {
		members[i] = map[string]interface{}{
			"user_id":   string(u.ID()),
			"username":  u.Username(),
			"is_active": u.IsActive(),
		}
	}

	return ctx.JSON(http.StatusCreated, map[string]interface{}{
		"team": map[string]interface{}{
			"team_name": team.TeamName,
			"members":   members,
		},
	})
}

func (s *Server) GetTeamGet(ctx echo.Context) error {
	teamName := ctx.QueryParam("team_name")
	if teamName == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": map[string]string{
				"code":    "INVALID_REQUEST",
				"message": "team_name is required",
			},
		})
	}

	team, err := s.teamService.GetTeam(ctx.Request().Context(), teamName)
	if err != nil {
		if err == services.ErrTeamNotFound {
			return ctx.JSON(http.StatusNotFound, map[string]interface{}{
				"error": map[string]string{
					"code":    "NOT_FOUND",
					"message": "team not found",
				},
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": map[string]string{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
	}

	users, err := s.userRepo.GetByTeamID(
		ctx.Request().Context(),
		entities.TeamID(team.ID),
	)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": map[string]string{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
	}

	members := make([]map[string]interface{}, len(users))
	for i, u := range users {
		members[i] = map[string]interface{}{
			"user_id":   string(u.ID()),
			"username":  u.Username(),
			"is_active": u.IsActive(),
		}
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"team_name": team.TeamName,
		"members":   members,
	})
}
