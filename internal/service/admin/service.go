package admin

import (
	"log/slog"
	"time"

	pg "Beanefits/internal/repository/postgres"
)

type Service struct {
	db  pg.DBTX
	txm pg.TxManager

	users pg.UsersRepo
	roles pg.RolesRepo
	rules pg.RulesRepo

	now func() time.Time
	log *slog.Logger
}

type Deps struct {
	DB  pg.DBTX
	TXM pg.TxManager

	Users pg.UsersRepo
	Roles pg.RolesRepo
	Rules pg.RulesRepo

	Now func() time.Time
	Log *slog.Logger
}

func New(deps Deps) *Service {
	n := deps.Now
	if n == nil {
		n = time.Now
	}

	l := deps.Log
	if l == nil {
		l = slog.Default()
	}
	// Локальные теги сервиса. Глобальные теги (env/app/version/request_id) — лучше добавлять выше, в app/transport.
	l = l.With("layer", "service", "svc", "admin")

	return &Service{
		db:    deps.DB,
		txm:   deps.TXM,
		users: deps.Users,
		roles: deps.Roles,
		rules: deps.Rules,
		now:   n,
		log:   l,
	}
}
