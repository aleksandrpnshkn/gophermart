package middlewares

import (
	"net/http"

	"github.com/aleksandrpnshkn/gophermart/internal/models"
	"github.com/aleksandrpnshkn/gophermart/internal/services"
	"github.com/aleksandrpnshkn/gophermart/internal/types"
	"go.uber.org/zap"
)

const AuthCookieName = "auth_token"

func NewAuthMiddleware(logger *zap.Logger, auther services.Auther) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ctx := req.Context()

			authCookie, err := req.Cookie(AuthCookieName)
			if err != nil && err != http.ErrNoCookie {
				logger.Error("unknown cookie error", zap.Error(err))
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			var user models.User

			if err == http.ErrNoCookie {
				res.WriteHeader(http.StatusUnauthorized)
				return
			}

			user, err = auther.ParseToken(ctx, types.RawToken(authCookie.Value))
			if err != nil {
				if err == services.ErrInvalidToken {
					res.WriteHeader(http.StatusUnauthorized)
					return
				} else {
					logger.Error("failed to parse token", zap.Error(err))
					res.WriteHeader(http.StatusInternalServerError)
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
