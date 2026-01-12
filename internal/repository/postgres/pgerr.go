package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const pgUniqueViolation = "23505"

func IsUniqueViolation(err error, constraintName string) bool {
	var pgerr *pgconn.PgError
	if errors.As(err, &pgerr) && pgerr != nil && pgerr.Code == pgUniqueViolation {
		if constraintName == "" {
			return true
		}
		return pgerr.ConstraintName == constraintName
	}
	return false
}
