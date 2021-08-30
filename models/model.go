package models

import "github.com/golang-jwt/jwt"

type Error struct {
	Message string `json:"message"`
	Code int16 `json:"code"`
}

type JWT struct {
	Authorization string `json:"Authorization" binding:"required"`
}

type User struct {
	ID int64 `json:"id"`
	Email string `json:"email"`
}

type UserClaims struct {
	User
	*jwt.StandardClaims
}
