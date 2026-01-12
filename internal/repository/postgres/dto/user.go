package dto

type UserRow struct {
	ID           int64
	Phone        string
	PasswordHash string
	IsActive     bool
	CreatedAt    Ts
}

type UserWithRoles struct {
	UserRow
	Roles []RoleCode
}
