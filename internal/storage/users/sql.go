package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aleksandrpnshkn/gophermart/internal/models"
	"github.com/aleksandrpnshkn/gophermart/internal/types"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type SQLStorage struct {
	pgxpool *pgxpool.Pool
}

func (s *SQLStorage) Ping(ctx context.Context) error {
	return s.pgxpool.Ping(ctx)
}

func (s *SQLStorage) GetByID(ctx context.Context, id int64) (models.User, error) {
	var user models.User

	row := s.pgxpool.QueryRow(ctx, `
        SELECT id, login, password_hash FROM users 
        WHERE id = $1
    `, id)
	err := row.Scan(&user.ID, &user.Login, &user.Hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, err
	}

	return user, nil
}

func (s *SQLStorage) GetByLogin(ctx context.Context, login string) (models.User, error) {
	var user models.User

	row := s.pgxpool.QueryRow(ctx, `
        SELECT id, login, password_hash FROM users 
        WHERE login = $1
    `, login)
	err := row.Scan(&user.ID, &user.Login, &user.Hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, err
	}

	return user, nil
}

func (s *SQLStorage) Create(
	ctx context.Context,
	login string,
	password types.PasswordHash,
) (models.User, error) {
	var user models.User

	row := s.pgxpool.QueryRow(ctx, `
        INSERT INTO users (id, login, password_hash) 
        VALUES (DEFAULT, @login, @password_hash) 
        RETURNING id, login, password_hash
    `, pgx.NamedArgs{
		"login":         login,
		"password_hash": password,
	})
	err := row.Scan(&user.ID, &user.Login, &user.Hash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return user, ErrUserAlreadyExists
		}

		return user, err
	}

	return user, nil
}

func (s *SQLStorage) Close() error {
	s.pgxpool.Close()
	return nil
}

func NewSQLStorage(ctx context.Context, databaseDSN string) (*SQLStorage, error) {
	pool, err := pgxpool.New(ctx, databaseDSN)
	if err != nil {
		return nil, err
	}

	storage := SQLStorage{
		pgxpool: pool,
	}

	err = storage.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return &storage, nil
}
