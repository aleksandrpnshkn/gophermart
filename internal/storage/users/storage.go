package users

import (
	"context"
	"errors"

	"github.com/aleksandrpnshkn/gophermart/internal/models"
	"github.com/aleksandrpnshkn/gophermart/internal/types"
)

type Storage interface {
	Ping(ctx context.Context) error

	GetByID(ctx context.Context, id int64) (models.User, error)

	GetByLogin(ctx context.Context, login string) (models.User, error)

	Create(ctx context.Context, login string, hash types.PasswordHash) (models.User, error)

	Close() error
}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)
