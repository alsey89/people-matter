package schema

import (
	"time"

	"gorm.io/gorm"
)

// ======================
//  COMPANY & USER
// ======================

type Company struct {
	gorm.Model
	TenantID string `json:"-" gorm:"not null;index"`

	Name    string `json:"name"    gorm:"type:varchar(255);not null"`
	LogoURL string `json:"logoUrl" gorm:"type:text"`

	// Contact Information
	Email   string `json:"email"   gorm:"type:varchar(255);not null"`
	Phone   string `json:"phone"   gorm:"type:varchar(255);not null"`
	Website string `json:"website" gorm:"type:varchar(255);not null"`

	// Address Information (embedded addresses)
	ContactAddress Address `json:"contactAddress" gorm:"embedded;embeddedPrefix:contact_"`
	BillingAddress Address `json:"billingAddress" gorm:"embedded;embeddedPrefix:billing_"`

	// Account Quotas
	LocationQuota int `json:"branchQuota"  gorm:"default:1"`
	EmployeeQuota int `json:"employeeQuota" gorm:"default:10"`

	// Associations
	Users     []User     `json:"users"     gorm:"foreignKey:CompanyID"`
	Documents []Document `json:"documents" gorm:"polymorphic:Documentable;"`
}

type User struct {
	gorm.Model
	CompanyID    uint   `json:"companyId" gorm:"uniqueIndex:idx_company_email;not null"`
	Name         string `json:"name"`
	Email        string `json:"email" gorm:"uniqueIndex:idx_company_email;not null"`
	PasswordHash string `json:"-"`

	EmailVerified bool `json:"emailVerified" gorm:"default:false"`

	// Could be null if a user doesn’t have an active compensation
	ActiveCompensationID *uint `json:"-" gorm:"index;default:null"`

	Documents []Document `json:"documents" gorm:"polymorphic:Documentable;"`
}

// ======================
//  LOCATION
// ======================

type Location struct {
	gorm.Model
	CompanyID uint `json:"companyId" gorm:"not null;index"`

	Name string `json:"name" gorm:"type:varchar(255);not null"`

	// Use an embedded Address for location details.
	Address Address `json:"address" gorm:"embedded;embeddedPrefix:address_"`

	Email   *string `json:"email"   gorm:"type:varchar(255)"`
	Phone   *string `json:"phone"   gorm:"type:varchar(255)"`
	Website *string `json:"website" gorm:"type:varchar(255)"`

	// For referencing employees at this location.
	UserPositions []UserPosition `json:"userPositions" gorm:"foreignKey:LocationID"`

	Documents []Document `json:"documents" gorm:"polymorphic:Documentable;"`
}

// ======================
//  POSITIONS & PERMISSIONS
// ======================

type Position struct {
	gorm.Model
	CompanyID        uint   `json:"companyId" gorm:"not null;index"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Qualifications   string `json:"qualifications"`
	Responsibilities string `json:"responsibilities"`

	SalaryMin      *float64 `json:"salaryMin"`
	SalaryMax      *float64 `json:"salaryMax"`
	SalaryCurrency *string  `json:"salaryCurrency"`

	Permissions []Permission `gorm:"many2many:position_permissions;"`

	Documents []Document `json:"documents" gorm:"polymorphic:Documentable;"`
}

type Permission struct {
	gorm.Model
	CompanyID   uint   `json:"companyId" gorm:"not null;index"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Documents []Document `json:"documents" gorm:"polymorphic:Documentable;"`
}

type PositionPermission struct {
	gorm.Model
	CompanyID    uint `json:"companyId" gorm:"not null;index"`
	PositionID   uint `json:"positionId"   gorm:"not null;index"`
	PermissionID uint `json:"permissionId" gorm:"not null;index"`

	Documents []Document `json:"documents" gorm:"polymorphic:Documentable;"`
}

type UserPosition struct {
	gorm.Model
	CompanyID  uint `json:"companyId"  gorm:"not null;index"`
	UserID     uint `json:"userId"     gorm:"not null;index"`
	PositionID uint `json:"positionId" gorm:"not null;index"`
	LocationID uint `json:"locationId" gorm:"not null;index"`

	StartedAt time.Time  `json:"startedAt" gorm:"not null"`
	EndedAt   *time.Time `json:"endedAt"   gorm:"default:null"`

	// Associations
	User     User     `gorm:"foreignKey:UserID;onDelete:CASCADE"`
	Position Position `gorm:"foreignKey:PositionID;onDelete:CASCADE"`
	Location Location `gorm:"foreignKey:LocationID;onDelete:CASCADE"`

	Documents []Document `json:"documents" gorm:"polymorphic:Documentable;"`
}

// ======================
//  COMPENSATION & PAYROLL
// ======================

