package dto

// RegisterIn — вход в usecase Auth.RegisterClient.
type RegisterIn struct {
	Phone    string `validate:"required,phone"`
	Password string `validate:"required,min=6,max=128"`
}

// LoginIn — вход в usecase Auth.Login.
type LoginIn struct {
	Phone    string `validate:"required,phone"`
	Password string `validate:"required,min=6,max=128"`
}

// AuthOut — выход для register/login.
type AuthOut struct {
	AccessToken string        `validate:"required"`
	User        UserWithRoles `validate:"required"`
	Account     AccountBase   `validate:"required"`
}
