package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// * User model ---------------------------------------------------------------
type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"-" bson:"password"` //! this omits password from json, MongoDB is bson
	IsAdmin   bool               `json:"isAdmin" bson:"isAdmin" default:"false"`
	AvatarURL string             `json:"avatarUrl" bson:"avatarUrl"`
	CreatedAt string             `json:"createdAt" bson:"createdAt"`
	UpdatedAt string             `json:"updatedAt" bson:"updatedAt"`
}
