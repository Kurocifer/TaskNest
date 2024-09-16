package models

import "github.com/dgrijalva/jwt-go"

type UserAuthRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserAuthResponseBody struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

type User struct {
	ID       int
	Username string
	Password string
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