type Compensation struct {
	gorm.Model
	CompanyID uint `json:"companyId" gorm:"not null;index"`
	UserID    uint `json:"userId"    gorm:"not null;index"`

	Amount   float64 `json:"amount"   gorm:"not null"`
	Currency string  `json:"currency" gorm:"not null"`
	Interval string  `json:"interval" gorm:"not null"` // e.g., "monthly", "bi-weekly"

	ChannelType string  `json:"channelType"    gorm:"not null"`
	ChannelCode *string `json:"channelId"` // e.g. bank code
	Account     string  `json:"channelAccount" gorm:"not null"`

	StartedAt time.Time  `json:"startedAt" gorm:"not null"`
	EndedAt   *time.Time `json:"endedAt"   gorm:"default:null"`

	// Associations
	User     User      `gorm:"foreignKey:UserID"`
	Payments []Payment `gorm:"foreignKey:CompensationID"`

	Documents []Document `json:"documents" gorm:"polymorphic:Documentable;"`
}

type Payment struct {
	gorm.Model
	CompanyID      uint      `json:"companyId"      gorm:"not null;index"`
	UserID         uint      `json:"userId"         gorm:"not null;index"`
	CompensationID uint      `json:"compensationId" gorm:"not null;index"`
	Currency       string    `json:"currency"       gorm:"not null"`
	PaidAt         time.Time `json:"paidAt"         gorm:"not null"`

	// Associations
	Compensation Compensation `gorm:"foreignKey:CompensationID"` //base compensation
	Bonuses      []Bonus      `gorm:"foreignKey:PaymentID"`      //extra bonus
	Expenses     []Expense    `gorm:"foreignKey:PaymentID"`      //expense claims
	Adjustments  []Adjustment `gorm:"foreignKey:PaymentID"`      //admin adjustments + or -

	Documents []Document `json:"documents" gorm:"polymorphic:Documentable;"`
}

type Bonus struct {
	gorm.Model
	CompanyID   uint    `json:"companyId" gorm:"not null;index"`
	PaymentID   uint    `json:"paymentId" gorm:"not null;index"`
	Amount      float64 `json:"amount"    gorm:"not null"`
	Description string  `json:"description"`

	Documents []Document `json:"documents" gorm:"polymorphic:Documentable;"`
}

type Adjustment struct {
	gorm.Model
	CompanyID   uint    `json:"companyId" gorm:"not null;index"`
	PaymentID   uint    `json:"paymentId" gorm:"not null;index"`
	Amount      float64 `json:"amount"    gorm:"not null"`
	Description string  `json:"description"`

	Documents []Document `json:"documents" gorm:"polymorphic:Documentable;"`
}

type Expense struct {
	gorm.Model
	CompanyID   uint    `json:"companyId" gorm:"not null;index"`
	PaymentID   uint    `json:"paymentId" gorm:"not null;index"`
	Amount      float64 `json:"amount"    gorm:"not null"`
	Description string  `json:"description"`

	Documents []Document `json:"documents" gorm:"polymorphic:Documentable;"`
}

// ======================
//  DOCUMENTS (Polymorphic Association)
// ======================

// Document represents a file or resource that can be attached to various entities.
// It uses a polymorphic association to allow it to be attached to multiple entities.
type Document struct {
	gorm.Model
	Name        string  `json:"name" gorm:"not null"`
	URL         string  `json:"url"  gorm:"not null"`
	Description *string `json:"description"`

	// Polymorphic association fields.
	DocumentableID   uint   `json:"documentableId"`
	DocumentableType string `json:"documentableType" gorm:"type:varchar(255)"`
}

// ======================
//  TIMEKEEPING (Clock In/Out)
// ======================

type TimeLog struct {
	gorm.Model
	CompanyID  uint       `json:"companyId" gorm:"not null;index"`
	UserID     uint       `json:"userId"    gorm:"not null;index"`
	LocationID *uint      `json:"locationId" gorm:"index"` // Optional if you track location
	PositionID *uint      `json:"positionId" gorm:"index"` // Added to associate with a position, if needed
	ClockIn    time.Time  `json:"clockIn"   gorm:"not null"`
	ClockOut   *time.Time `json:"clockOut"` // Nil if not yet clocked out
	Notes      *string    `json:"notes"`

	// Associations
	User     User      `gorm:"foreignKey:UserID"`
	Position *Position `gorm:"foreignKey:PositionID"`
	Location *Location `gorm:"foreignKey:LocationID"`

	Documents []Document `json:"documents" gorm:"polymorphic:Documentable;"`
}

// ======================
// EMBEDDED NON-TABLE STRUCTS
// ======================

type Address struct {
	Street     string `json:"street"     gorm:"type:text;not null"`
	City       string `json:"city"       gorm:"type:varchar(255);not null"`
	Country    string `json:"country"    gorm:"type:varchar(255);not null"`
	PostalCode string `json:"postalCode" gorm:"type:varchar(255);not null"`
}
