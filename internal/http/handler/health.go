package handler

import (
	"context"
	"net/http"
	"time"

	"Beanefits/internal/api"
)

// GET /healthz
func (h *Handler) GetHealthz(w http.ResponseWriter, r *http.Request) {
	h.helpers.JSON(w, http.StatusOK, api.HealthResponse{Status: "ok"})
}

// GET /readyz
func (h *Handler) GetReadyz(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		detail := "db is not configured"
		WriteProblem(w, http.StatusServiceUnavailable, "Not ready", &detail, ptr("DB_NOT_CONFIGURED"), instanceFromRequest(r))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.db.Ping(ctx); err != nil {
		detail := err.Error()
		WriteProblem(w, http.StatusServiceUnavailable, "Not ready", &detail, ptr("DB_UNAVAILABLE"), instanceFromRequest(r))
		return
	}

	h.helpers.JSON(w, http.StatusOK, api.HealthResponse{Status: "ok"})
}

// GET /metrics
func (h *Handler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	if h.metrics == nil {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte{})
		return
	}
	h.metrics.ServeHTTP(w, r)
}
