package user

import (
	"Beanefits/internal/domain/errs"
	"regexp"
)

type Phone string

var rePhone = regexp.MustCompile(`^\+?[1-9]\d{10,14}$`)

func ParsePhone(s string) (Phone, error) {
	if !rePhone.MatchString(s) {
		return "", ErrInvalidPhone
	}
	return Phone(s), nil
}

var (
	ErrInvalidPhone = errs.New(errs.CodeInvalidPhone, "invalid phone format")
)
