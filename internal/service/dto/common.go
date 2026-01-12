package dto

import "time"

type RoleCode string

const (
	RoleClient  RoleCode = "CLIENT"
	RoleCashier RoleCode = "CASHIER"
	RoleAdmin   RoleCode = "ADMIN"
)

type UserBase struct {
	ID        int64     `validate:"required,gt=0"`
	Phone     string    `validate:"required,phone"`
	IsActive  bool      `validate:"-"`
	CreatedAt time.Time `validate:"required"`
}

type UserWithRoles struct {
	UserBase
	Roles []RoleCode `validate:"required,min=1,dive,oneof=CLIENT CASHIER ADMIN"`
}

type AccountBase struct {
	ID              int64     `validate:"required,gt=0"`
	PublicCode      string    `validate:"required,min=6,max=64"`
	BalancePoints   int       `validate:"required,gte=0"`
	TotalSpendMoney string    `validate:"required"`
	LevelCode       string    `validate:"required,min=1,max=64"`
	CreatedAt       time.Time `validate:"required"`
}
