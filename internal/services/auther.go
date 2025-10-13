package services

import (
	"context"
	"errors"

	"github.com/aleksandrpnshkn/gophermart/internal/models"
	"github.com/aleksandrpnshkn/gophermart/internal/storage/users"
	"github.com/aleksandrpnshkn/gophermart/internal/types"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int64
}

type Auther interface {
	ParseToken(ctx context.Context, token types.RawToken) (models.User, error)

	RegisterUser(ctx context.Context, login string, password string) (models.User, types.RawToken, error)

	LoginUser(ctx context.Context, login string, password string) (models.User, types.RawToken, error)

	FromUserContext(ctx context.Context) (models.User, error)
}

type JwtAuther struct {
	usersStorage users.Storage

	secretKey string
}

type ctxKey string

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrBadCredentials     = errors.New("bad credentials")
	ErrLoginAlreadyExists = errors.New("login already exists")
)

const ctxUserID ctxKey = "user_id"

func (a *JwtAuther) ParseToken(ctx context.Context, tokenString types.RawToken) (models.User, error) {
	token, err := jwt.ParseWithClaims(string(tokenString), &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(a.secretKey), nil
	})
	if err != nil {
		return models.User{}, ErrInvalidToken
	}

	claims := token.Claims.(*Claims)

	user, err := a.usersStorage.GetByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			return models.User{}, ErrInvalidToken
		} else {
			return models.User{}, err
		}
	}

	return user, nil
}

func (a *JwtAuther) RegisterUser(
	ctx context.Context,
	login string,
	password string,
) (models.User, types.RawToken, error) {
	token := types.RawToken("")
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, token, err
	}
	hash := types.PasswordHash(hashBytes)

	user, err := a.usersStorage.Create(ctx, login, hash)
	if err != nil {
		if errors.Is(err, users.ErrUserAlreadyExists) {
			return models.User{}, token, ErrLoginAlreadyExists
		}

		return models.User{}, token, err
	}

	token, err = a.createAuthToken(user.ID)
	if err != nil {
		return models.User{}, token, err
	}

	return user, token, nil
}

func (a *JwtAuther) LoginUser(
	ctx context.Context,
	login string,
	password string,
) (models.User, types.RawToken, error) {
	token := types.RawToken("")
	user, err := a.usersStorage.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			return models.User{}, token, ErrBadCredentials
		}

		return models.User{}, token, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
	if err != nil {
		return models.User{}, token, ErrBadCredentials
	}

	token, err = a.createAuthToken(user.ID)
	if err != nil {
		return models.User{}, token, err
	}

	return user, token, nil
}

func (a *JwtAuther) FromUserContext(ctx context.Context) (models.User, error) {
	userID, ok := ctx.Value(ctxUserID).(int64)
	if !ok {
		return models.User{}, errors.New("user id is not set")
	}

	user, err := a.usersStorage.GetByID(ctx, userID)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (a *JwtAuther) createAuthToken(userID int64) (types.RawToken, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{},
		UserID:           userID,
	})

	tokenString, err := token.SignedString([]byte(a.secretKey))
	if err != nil {
		return types.RawToken(""), err
	}

	return types.RawToken(tokenString), nil
}

func NewAuther(usersStorage users.Storage, secretKey string) *JwtAuther {
	return &JwtAuther{
		usersStorage: usersStorage,
		secretKey:    secretKey,
	}
}

func NewUserContext(ctx context.Context, user models.User) context.Context {
	return context.WithValue(ctx, ctxUserID, user.ID)
}
