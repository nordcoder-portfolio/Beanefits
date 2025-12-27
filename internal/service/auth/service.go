package auth

import (
	"context"
	"log/slog"
	"time"

	"Beanefits/internal/domain/errs"
	pg "Beanefits/internal/repository/postgres"
	pgdto "Beanefits/internal/repository/postgres/dto"
)

type PasswordHasher interface {
	Hash(ctx context.Context, password string) (string, error)
	Compare(ctx context.Context, password, passwordHash string) (bool, error)
}

type TokenIssuer interface {
	Issue(ctx context.Context, user pgdto.UserWithRoles) (string, error)
}

type PublicCodeGenerator interface {
	New(ctx context.Context) (string, error)
}

type Clock func() time.Time

type Service struct {
	db  pg.DBTX
	txm pg.TxManager

	users    pg.UsersRepo
	roles    pg.RolesRepo
	accounts pg.AccountsRepo

	hasher  PasswordHasher
	issuer  TokenIssuer
	codeGen PublicCodeGenerator

	now Clock
	log *slog.Logger

	initialLevelCode  string
	publicCodeRetries int
}

type Deps struct {
	DB  pg.DBTX
	TXM pg.TxManager

	Users    pg.UsersRepo
	Roles    pg.RolesRepo
	Accounts pg.AccountsRepo

	Hasher  PasswordHasher
	Issuer  TokenIssuer
	CodeGen PublicCodeGenerator

	Now Clock
	Log *slog.Logger

	InitialLevelCode  string
	PublicCodeRetries int
}

const (
	defaultInitialLevelCode  = "Green Bean"
	defaultPublicCodeRetries = 8
)

func New(deps Deps) *Service {
	n := deps.Now
	if n == nil {
		n = time.Now
	}

	l := deps.Log
	if l == nil {
		l = slog.Default()
	}
	l = l.With("layer", "service", "svc", "auth")

	level := deps.InitialLevelCode
	if level == "" {
		level = defaultInitialLevelCode
	}

	retries := deps.PublicCodeRetries
	if retries <= 0 {
		retries = defaultPublicCodeRetries
	}

	// wiring invariants: fail fast if misconfigured
	if deps.DB == nil {
		panic("auth.New: deps.DB is nil")
	}
	if deps.TXM == nil {
		panic("auth.New: deps.TXM is nil")
	}
	if deps.Users == nil || deps.Roles == nil || deps.Accounts == nil {
		panic("auth.New: repos are nil")
	}
	if deps.Hasher == nil {
		panic("auth.New: deps.Hasher is nil")
	}
	if deps.Issuer == nil {
		panic("auth.New: deps.Issuer is nil")
	}
	if deps.CodeGen == nil {
		panic("auth.New: deps.CodeGen is nil")
	}

	return &Service{
		db:                deps.DB,
		txm:               deps.TXM,
		users:             deps.Users,
		roles:             deps.Roles,
		accounts:          deps.Accounts,
		hasher:            deps.Hasher,
		issuer:            deps.Issuer,
		codeGen:           deps.CodeGen,
		now:               n,
		log:               l,
		initialLevelCode:  level,
		publicCodeRetries: retries,
	}
}

func (s *Service) invalidCredentials() error {
	// единая точка истины (важно для безопасности: одинаковая ошибка для "нет юзера" и "не тот пароль")
	return errs.New(errs.CodeInvalidCredentials, "invalid credentials")
}
