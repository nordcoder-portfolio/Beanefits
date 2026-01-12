package ledger

import (
	"Beanefits/internal/domain/errs"
)

var (
	ErrInvalidPoints = errs.New(errs.CodeInvalidPoints, "points must be positive (or non-negative where allowed)")
)
