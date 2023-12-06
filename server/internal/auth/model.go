package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Claims struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	IsAdmin bool               `json:"isAdmin" default:"false"`
	Email   string             `json:"email"`
	jwt.RegisteredClaims
}

type Credentials struct {
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}
