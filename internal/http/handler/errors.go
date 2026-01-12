package handler

import (
	"errors"
	"net/http"

	"Beanefits/internal/domain/errs"
)

type domainCoder interface{ Code() errs.Code }

func domainCode(err error) (errs.Code, bool) {
	if err == nil {
		return "", false
	}

	var dc domainCoder
	if errors.As(err, &dc) {
		c := dc.Code()
		if c != "" {
			return c, true
		}
	}

	type stringCoder interface{ Code() string }
	var sc stringCoder
	if errors.As(err, &sc) {
		c := sc.Code()
		if c != "" {
			return errs.Code(c), true
		}
	}

	return "", false
}

type problemSpec struct {
	status int
	title  string
}

func problemForCode(code errs.Code) (problemSpec, bool) {
	switch code {

	// 422
	case errs.CodeInvalidPhone,
		errs.CodeInvalidPublicCode,
		errs.CodeInvalidPoints,
		errs.CodeInvalidPurchaseAmount,
		errs.CodeInvalidRuleset,
		errs.CodeInvalidLevels,
		errs.CodeInvalidMoney:
		return problemSpec{status: http.StatusUnprocessableEntity, title: "Validation error"}, true

	// 409
	case errs.CodeNotEnoughBalance:
		return problemSpec{status: http.StatusConflict, title: "Not enough balance"}, true
	case errs.CodePhoneAlreadyExists:
		return problemSpec{status: http.StatusConflict, title: "Phone already exists"}, true
	case errs.CodePublicCodeCollision:
		return problemSpec{status: http.StatusConflict, title: "Conflict"}, true

	// 401
	case errs.CodeInvalidCredentials:
		return problemSpec{status: http.StatusUnauthorized, title: "Invalid credentials"}, true

	// 403
	case errs.CodeUserInactive:
		return problemSpec{status: http.StatusForbidden, title: "User inactive"}, true

	// 404
	case errs.CodeAccountNotFound:
		return problemSpec{status: http.StatusNotFound, title: "Not Found"}, true

	// 500
	case errs.CodeRolesNotFound, errs.CodeInternal:
		return problemSpec{status: http.StatusInternalServerError, title: "Internal Server Error"}, true
	}

	return problemSpec{}, false
}

func (h *Handler) WriteServiceError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	if c, ok := domainCode(err); ok {
		if ps, ok := problemForCode(c); ok {
			detail := err.Error()
			code := string(c)
			WriteProblem(w, ps.status, ps.title, &detail, &code, instanceFromRequest(r))
			return
		}
	}

	detail := err.Error()
	code := string(errs.CodeInternal)
	WriteProblem(w, http.StatusInternalServerError, "Internal Server Error", &detail, &code, instanceFromRequest(r))
}
