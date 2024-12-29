package schema

import "time"

type User struct {
	BaseModel
	TenantID uint `json:"tenantId" gorm:"not null;uniqueIndex:idx_tenant_email"`
	// User Information ---------------------
	FirstName  string  `json:"firstName"`
	MiddleName *string `json:"middleName,omitempty"`
	LastName   string  `json:"lastName"`
	AvatarURL  string  `json:"avatarUrl,omitempty"`
	// Account Information ---------------------
	Email             string     `json:"email" gorm:"not null;uniqueIndex:idx_tenant_email"`
	PasswordHash      string     `json:"-"`
	EmailConfirmed    bool       `json:"emailConfirmed"`
	LockOutEnabled    bool       `json:"lockOutEnabled"`
	LockOutExpiry     *time.Time `json:"lockOutExpiry,omitempty"`
	AccessFailedCount int        `json:"accessFailedCount"`
	// Comp Information ---------------------
	Comps []Comp `json:"comps" gorm:"foreignKey:UserID"`
}
