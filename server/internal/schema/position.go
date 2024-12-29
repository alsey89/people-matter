package schema

type Position struct {
	BaseModel
	TenantIdentifier string `json:"-" gorm:"not null;index"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Qualifications   string `json:"qualifications"`
	Responsibilities string `json:"responsibilities"`
	SalaryRange      string `json:"salaryRange"`

	// Permissions
	Permissions []Permission `json:"permissions" gorm:"many2many:position_permissions;"`
}

type Permission struct {
	BaseModel
	TenantIdentifier string `json:"-" gorm:"not null;index"`
	Name             string `json:"name"`
	Description      string `json:"description"`
}

type UserPosition struct {
	BaseModel
	TenantIdentifier string `json:"-" gorm:"not null;index"`

	UserID     uint `json:"userId" gorm:"not null;index"`
	PositionID uint `json:"positionId" gorm:"not null;index"`
	LocationID uint `json:"locationId" gorm:"not null;index"`

	// Associations
	User     *User     `json:"user" gorm:"foreignKey:UserID"`
	Position *Position `json:"position" gorm:"foreignKey:PositionID"`
	Location *Location `json:"location" gorm:"foreignKey:LocationID"`
}
