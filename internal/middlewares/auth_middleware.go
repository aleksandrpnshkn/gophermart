package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/aleksandrpnshkn/gophermart/internal/models"
	"github.com/aleksandrpnshkn/gophermart/internal/services"
	"github.com/aleksandrpnshkn/gophermart/internal/types"
	"go.uber.org/zap"
)

const AuthCookieName = "auth_token"

type TokenParser interface {
	ParseToken(ctx context.Context, token types.RawToken) (models.User, error)
}

func NewAuthMiddleware(
	responser *services.Responser,
	logger *zap.Logger,
	tokenParser TokenParser,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ctx := req.Context()

			authCookie, err := req.Cookie(AuthCookieName)
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					res.WriteHeader(http.StatusUnauthorized)
					return
				} else {
					logger.Error("unknown cookie error", zap.Error(err))
					responser.WriteInternalServerError(ctx, res)
					return
				}
			}

			var user models.User

			user, err = tokenParser.ParseToken(ctx, types.RawToken(authCookie.Value))
			if err != nil {
				if errors.Is(err, services.ErrInvalidToken) {
					res.WriteHeader(http.StatusUnauthorized)
					return
				} else {
					logger.Error("failed to parse token", zap.Error(err))
					responser.WriteInternalServerError(ctx, res)
					return
				}
			}

			req = req.WithContext(services.NewUserContext(ctx, user))

			next.ServeHTTP(res, req)
		})
	}
}

func SetAuthCookie(res http.ResponseWriter, token types.RawToken) {
	authCookie := &http.Cookie{
		Name:  AuthCookieName,
		Value: string(token),

		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   false,
	}
	http.SetCookie(res, authCookie)
}
