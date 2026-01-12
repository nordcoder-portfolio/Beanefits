package user

type RoleCode string

const (
	RoleClient  RoleCode = "CLIENT"
	RoleCashier RoleCode = "CASHIER"
	RoleAdmin   RoleCode = "ADMIN"
)
