package dto

type RegisterIn struct {
	Phone    string `validate:"required,phone"`
	Password string `validate:"required,min=6,max=128"`
}

type LoginIn struct {
	Phone    string `validate:"required,phone"`
	Password string `validate:"required,min=6,max=128"`
}

type AuthOut struct {
	AccessToken string        `validate:"required"`
	User        UserWithRoles `validate:"required"`
	Account     AccountBase   `validate:"required"`
}
