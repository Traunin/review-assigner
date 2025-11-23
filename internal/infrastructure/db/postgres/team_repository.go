package postgres

import (
	"context"

	"github.com/Traunin/review-assigner/internal/domain/entities"
	"github.com/Traunin/review-assigner/internal/infrastructure/db/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type TeamRepository struct {
	db *DB
}

func NewTeamRepository(db *DB) *TeamRepository {
	return &TeamRepository{
		db: db,
	}
}

func teamIdToPgInt4(id entities.TeamID) pgtype.Int4 {
	return pgtype.Int4{
		Int32: int32(id),
		Valid: true,
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

	// a user can only be a member of one team
	memberRows, err := r.db.Queries.GetUsersByTeamID(
		ctx,
		teamIdToPgInt4(team.ID()),
	)
	if err != nil {
		return nil, err
	}

	for _, row := range memberRows {
		if err := team.AddMember(entities.UserID(row.UserID)); err != nil {
			return nil, err
		}
	}

	return team, nil
}

func (r *TeamRepository) Create(
	ctx context.Context,
	team *entities.Team,
) error {
	_, err := r.db.Queries.CreateTeam(ctx, team.Name())
	return err
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

func (r *TeamRepository) FindByUserID(
	ctx context.Context,
	id entities.UserID,
) (*entities.Team, error) {
	teamRow, err := r.db.Queries.GetTeamByUserID(ctx, id.String())
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

func (r *TeamRepository) Update(
	ctx context.Context,
	team *entities.Team,
) error {
	return r.db.Queries.UpdateTeam(ctx, sqlc.UpdateTeamParams{
		ID:       int32(team.ID()),
		TeamName: team.Name(),
	})
}

func (r *TeamRepository) FindActiveReviewersByTeamID(
    ctx context.Context,
    id entities.TeamID,
) ([]*entities.User, error) {
    userRows, err := r.db.Queries.GetActiveUsersByTeamID(
        ctx,
        teamIdToPgInt4(id),
    )
    if err != nil {
        return nil, err
    }
    users := make([]*entities.User, 0, len(userRows))
    for _, row := range userRows {
        var teamID *entities.TeamID
        if row.TeamID.Valid {
            tid := entities.TeamID(row.TeamID.Int32)
            teamID = &tid
        }

        user, err := entities.NewUser(
            entities.UserID(row.UserID),
            row.Username,
            row.IsActive,
            teamID,
        )
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    return users, nil
}

func (r *TeamRepository) TeamExists(
	ctx context.Context,
	name string,
) (bool, error) {
	return r.db.Queries.TeamExists(ctx, name)
}
