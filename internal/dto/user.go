package dto

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegistrationUser struct {
	Username string `json:"username" validate:"max=20,min=3"`
	Password string `json:"password" validate:"max=16,min=8"`
}
