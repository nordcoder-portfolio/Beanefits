package middleware

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"Beanefits/internal/api"
)

type Claims struct {
	UserID int64
	Roles  []string
}

type TokenVerifier interface {
	Verify(ctx context.Context, token string) (Claims, error)
}

type ctxKey int

const claimsKey ctxKey = 1

func WithClaims(ctx context.Context, c Claims) context.Context {
	return context.WithValue(ctx, claimsKey, c)
}

func ClaimsFromContext(ctx context.Context) (Claims, bool) {
	v := ctx.Value(claimsKey)
	c, ok := v.(Claims)
	return c, ok
}

// Auth — oapi-codegen middleware.
// Требует Bearer только для операций, где wrapper положил api.BearerAuthScopes в ctx.
func Auth(verifier TokenVerifier, log *slog.Logger) api.MiddlewareFunc {
	if log == nil {
		log = slog.Default()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// публичная операция (openapi security: [])
			if r.Context().Value(api.BearerAuthScopes) == nil {
				next.ServeHTTP(w, r)
				return
			}

			if verifier == nil {
				detail := "token verifier is not configured"
				writeProblem(w, r, http.StatusUnauthorized, "Unauthorized", &detail, ptr("AUTH_NOT_CONFIGURED"))
				return
			}

			token, ok := bearerToken(r)
			if !ok {
				detail := "missing or invalid Authorization header"
				writeProblem(w, r, http.StatusUnauthorized, "Unauthorized", &detail, ptr("UNAUTHORIZED"))
				return
			}

			claims, err := verifier.Verify(r.Context(), token)
			if err != nil {
				detail := err.Error()
				writeProblem(w, r, http.StatusUnauthorized, "Unauthorized", &detail, ptr("UNAUTHORIZED"))
				return
			}

			ctx := WithClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func bearerToken(r *http.Request) (string, bool) {
	v := r.Header.Get("Authorization")
	const prefix = "Bearer "
	if !strings.HasPrefix(v, prefix) {
		return "", false
	}
	tok := strings.TrimSpace(strings.TrimPrefix(v, prefix))
	return tok, tok != ""
}

func writeProblem(w http.ResponseWriter, r *http.Request, status int, title string, detail *string, code *string) {
	instance := ""
	if r != nil && r.URL != nil {
		instance = r.URL.Path
	}

	p := api.Problem{
		Type:     "about:blank",
		Title:    title,
		Status:   status,
		Detail:   detail,
		Code:     code,
		Instance: &instance,
	}

	w.Header().Set("Content-Type", "application/problem+json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(p)
}

func ptr[T any](v T) *T { return &v }
