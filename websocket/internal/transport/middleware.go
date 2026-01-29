package transport

import (
	"fmt"
	"net/http"

	"github.com/1ocknight/mess/shared/auth"
	"github.com/1ocknight/mess/shared/logger"
	"github.com/1ocknight/mess/websocket/internal/ctxkey"
)

const (
	Bearer = "Bearer"
)

func SubjectMiddleware(auth auth.Service, lg logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.URL.Query().Get("token")
			if token == "" {
				err := fmt.Errorf("not found token")
				lg.Error(err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			token = fmt.Sprintf("%v %v", Bearer, token)

			sub, err := auth.Verify(token)
			if err != nil {
				err = fmt.Errorf("verify token: %v", err)
				lg.Error(err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := ctxkey.WithSubject(r.Context(), sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
