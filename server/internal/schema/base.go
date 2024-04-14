package schema

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint `json:"Id" gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

//USER MODEL ---------------------------------------------------------------

type User struct {
	//* account information
	BaseModel      // This includes fields ID, CreatedAt, UpdatedAt, DeletedAt
	CompanyID uint `json:"companyId" gorm:"onUpdate:CASCADE;onDelete:CASCADE"`

	Email     string     `json:"email"      gorm:"uniqueIndex"`
	Password  string     `json:"-"          gorm:"type:varchar(100)"` //* Password is not returned in JSON
	AvatarURL string     `json:"avatarUrl"  gorm:"type:text"`
	Role      string     `json:"role"       gorm:"enum:admin,manager,user;default:user"`
	IsActive  bool       `json:"isActive"   gorm:"default:true"`
	LastLogin *time.Time `json:"lastLogin"  gorm:"default:null"`

	//* personal information
	FirstName        string            `json:"firstName"`
	MiddleName       *string           `json:"middleName"`
	LastName         string            `json:"lastName"`
	Nickname         string            `json:"nickname"`
	ContactInfo      *ContactInfo      `json:"contactInfo"       gorm:"foreignKey:UserID"`
	EmergencyContact *EmergencyContact `json:"emergencyContact"  gorm:"foreignKey:UserID"`

	//*job information
	UserPosition *UserPosition `json:"assignedJob" gorm:"foreignKey:UserID"`

	//* salary information
	SalaryID *uint     `json:"salaryId"  gorm:"foreignKey:UserID"`
	Payments []Payment `json:"Payments" gorm:"foreignKey:UserID"`

	//* leave & attendance
	Leave      []Leave      `json:"leave"`
	Attendance []Attendance `json:"attendance"`
}

// enum constants for User.Role
const (
	RoleAdmin   = "admin"
	RoleManager = "manager"
	RoleUser    = "user"
)

type ContactInfo struct {
	BaseModel
	UserID     uint   `json:"userId"`
	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
	Mobile     string `json:"mobile"`
	Email      string `json:"email"`
}

type EmergencyContact struct {
	BaseModel
	UserID     uint    `json:"userId"`
	FirstName  string  `json:"firstName"`
	MiddleName *string `json:"middleName"`
	LastName   string  `json:"lastName"`
	Relation   string  `json:"relation"`
	Mobile     string  `json:"mobile"`
	Email      string  `json:"email"`
}

// COMPANY MODEL ---------------------------------------------------------------

type Company struct {
	BaseModel
	Name       string  `json:"name"`
	LogoURL    *string `json:"logoUrl"`
	Website    *string `json:"website"`
	Email      *string `json:"email"`
	Phone      *string `json:"phone"`
	Address    *string `json:"address"`
	City       *string `json:"city"`
	State      *string `json:"state"`
	Country    *string `json:"country"`
	PostalCode *string `json:"postalCode"`

	// Relationships
	Departments []Department `json:"departments" gorm:"foreignKey:CompanyID"`
	Titles      []Title      `json:"titles" gorm:"foreignKey:CompanyID"`
	Locations   []Location   `json:"locations" gorm:"foreignKey:CompanyID"`
	Positions   []Position   `json:"jobs" gorm:"foreignKey:CompanyID"`
}

type Title struct {
	BaseModel
	CompanyID uint `json:"companyId"`

	Name        string `json:"name"`
	Description string `json:"description"`
}

type Department struct {
	BaseModel
	CompanyID uint `json:"companyId" gorm:"onDelete:CASCADE onUpdate:CASCADE"`

	Name        string `json:"name"`
	Description string `json:"description"`
}

type Location struct {
	BaseModel
	CompanyID uint `json:"companyId" gorm:"onDelete:CASCADE onUpdate:CASCADE"`

	Name         string `json:"name"`
	IsHeadOffice bool   `json:"isHeadOffice"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
	PostalCode   string `json:"postalCode"`
}

// POSITION MODEL ---------------------------------------------------------------

type Position struct {
	BaseModel
	// Foreign Keys
	DepartmentID uint `json:"departmentId"`
	LocationID   uint `json:"locationId"`
	CompanyID    uint `json:"companyId" gorm:"onUpdate:CASCADE;onDelete:CASCADE"`

	// Basic Position Details
	Title          string `json:"title"`
	Description    string `json:"description"`
	Duties         string `json:"duties"`
	Qualifications string `json:"qualifications"`
	Experience     string `json:"experience"`
	MinSalary      int    `json:"minSalary"`
	MaxSalary      int    `json:"maxSalary"`

	// Associated Structs
	Department *Department `json:"department" gorm:"foreignKey:DepartmentID"`
	Location   *Location   `json:"location" gorm:"foreignKey:LocationID"`

	// Hierarchical Relationship
	ManagerID    *uint       `json:"managerId"`                                //* if not reference, will cause foreign key constraint error when null
	Subordinates []*Position `json:"subordinates" gorm:"foreignKey:ManagerID"` // Jobs where this job is the manager

	// Other Related Data
	UserPositions []*UserPosition `json:"assignedJobs"`
}

