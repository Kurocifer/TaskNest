package models

type UserAuthRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRegistrationResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type User struct {
	ID       int
	Username string
	Password string
}
