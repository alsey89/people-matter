package schema

import (
	"gorm.io/gorm"
)

type Document struct {
	gorm.Model
	UserID uint   `json:"userId" gorm:"not null"`
	URL    string `json:"url"    gorm:"type:text;not null"`
}
