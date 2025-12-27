package httpserver

import (
	"Beanefits/internal/api"
	apphandler "Beanefits/internal/http/handler"
	appmw "Beanefits/internal/http/middleware"
	"context"
	"errors"
	"github.com/go-chi/cors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	httpServer *http.Server
	log        *slog.Logger
	cfg        Config
}

type TokenVerifier = appmw.TokenVerifier

func New(cfg Config, h *apphandler.Handler, verifier TokenVerifier, log *slog.Logger) *Server {
	if log == nil {
		log = slog.Default()
	}

	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(cfg.RequestTimeout))

	r.Use(appmw.AccessLog(log))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"X-Request-Id"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	api.HandlerWithOptions(h, api.ChiServerOptions{
		BaseURL:    cfg.BaseURL,
		BaseRouter: r,
		Middlewares: []api.MiddlewareFunc{
			appmw.Auth(verifier, log),
		},
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			reqID := chimw.GetReqID(r.Context())
			log.Error("request validation/binding error", reqID)
			h.WriteServiceError(w, r, err)
		},
	})

	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &Server{httpServer: srv, log: log, cfg: cfg}
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		s.log.Info("http server starting", "addr", s.cfg.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		s.log.Info("shutdown requested by context")
	case <-stop:
		s.log.Info("shutdown requested by signal")
	case err := <-errCh:
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		return err
	}
	return <-errCh
}

func (s *Server) Close() error {
	return s.httpServer.Close()
}
