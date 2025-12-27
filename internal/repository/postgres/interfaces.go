package postgres

import (
	"context"
	"time"

	"Beanefits/internal/repository/postgres/dto"
)

type UsersRepo interface {
	Create(ctx context.Context, db DBTX, phone, passwordHash string) (userID int64, err error)
	GetByPhone(ctx context.Context, db DBTX, phone string) (dto.UserRow, bool, error)
	GetByID(ctx context.Context, db DBTX, id int64) (dto.UserRow, bool, error)
	List(ctx context.Context, db DBTX, q string, limit, offset int) ([]dto.UserRow, error)
	Deactivate(ctx context.Context, db DBTX, id int64) error
}

type RolesRepo interface {
	GetRoles(ctx context.Context, db DBTX, userID int64) ([]dto.RoleCode, error)
	AddRole(ctx context.Context, db DBTX, userID int64, role dto.RoleCode) error
	RemoveRole(ctx context.Context, db DBTX, userID int64, role dto.RoleCode) error
}

type AccountsRepo interface {
	CreateForUser(ctx context.Context, db DBTX, userID int64, publicCode string, initialLevelCode string) (dto.AccountRow, error)

	GetByUserID(ctx context.Context, db DBTX, userID int64) (dto.AccountRow, bool, error)
	GetByPublicCode(ctx context.Context, db DBTX, publicCode string) (dto.AccountRow, bool, error)

	// LockByID must use SELECT ... FOR UPDATE to protect concurrent spend/earn.
	LockByID(ctx context.Context, db DBTX, accountID int64) (dto.AccountRow, error)

	UpdateAfterEarn(ctx context.Context, db DBTX, accountID int64, balancePoints int, totalSpend dto.Money, levelCode string) (dto.AccountRow, error)
	UpdateAfterSpend(ctx context.Context, db DBTX, accountID int64, balancePoints int) (dto.AccountRow, error)
}

type RulesRepo interface {
	// GetEffectiveAt returns the ruleset effective at the provided timestamp (effective_from <= at, newest).
	GetEffectiveAt(ctx context.Context, db DBTX, at time.Time) (dto.RulesetWithLevels, bool, error)

	CreateRuleset(ctx context.Context, db DBTX, effectiveFrom time.Time, baseRubPerPoint dto.Money, levels []dto.LevelRuleRow) (dto.RulesetWithLevels, error)
	ListRulesets(ctx context.Context, db DBTX, limit, offset int) ([]dto.RulesetWithLevels, error)
}

type EventsRepo interface {
	Insert(ctx context.Context, db DBTX, in dto.EventInsert) (dto.EventRow, error)

	// ListByAccount returns newest-first; beforeTs is optional for pagination.
	ListByAccount(ctx context.Context, db DBTX, accountID int64, limit int, beforeTs *time.Time) ([]dto.EventRow, error)
}

type OperationsRepo interface {
	// Get returns cached operation (for idempotency replay).
	Get(ctx context.Context, db DBTX, accountID int64, opType dto.OperationType, operationID string) (dto.OperationRecord, bool, error)

	// InsertPending tries to insert a new operation record.
	// It must be implemented as INSERT ... ON CONFLICT DO NOTHING (or equivalent),
	// and returns inserted=true if a new row was created.
	InsertPending(ctx context.Context, db DBTX, in dto.OperationPendingInsert) (inserted bool, err error)

	// Finalize stores the HTTP status and response JSON for replay.
	Finalize(ctx context.Context, db DBTX, in dto.OperationFinalize) error
}
