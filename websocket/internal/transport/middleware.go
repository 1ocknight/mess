package transport

import (
	"fmt"
	"net/http"

	"github.com/TATAROmangol/mess/shared/auth"
	"github.com/TATAROmangol/mess/shared/logger"
	"github.com/TATAROmangol/mess/shared/requestmeta"
	"github.com/TATAROmangol/mess/websocket/internal/ctxkey"
	"github.com/TATAROmangol/mess/websocket/internal/loglables"
	"github.com/google/uuid"
)

func LoggerMiddleware(lg logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := ctxkey.WithLogger(r.Context(), lg)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequestMetadataMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			meta := requestmeta.GetFromHTTPRequest(r)

			id := uuid.NewString()

			lg, err := ctxkey.ExtractLogger(r.Context())
			if err != nil {
				http.Error(w, fmt.Sprintf("extract logger: %v", err), http.StatusInternalServerError)
				return
			}

			lg = lg.With(loglables.RequestMetadata, *meta)
			lg = lg.With(loglables.RequestID, id)

			ctx := ctxkey.WithLogger(r.Context(), lg)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type bodyWriter struct {
	http.ResponseWriter
	body       []byte
	statusCode int
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return w.ResponseWriter.Write(b)
}

func (w *bodyWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func LogResponseMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bw := &bodyWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(bw, r)

			lg, err := ctxkey.ExtractLogger(r.Context())
			if err != nil {
				http.Error(w, fmt.Sprintf("extract logger: %v", err), http.StatusInternalServerError)
				return
			}

			lg = lg.With(loglables.StatusResponse, bw.statusCode)
			lg = lg.With(loglables.Response, string(bw.body))
			lg.Info("request completed")
		})
	}
}

func SubjectMiddleware(auth auth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")

			sub, err := auth.Verify(token)
			if err != nil {
				http.Error(w, fmt.Sprintf("verify token: %v", err), http.StatusUnauthorized)
				return
			}

			ctx := ctxkey.WithSubject(r.Context(), sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
