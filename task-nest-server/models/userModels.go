package models

type UserRegistrationRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRegistrationResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}
