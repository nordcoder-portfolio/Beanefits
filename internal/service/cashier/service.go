package cashier

import (
	"log/slog"
	"time"

	pg "Beanefits/internal/repository/postgres"
)

type Clock func() time.Time

type Service struct {
	db  pg.DBTX
	txm pg.TxManager

	accounts   pg.AccountsRepo
	events     pg.EventsRepo
	operations pg.OperationsRepo
	rules      pg.RulesRepo

	now Clock
	log *slog.Logger
}

type Deps struct {
	DB  pg.DBTX
	TXM pg.TxManager

	Accounts   pg.AccountsRepo
	Events     pg.EventsRepo
	Operations pg.OperationsRepo
	Rules      pg.RulesRepo

	Now Clock
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
	l = l.With("layer", "service", "svc", "cashier")

	return &Service{
		db:         deps.DB,
		txm:        deps.TXM,
		accounts:   deps.Accounts,
		events:     deps.Events,
		operations: deps.Operations,
		rules:      deps.Rules,
		now:        n,
		log:        l,
	}
}
