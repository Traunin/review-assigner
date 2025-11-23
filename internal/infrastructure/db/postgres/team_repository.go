package postgres

import (
	"context"

	"github.com/Traunin/review-assigner/internal/domain/entities"
	"github.com/Traunin/review-assigner/internal/infrastructure/db/sqlc"
	"github.com/jackc/pgx/v5"
)

type TeamRepository struct {
	db *DB
}

func NewTeamRepository(db *DB) *TeamRepository {
	return &TeamRepository{
		db: db,
	}
}

func (r *TeamRepository) buildTeamWithMembers(
	ctx context.Context,
	teamID int32,
	teamName string,
) (*entities.Team, error) {
	team, err := entities.NewTeam(teamName, entities.TeamID(teamID))
	if err != nil {
		return nil, err
	}

	memberRows, err := r.db.Queries.GetTeamWithMembers(ctx, teamName)
	if err != nil {
		return nil, err
	}

	for _, row := range memberRows {
		if row.UserID.Valid {
			if err := team.AddMember(entities.UserID(row.UserID.String)); err != nil {
				return nil, err
			}
		}
	}

	return team, nil
}

func (r *TeamRepository) Create(
	ctx context.Context,
	team *entities.Team,
) error {
	return r.db.execTx(ctx, func(q *sqlc.Queries) error {
		if _, err := q.CreateTeam(ctx, team.Name()); err != nil {
			return err
		}

		for _, userID := range team.Members() {
			if err := q.AddTeamMember(ctx, sqlc.AddTeamMemberParams{
				TeamID: int32(team.ID()),
				UserID: userID.String(),
			}); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *TeamRepository) DeleteByID(
	ctx context.Context,
	id entities.TeamID,
) error {
	return r.db.Queries.DeleteTeam(ctx, int32(id))
}

func (r *TeamRepository) FindByID(
	ctx context.Context,
	id entities.TeamID,
) (*entities.Team, error) {
	teamRow, err := r.db.Queries.GetTeamByID(ctx, int32(id))
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return r.buildTeamWithMembers(ctx, teamRow.ID, teamRow.TeamName)
}

func (r *TeamRepository) FindByName(
	ctx context.Context,
	name string,
) (*entities.Team, error) {
	teamRow, err := r.db.Queries.GetTeamByName(ctx, name)
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return r.buildTeamWithMembers(ctx, teamRow.ID, teamRow.TeamName)
}

func (r *TeamRepository) FindAll(
	ctx context.Context,
) ([]*entities.Team, error) {
	teams, err := r.db.Queries.GetTeams(ctx)
	if err != nil {
		return nil, err
	}

	teamEntities := make([]*entities.Team, len(teams))
	for i, team := range teams {
		teamEntity, err := r.buildTeamWithMembers(ctx, team.ID, team.TeamName)
		if err != nil {
			return nil, err
		}
		teamEntities[i] = teamEntity
	}

	return teamEntities, nil
}

func (r *TeamRepository) FindTeamsByUserID(
	ctx context.Context,
	id entities.UserID,
) ([]*entities.Team, error) {
	teamRows, err := r.db.Queries.GetTeamsByUserID(ctx, id.String())
	if err != nil {
		return nil, err
	}

	teams := make([]*entities.Team, len(teamRows))
	for i, teamRow := range teamRows {
		team, err := r.buildTeamWithMembers(ctx, teamRow.ID, teamRow.TeamName)
		if err != nil {
			return nil, err
		}
		teams[i] = team
	}

	return teams, nil
}

func (r *TeamRepository) Update(
	ctx context.Context,
	team *entities.Team,
) error {
	return r.db.execTx(ctx, func(q *sqlc.Queries) error {
		if err := q.UpdateTeam(ctx, sqlc.UpdateTeamParams{
			ID:       int32(team.ID()),
			TeamName: team.Name(),
		}); err != nil {
			return err
		}

		if err := q.DeleteTeamMembers(ctx, int32(team.ID())); err != nil {
			return err
		}

		for _, userID := range team.Members() {
			if err := q.AddTeamMember(ctx, sqlc.AddTeamMemberParams{
				TeamID: int32(team.ID()),
				UserID: userID.String(),
			}); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *TeamRepository) FindActiveReviewersByTeamID(
	ctx context.Context,
	id entities.TeamID,
) ([]*entities.User, error) {
	userRows, err := r.db.Queries.GetActiveTeamMembers(ctx, int32(id))
	if err != nil {
		return nil, err
	}

	users := make([]*entities.User, 0, len(userRows))
	for _, row := range userRows {
		user, err := entities.NewUser(
			entities.UserID(row.UserID),
			row.Username,
			row.IsActive,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
