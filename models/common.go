package models

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type Error struct {
	Message string `json:"message"`
	Code    int16  `json:"code"`
}

type Timestamp time.Time

type DateFields struct {
	CreatedAt Timestamp
	UpdatedAt Timestamp
}

type JWT struct {
	Authorization string `json:"Authorization" binding:"required"`
}

type User struct {
	ID         uint64 `json:"id"`
	Email      string `json:"email"`
	Password   string `json:"-" binding:"required"`
	Firstname  string `json:"firstname" binding:"required"`
	Lastname   string `json:"lastname"`
	Phone      string
	Agree      bool `json:"-"`
	AllowPromo bool `json:"-"`
}

type UserClaims struct {
	Refresh bool
	User    User
	*jwt.StandardClaims
}
