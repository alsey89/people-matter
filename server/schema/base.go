package schema

import (
	"time"
)

type BaseModel struct {
	ID        uint `json:"Id" gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}