type UserPosition struct {
	BaseModel
	CompanyID uint `json:"companyId" gorm:"onUpdate:CASCADE;onDelete:CASCADE"`

	PositionID uint `json:"positionId"` // Foreign key
	UserID     uint `json:"userId"`     // Foreign key

	Position Position `json:"job" gorm:"foreignKey:JobID"`

	IsActive  bool       `json:"isActive"`
	StartDate time.Time  `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`
}

// ATTENDANCE MODEL ---------------------------------------------------------------

type Attendance struct {
	BaseModel
	CompanyID uint `json:"companyId" gorm:"onUpdate:CASCADE;onDelete:CASCADE"`

	UserID   uint      `json:"user_id"   gorm:"not null"`
	Date     time.Time `json:"date"      gorm:"type:date;not null"`
	ClockIn  time.Time `json:"clock_in"  gorm:"type:time"`
	ClockOut time.Time `json:"clock_out" gorm:"type:time"`
	Notes    string    `json:"notes"     gorm:"type:text"`
}

// LEAVE MODEL ---------------------------------------------------------------

type Leave struct {
	BaseModel
	CompanyID uint `json:"companyId" gorm:"onUpdate:CASCADE;onDelete:CASCADE"`

	UserID    uint        `json:"user_id" gorm:"not null"`
	Type      LeaveType   `json:"type" gorm:"type:varchar(100);not null"`
	Status    LeaveStatus `json:"status" gorm:"type:varchar(100);not null"`
	StartDate time.Time   `json:"startDate" gorm:"type:date;not null"`
	EndDate   time.Time   `json:"endDate" gorm:"type:date;not null"`
	Reason    string      `json:"reason" gorm:"type:text"`
	Documents []Document  `json:"documents" gorm:"many2many:leave_documents;"`
}

type LeaveStatus string

const (
	StatusPending  LeaveStatus = "pending"
	StatusApproved LeaveStatus = "approved"
	StatusRejected LeaveStatus = "rejected"
)

type LeaveType string

const (
	TypeSick        LeaveType = "sick"
	TypeVacation    LeaveType = "vacation"
	TypeMaternity   LeaveType = "maternity"
	TypePaternity   LeaveType = "paternity"
	TypeBereavement LeaveType = "bereavement"
	TypeOther       LeaveType = "other"
)

// add checks before creating entry
func (l *Leave) BeforeCreate(tx *gorm.DB) (err error) {
	if !IsValidLeaveStatus(l.Status) {
		return gorm.ErrInvalidData
	}
	if !IsValidLeaveType(l.Type) {
		return gorm.ErrInvalidData
	}
	return nil
}

func IsValidLeaveStatus(status LeaveStatus) bool {
	switch status {
	case StatusPending, StatusApproved, StatusRejected:
		return true
	default:
		return false
	}
}

func IsValidLeaveType(leaveType LeaveType) bool {
	switch leaveType {
	case TypeSick, TypeVacation, TypeMaternity, TypePaternity, TypeBereavement, TypeOther:
		return true
	default:
		return false
	}
}

// SALARY MODEL ---------------------------------------------------------------

type Salary struct {
	BaseModel
	CompanyID uint `json:"companyId" gorm:"onUpdate:CASCADE;onDelete:CASCADE"`

	UserID          uint      `json:"userId" gorm:"onUpdate:CASCADE;onDelete:CASCADE"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	PaymentInterval string    `json:"paymentInterval"`
	EffectiveDate   time.Time `json:"effectiveDate"`
	IsActive        bool      `json:"isActive" gorm:"default:false"`
	IsApproved      bool      `json:"isApproved" gorm:"default:false"`
}

type Payment struct {
	BaseModel
	CompanyID uint `json:"companyId" gorm:"onUpdate:CASCADE;onDelete:CASCADE"`

	UserID        uint      `json:"userId" gorm:"onUpdate:CASCADE;onDelete:CASCADE"`
	SalaryID      uint      `json:"salaryId" gorm:"onUpdate:CASCADE;onDelete:SET NULL"`
	PaymentDate   time.Time `json:"paymentDate"`
	Amount        float64   `json:"amount"`
	PaymentMethod string    `json:"paymentMethod"`
	Status        string    `json:"status"`
	PeriodStart   time.Time `json:"periodStart"`
	PeriodEnd     time.Time `json:"periodEnd"`
	Deductions    float64   `json:"deductions"`
	Bonuses       float64   `json:"bonuses"`
	Notes         string    `json:"notes"`
}

// DOCUMENT MODEL ---------------------------------------------------------------

type Document struct {
	BaseModel
	CompanyID uint `json:"companyId" gorm:"onUpdate:CASCADE;onDelete:CASCADE"`

	UserID uint   `json:"userId" gorm:"not null"`
	URL    string `json:"url"    gorm:"type:text;not null"`
}
