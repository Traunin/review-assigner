package postgres

import (
	"context"

	"github.com/Traunin/review-assigner/internal/domain/entities"
	"github.com/Traunin/review-assigner/internal/infrastructure/db/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(
	ctx context.Context,
	user *entities.User,
) error {
	var pgTeamID pgtype.Int4
	if user.TeamID() != nil {
		pgTeamID = pgtype.Int4{Int32: int32(*user.TeamID()), Valid: true}
	}

	_, err := r.db.Queries.CreateUser(ctx, sqlc.CreateUserParams{
		UserID:   user.ID().String(),
		Username: user.Username(),
		IsActive: user.IsActive(),
		TeamID:   pgTeamID,
	})

	return err
}

func (r *UserRepository) DeleteByID(
	ctx context.Context,
	id entities.UserID,
) error {
	return r.db.Queries.DeleteUser(ctx, id.String())
}

func (r *UserRepository) FindByID(
	ctx context.Context,
	id entities.UserID,
) (*entities.User, error) {
	user, err := r.db.Queries.GetUserByID(ctx, id.String())
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var teamID *entities.TeamID
	if user.TeamID.Valid {
		tid := entities.TeamID(user.TeamID.Int32)
		teamID = &tid
	}

	return entities.NewUser(
		entities.UserID(user.UserID),
		user.Username,
		user.IsActive,
		teamID,
	)
}

func (r *UserRepository) FindAll(
	ctx context.Context,
) ([]*entities.User, error) {
	users, err := r.db.Queries.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	userEntities := make([]*entities.User, len(users))
	for i, user := range users {
		var teamID *entities.TeamID
		if user.TeamID.Valid {
			tid := entities.TeamID(user.TeamID.Int32)
			teamID = &tid
		}

		userEntity, err := entities.NewUser(
			entities.UserID(user.UserID),
			user.Username,
			user.IsActive,
			teamID,
		)
		if err != nil {
			return nil, err
		}
		userEntities[i] = userEntity
	}

	return userEntities, nil
}

func (r *UserRepository) Update(
	ctx context.Context,
	user *entities.User,
) error {
	var pgTeamID pgtype.Int4
	if user.TeamID() != nil {
		pgTeamID = pgtype.Int4{Int32: int32(*user.TeamID()), Valid: true}
	}

	return r.db.Queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		UserID:   user.ID().String(),
		Username: user.Username(),
		IsActive: user.IsActive(),
		TeamID:   pgTeamID,
	})
}

func (r *UserRepository) GetActiveUsers(
	ctx context.Context,
) ([]*entities.User, error) {
	users, err := r.db.Queries.GetActiveUsers(ctx)
	if err != nil {
		return nil, err
	}

	userEntities := make([]*entities.User, len(users))
	for i, user := range users {
		var teamID *entities.TeamID
		if user.TeamID.Valid {
			tid := entities.TeamID(user.TeamID.Int32)
			teamID = &tid
		}

		userEntity, err := entities.NewUser(
			entities.UserID(user.UserID),
			user.Username,
			user.IsActive,
			teamID,
		)
		if err != nil {
			return nil, err
		}
		userEntities[i] = userEntity
	}

	return userEntities, nil
}
