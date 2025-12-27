package app

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log/slog"
	"time"

	"Beanefits/internal/config"
	"Beanefits/internal/http/handler"
	"Beanefits/internal/http/httpserver"
	"Beanefits/internal/infra/jwt"
	"Beanefits/internal/infra/jwtverifier"
	kafkainfra "Beanefits/internal/infra/kafka"
	"Beanefits/internal/infra/publiccode"
	"Beanefits/internal/infra/security"
	"Beanefits/internal/repository/postgres"
	"Beanefits/internal/repository/postgres/repo"
	"Beanefits/internal/repository/postgres/sqlc/gen"
	"Beanefits/internal/service/admin"
	"Beanefits/internal/service/auth"
	"Beanefits/internal/service/cashier"
	"Beanefits/internal/service/client"
	"Beanefits/internal/service/validation"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type App struct {
	DB   *pgxpool.Pool
	HTTP *httpserver.Server

	// Важно: не храним *Producer/*AutoPublisher в структуре, только cleanup-функции
	stopKafka  func()
	closeKafka func() error
}

func New(ctx context.Context, cfg config.Config, log *slog.Logger) (*App, error) {

	sum := sha256.Sum256([]byte(cfg.JWTSecret))
	log.Info("jwt config",
		"issuer", cfg.JWTIssuer,
		"secret_len", len(cfg.JWTSecret),
		"secret_fp8", fmt.Sprintf("%x", sum[:8]),
	)

	l := baseLogger(log)

	l.InfoContext(ctx, "app.init start")

	v, err := newValidator()
	if err != nil {
		l.ErrorContext(ctx, "app.init validator failed", "err", err)
		return nil, err
	}
	_ = v // хендлеру, вероятно, пригодится позже; сейчас у тебя в handler.New валидация не прокидывается

	pool, err := postgres.NewPool(ctx, cfg.PostgresURL)
	if err != nil {
		l.ErrorContext(ctx, "app.init postgres pool failed", "err", err)
		return nil, err
	}

	q := gen.New()

	// repositories
	usersRepo := repo.NewUsersRepo(q)
	accountsRepo := repo.NewAccountsRepo(q)
	rolesRepo := repo.NewRolesRepo(q)
	rulesRepo := repo.NewRulesRepo(q)
	opsRepo := repo.NewOperationsRepo(q)
	eventsRepo := repo.NewEventsRepo(q)

	txm := postgres.NewTxManager(pool)

	// infra
	hasher := security.NewHasher()
	jwtIssuer := jwt.NewIssuer(cfg.JWTSecret, cfg.JWTIssuer)
	codeGen := publiccode.NewGenerator()
	verifier := jwtverifier.New(cfg.JWTSecret, cfg.JWTIssuer)

	// time source
	now := func() time.Time { return time.Now().UTC() }

	// services
	authSvc := auth.New(auth.Deps{
		DB:       pool,
		TXM:      txm,
		Users:    usersRepo,
		Roles:    rolesRepo,
		Accounts: accountsRepo,
		Hasher:   hasher,
		Issuer:   jwtIssuer,
		CodeGen:  codeGen,
		Now:      now,
		Log:      l,
	})

	clientSvc := client.New(client.Deps{
		DB:       pool,
		Users:    usersRepo,
		Roles:    rolesRepo,
		Accounts: accountsRepo,
		Events:   eventsRepo,
		Now:      now,
		Log:      l,
	})

	cashierSvc := cashier.New(cashier.Deps{
		DB:         pool,
		TXM:        txm,
		Accounts:   accountsRepo,
		Events:     eventsRepo,
		Operations: opsRepo,
		Rules:      rulesRepo,
		Now:        now,
		Log:        l,
	})

	adminSvc := admin.New(admin.Deps{
		DB:    pool,
		TXM:   txm,
		Users: usersRepo,
		Roles: rolesRepo,
		Rules: rulesRepo,
		Now:   now,
		Log:   l,
	})

	// http handler
	h := handler.New(handler.Deps{
		Auth:    authSvc,
		Client:  clientSvc,
		Cashier: cashierSvc,
		Admin:   adminSvc,
		DB:      pool,
		Metrics: promhttp.Handler(),
	})

	httpCfg := httpserver.DefaultConfig()
	httpCfg.Addr = cfg.HTTPAddr
	httpCfg.BaseURL = cfg.HTTPBaseURL
	httpCfg.ShutdownTimeout = cfg.ShutdownTimeout

	srv := httpserver.New(httpCfg, h, verifier, l.With("component", "http"))

	// ---- Kafka wiring (создаётся внутри app.New) ----
	prod, err := kafkainfra.NewProducer(kafkainfra.ProducerConfig{
		Brokers:      cfg.KafkaBrokers,
		Topic:        cfg.KafkaTopic,
		WriteTimeout: cfg.KafkaWriteTimeout,
		BatchSize:    cfg.KafkaBatchSize,
		BatchTimeout: cfg.KafkaBatchTimeout,
	})
	if err != nil {
		l.ErrorContext(ctx, "app.init kafka producer failed", "err", err)
		_ = srv.Close()
		pool.Close()
		return nil, err
	}

	app := &App{
		DB:         pool,
		HTTP:       srv,
		stopKafka:  func() {},
		closeKafka: prod.Close,
	}

	if cfg.KafkaAutoPublish {
		stop := kafkainfra.StartAutoPublish(ctx, prod, kafkainfra.AutoPublishConfig{
			Interval: cfg.KafkaPublishInterval,
			Min:      -100,
			Max:      100,
			Log:      l,
		})
		app.stopKafka = stop
	}

	l.InfoContext(ctx, "app.init ok", "httpAddr", httpCfg.Addr)

	return app, nil
}

func (a *App) Close() error {
	var first error

	if a.stopKafka != nil {
		a.stopKafka()
	}

	if a.closeKafka != nil {
		if err := a.closeKafka(); err != nil {
			first = err
		}
	}

	if a.HTTP != nil {
		if err := a.HTTP.Close(); err != nil && first == nil {
			first = err
		}
	}

	if a.DB != nil {
		a.DB.Close()
	}

	return first
}

func baseLogger(log *slog.Logger) *slog.Logger {
	if log == nil {
		log = slog.Default()
	}
	return log.With("app", "Beanefits", "layer", "app")
}

func newValidator() (*validator.Validate, error) {
	v := validator.New()
	if err := validation.RegisterValidations(v); err != nil {
		return nil, err
	}
	return v, nil
}
