package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

// * JWT claims struct
type Claims struct {
	ID        uint   `json:"id" bson:"id"`
	CompanyID uint   `json:"companyId"`
	Role      string `json:"isAdmin" default:"user"`
	Email     string `json:"email"`
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

//------------------------------------------------------------
