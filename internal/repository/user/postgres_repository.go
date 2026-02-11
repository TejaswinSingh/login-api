package user

import (
	"context"
	"errors"
	"log/slog"

	"github.com/TejaswinSingh/login-api/internal/config"
	"github.com/TejaswinSingh/login-api/internal/constants"
	common_errors "github.com/TejaswinSingh/login-api/internal/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type postgresUserRepository struct {
	dbPool *pgxpool.Pool
	logger *slog.Logger
	config config.Config
}

func NewPostgresUserRepository(dbPool *pgxpool.Pool, logger *slog.Logger, config config.Config) UserRepository {
	return &postgresUserRepository{
		dbPool: dbPool,
		logger: logger,
		config: config,
	}
}

func (r *postgresUserRepository) GetUserFromID(ctx context.Context, id int) (*User, error) {
	query := "SELECT * FROM public.users WHERE id = $1"

	rows, _ := r.dbPool.Query(ctx, query, id)
	user, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.ErrorContext(ctx, "user not found")
			return nil, common_errors.ErrNotFound
		}
		r.logger.ErrorContext(ctx, "unable to get user", "error", err.Error())
		return nil, common_errors.ErrInternalServer
	}
	return user, nil
}

func (r *postgresUserRepository) CreateNewUser(ctx context.Context, user *User) error {
	query := "INSERT INTO public.users (username, passhash) VALUES (@username, @passhash)"
	args := pgx.NamedArgs{
		"username": user.Username,
		"passhash": user.Passhash,
	}

	var pgerr *pgconn.PgError

	if _, err := r.dbPool.Exec(ctx, query, args); err != nil {
		if errors.As(err, &pgerr) && pgerr.Code == string(constants.UNIQUE_CONSTRAINT_VIOLATION) {
			r.logger.ErrorContext(ctx, "user already exists", "error", err.Error())
			return common_errors.ErrResourceAlreadyExists
		}
		r.logger.ErrorContext(ctx, "unable to create user", "error", err.Error())
		return common_errors.ErrInternalServer
	}

	return nil
}

func (r *postgresUserRepository) GeneratePasswordHash(ctx context.Context, password string) (string, error) {
	passhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		r.logger.ErrorContext(ctx, "unable to generate password hash", "error", err.Error())
		return "", common_errors.ErrInternalServer
	}
	return string(passhash), nil
}

func (r *postgresUserRepository) ValidatePassword(ctx context.Context, username, password string) error {
	query := "SELECT passhash FROM public.users WHERE username = @username"
	args := pgx.NamedArgs{
		"username": username,
	}

	var dbhash string
	err := r.dbPool.QueryRow(ctx, query, args).Scan(&dbhash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.ErrorContext(ctx, "user not found")
			return common_errors.ErrNotFound
		}
		r.logger.ErrorContext(ctx, "unable to get user passhash", "error", err.Error())
		return common_errors.ErrInternalServer
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbhash), []byte(password)); err != nil {
		r.logger.ErrorContext(ctx, "password validation failed", "error", err.Error())
		return common_errors.ErrInvalidCreds
	}

	return nil
}
