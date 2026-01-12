package handler

import (
	"net/http"

	"Beanefits/internal/api"
	"Beanefits/internal/http/middleware"
)

func (h *Handler) requireClient(w http.ResponseWriter, r *http.Request) (int64, bool) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		detail := "missing auth context"
		WriteProblem(w, http.StatusUnauthorized, "Unauthorized", &detail, ptr("UNAUTHORIZED"), instanceFromRequest(r))
		return 0, false
	}

	if !hasRole(claims.Roles, string(api.CLIENT)) {
		detail := "client role required"
		WriteProblem(w, http.StatusForbidden, "Forbidden", &detail, ptr("FORBIDDEN"), instanceFromRequest(r))
		return 0, false
	}

	return claims.UserID, true
}
