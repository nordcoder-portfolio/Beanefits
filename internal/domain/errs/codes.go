package errs

type Code string

const (
	CodeInvalidPhone          Code = "INVALID_PHONE"
	CodeInvalidPublicCode     Code = "INVALID_PUBLIC_CODE"
	CodeInvalidPoints         Code = "INVALID_POINTS"
	CodeNotEnoughBalance      Code = "NOT_ENOUGH_BALANCE"
	CodeInvalidPurchaseAmount Code = "INVALID_PURCHASE_AMOUNT"

	CodeInvalidUserID Code = "INVALID_USER_ID"

	CodeInvalidRuleset Code = "INVALID_RULESET"
	CodeRulesetNotFound Code = "RULESET_NOT_FOUND"
	CodeInvalidLevels  Code = "INVALID_LEVELS"
	CodeInvalidMoney   Code = "INVALID_MONEY"

	CodePhoneAlreadyExists Code = "PHONE_ALREADY_EXISTS"
	CodeInvalidCredentials Code = "INVALID_CREDENTIALS"
	CodeUserInactive       Code = "USER_INACTIVE"
	CodeAccountNotFound    Code = "ACCOUNT_NOT_FOUND"
	CodeUserNotFound       Code = "USER_NOT_FOUND"
	CodeRolesNotFound      Code = "ROLES_NOT_FOUND"

	CodePublicCodeCollision Code = "PUBLIC_CODE_COLLISION"
	CodeInternal            Code = "INTERNAL_ERROR"
)
