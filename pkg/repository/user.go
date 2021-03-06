package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aodin/date"
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
	"github.com/neoxelox/odin/pkg/model"
)

const USER_TABLE = "user"

var (
	ErrUserGeneric = internal.NewError("User query failed")
	ErrUserExists  = internal.NewError("User already exists")
)

type UserRepository struct {
	class.Repository
}

func NewUserRepository(configuration internal.Configuration, logger core.Logger, database database.Database) *UserRepository {
	return &UserRepository{
		Repository: *class.NewRepository(configuration, logger, database),
	}
}

func (self *UserRepository) Create(ctx context.Context, user model.User) (*model.User, error) {
	var u model.User

	query := fmt.Sprintf(`INSERT INTO "%s"
						  ("id", "phone", "name", "email", "picture", "birthday",
						  "language", "last_session_id", "is_banned", "created_at", "deleted_at")
						  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
						  RETURNING *;`, USER_TABLE)

	err := self.Database.Query(
		ctx, query, user.ID, user.Phone, user.Name, user.Email, user.Picture, user.Birthday,
		user.Language, user.LastSessionID, user.IsBanned, user.CreatedAt, user.DeletedAt).Scan(&u)
	switch {
	case err == nil:
		return &u, nil
	case database.ErrIntegrityViolation().Is(err):
		return nil, ErrUserExists().Wrap(err)
	default:
		return nil, ErrUserGeneric().Wrap(err)
	}
}

func (self *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var u model.User

	query := fmt.Sprintf(`SELECT * FROM "%s"
						  WHERE "id" = $1;`, USER_TABLE)

	err := self.Database.Query(ctx, query, id).Scan(&u)
	switch {
	case err == nil:
		return &u, nil
	case database.ErrNoRows().Is(err):
		return nil, nil
	default:
		return nil, ErrUserGeneric().Wrap(err)
	}
}

func (self *UserRepository) GetByIDs(ctx context.Context, ids []string) ([]model.User, error) {
	var us []model.User

	query := fmt.Sprintf(`SELECT * FROM "%s"
						  WHERE "id" = ANY ($1);`, USER_TABLE)

	err := self.Database.Query(ctx, query, ids).Scan(&us)
	switch {
	case err == nil:
		return us, nil
	case database.ErrNoRows().Is(err):
		return []model.User{}, nil
	default:
		return nil, ErrUserGeneric().Wrap(err)
	}
}

func (self *UserRepository) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	var u model.User

	query := fmt.Sprintf(`SELECT * FROM "%s"
						  WHERE "phone" = $1 AND "deleted_at" IS NULL;`, USER_TABLE)

	err := self.Database.Query(ctx, query, phone).Scan(&u)
	switch {
	case err == nil:
		return &u, nil
	case database.ErrNoRows().Is(err):
		return nil, nil
	default:
		return nil, ErrUserGeneric().Wrap(err)
	}
}

func (self *UserRepository) UpdateSession(ctx context.Context, id string, sessionID string) error {
	query := fmt.Sprintf(`UPDATE "%s"
						  SET "last_session_id" = $1
						  WHERE "id" = $2;`, USER_TABLE)

	affected, err := self.Database.Exec(ctx, query, sessionID, id)
	if err != nil {
		return ErrUserGeneric().Wrap(err)
	}

	if affected != 1 {
		return ErrUserGeneric()
	}

	return nil
}

func (self *UserRepository) UpdateProfile(ctx context.Context, id string, name string, picture string, birthday date.Date) error {
	query := fmt.Sprintf(`UPDATE "%s"
						  SET "name" = $1, "picture" = $2, "birthday" = $3
						  WHERE "id" = $4;`, USER_TABLE)

	affected, err := self.Database.Exec(ctx, query, name, picture, birthday, id)
	if err != nil {
		return ErrUserGeneric().Wrap(err)
	}

	if affected != 1 {
		return ErrUserGeneric()
	}

	return nil
}

func (self *UserRepository) UpdateEmail(ctx context.Context, id string, email string) error {
	query := fmt.Sprintf(`UPDATE "%s"
						  SET "email" = $1
						  WHERE "id" = $2;`, USER_TABLE)

	affected, err := self.Database.Exec(ctx, query, email, id)
	if err != nil {
		return ErrUserGeneric().Wrap(err)
	}

	if affected != 1 {
		return ErrUserGeneric()
	}

	return nil
}

func (self *UserRepository) UpdatePhone(ctx context.Context, id string, phone string) error {
	query := fmt.Sprintf(`UPDATE "%s"
						  SET "phone" = $1
						  WHERE "id" = $2;`, USER_TABLE)

	affected, err := self.Database.Exec(ctx, query, phone, id)
	switch {
	case err == nil:
		if affected != 1 {
			return ErrUserGeneric()
		}

		return nil
	case database.ErrIntegrityViolation().Is(err):
		return ErrUserExists().Wrap(err)
	default:
		return ErrUserGeneric().Wrap(err)
	}
}

func (self *UserRepository) UpdateDeletedAt(ctx context.Context, id string, deletedAt time.Time) error {
	query := fmt.Sprintf(`UPDATE "%s"
						  SET "deleted_at" = $1
						  WHERE "id" = $2;`, USER_TABLE)

	affected, err := self.Database.Exec(ctx, query, deletedAt, id)
	if err != nil {
		return ErrUserGeneric().Wrap(err)
	}

	if affected != 1 {
		return ErrUserGeneric()
	}

	return nil
}
