package schema

import (
	"time"

	"gorm.io/gorm"
)

// USER & IDENTITY ---------------------

type User struct {
	gorm.Model
	CompanyID uint `json:"companyId" gorm:"uniqueIndex:idx_company_email;not null"`

	Name           string `json:"name"`
	Email          string `json:"email" gorm:"uniqueIndex:idx_company_email;not null"`
	PasswordHash   string `json:"-"`
	CompensationID uint   `json:"-" gorm:"not null;index"`
}

// COMPANY & LOCATION ---------------------

type Company struct {
	gorm.Model
	TenantIdentifier string `json:"-" gorm:"not null;index"`

	// Company Information ---------------------
	Name    string `json:"name" gorm:"type:varchar(255);not null"`
	LogoURL string `json:"logoUrl" gorm:"text"`
	// Contact Information ---------------------
	Email   string `json:"email" gorm:"type:varchar(255);not null"`
	Phone   string `json:"phone" gorm:"type:varchar(255);not null"`
	Website string `json:"website" gorm:"type:varchar(255);not null"`
	//Address Information ---------------------
	ContactAddress    string `json:"address" gorm:"type:text; not null"`
	ContactCity       string `json:"city" gorm:"type:varchar(255);not null"`
	ContactCountry    string `json:"country" gorm:"type:varchar(255);not null"`
	ContactPostalCode string `json:"postalCode" gorm:"type:varchar(255);not null"`

	BillingAddress    string `json:"billingAddress" gorm:"type:text; not null"`
	BillingCity       string `json:"billingCity" gorm:"type:varchar(255);not null"`
	BillingCountry    string `json:"billingCountry" gorm:"type:varchar(255);not null"`
	BillingPostalCode string `json:"billingPostalCode" gorm:"type:varchar(255);not null"`
	// Account Information ---------------------
	LocationQuota     int `json:"branchQuota" gorm:"default:1"`
	LocationQuotaUsed int `json:"branchQuotaUsed" gorm:"default:0"`
	EmployeeQuota     int `json:"employeeQuota" gorm:"default:10"`
	EmployeeQuotaUsed int `json:"employeeQuotaUsed" gorm:"default:0"`

	// Associations
	Users []User `json:"users" gorm:"foreignKey:CompanyID"`
}

type Location struct {
	gorm.Model
	CompanyID uint `json:"companyId" gorm:"not null;index"`

	Name       string  `json:"name" gorm:"type:varchar(255);not null"`
	Address    string  `json:"address" gorm:"type:text;not null"`
	City       string  `json:"city" gorm:"type:varchar(255);not null"`
	Country    string  `json:"country" gorm:"type:varchar(255);not null"`
	PostalCode string  `json:"postalCode" gorm:"type:varchar(255);not null"`
	Email      *string `json:"email" gorm:"type:varchar(255);not null"`
	Phone      *string `json:"phone" gorm:"type:varchar(255);not null"`
	Website    *string `json:"website" gorm:"type:varchar(255);not null"`

	// Associations
	UserPositions []UserPosition `json:"userPositions" gorm:"foreignKey:LocationID"`
}

// POSITIONS & PERMISSIONS ---------------------

type Position struct {
	gorm.Model
	CompanyID uint `json:"companyId" gorm:"not null;index"`

	Name             string `json:"name"`
	Description      string `json:"description"`
	Qualifications   string `json:"qualifications"`
	Responsibilities string `json:"responsibilities"`
	SalaryRange      string `json:"salaryRange"`

	// Associations
	Permissions []Permission `json:"permissions" gorm:"many2many:position_permissions;"`
}

type UserPosition struct {
	gorm.Model
	CompanyID uint `json:"companyId" gorm:"not null;index"`

	UserID     uint `json:"userId" gorm:"not null;index"`
	PositionID uint `json:"positionId" gorm:"not null;index"`
	LocationID uint `json:"locationId" gorm:"not null;index"`

	StartedAt time.Time  `json:"startedAt" gorm:"not null"`
	EndedAt   *time.Time `json:"endedAt" gorm:"default:null"`

	// Associations
	User     *User     `json:"user" gorm:"foreignKey:UserID;onDelete:CASCADE"`
	Position *Position `json:"position" gorm:"foreignKey:PositionID;onDelete:CASCADE"`
	Location *Location `json:"location" gorm:"foreignKey:LocationID;onDelete:CASCADE"`
}

type PositionPermission struct {
	gorm.Model
	CompanyID uint `json:"companyId" gorm:"not null;index"`

	PositionID   uint `json:"positionId" gorm:"not null;index"`
	PermissionID uint `json:"permissionId" gorm:"not null;index"`

	// Associations
	Position   *Position   `json:"position" gorm:"foreignKey:PositionID;onDelete:CASCADE"`
	Permission *Permission `json:"permission" gorm:"foreignKey:PermissionID;onDelete:CASCADE"`
}

type Permission struct {
	gorm.Model
	CompanyID uint `json:"companyId" gorm:"not null;index"`

	Name        string `json:"name"`
	Description string `json:"description"`
}

// COMPENSATION ---------------------

type Compensation struct {
	gorm.Model
	CompanyID uint `json:"companyId" gorm:"not null;index"`

	UserID   uint    `json:"userId" gorm:"not null;index"`
	Amount   float64 `json:"amount" gorm:"not null"`
	Currency string  `json:"currency" gorm:"not null"`
	Interval string  `json:"interval" gorm:"not null"`

	ChannelType string  `json:"channelType" gorm:"not null"`    // e.g. "bank", "paypal", "cash"
	ChannelID   *string `json:"channelId"`                      // e.g. "807"[Sinopac Bank]
	Account     string  `json:"channelAccount" gorm:"not null"` // e.g. "1234567890", "johndoe@paypal", "cash"

	StartedAt time.Time  `json:"startedAt" gorm:"not null"`
	EndedAt   *time.Time `json:"endedAt" gorm:"default:null"`

	// Associations
	User     *User     `json:"user" gorm:"foreignKey:UserID"`
	Payments []Payment `json:"payments" gorm:"foreignKey:CompensationID"`
}

type Payment struct {
	gorm.Model
	CompanyID uint `json:"companyId" gorm:"not null;index"`

	CompensationID uint      `json:"compensationId" gorm:"not null;index"`
	Amount         float64   `json:"amount" gorm:"not null"`
	Currency       string    `json:"currency" gorm:"not null"`
	Bonus          float64   `json:"bonus" gorm:"default:0"`
	Adjustments    float64   `json:"adjustments" gorm:"default:0"`
	PaidAt         time.Time `json:"paidAt" gorm:"not null"`
}
