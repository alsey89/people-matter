package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ID      uint   `json:"id,omitempty" bson:"id,omitempty"`
	IsAdmin bool   `json:"isAdmin" default:"false"`
	Email   string `json:"email"`
	jwt.RegisteredClaims
}

type SignupCredentials struct {
	Email           string `json:"email" bson:"email"`
	Password        string `json:"password" bson:"password"`
	ConfirmPassword string `json:"confirmPassword" bson:"confirmPassword"`
}

type SigninCredentials struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}
