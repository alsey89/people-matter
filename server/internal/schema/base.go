package schema

import (
	"time"
)

type BaseModel struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	// DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type BaseModelWithTime struct {
	BaseModel
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
