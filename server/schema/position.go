package schema

import (
	"time"
)

type Position struct {
	BaseModel
	CompanyID uint `json:"company_id" gorm:"constraint:OnDelete:CASCADE;not null"`
	// ------------------------------------------------------------------------------------------------
	Name           string `json:"name"`
	Description    string `json:"description"`
	Duties         string `json:"duties"`
	Qualifications string `json:"qualifications"`
	Experience     string `json:"experience"`
	MinSalary      int    `json:"minSalary"`
	MaxSalary      int    `json:"maxSalary"`
	// ------------------------------------------------------------------------------------------------
	DepartmentID *uint       `json:"departmentId" gorm:"constraint:OnDelete:SET NULL;foreignKey:DepartmentID"`
	Department   *Department `json:"department" gorm:"constraint:OnDelete:SET NULL;foreignKey:DepartmentID"`
	// ------------------------------------------------------------------------------------------------
	ManagerID    *uint       `json:"managerId"` //* should be nullable
	Subordinates []*Position `json:"subordinates" gorm:"foreignKey:ManagerID"`
	// ------------------------------------------------------------------------------------------------
	UserPositions []*UserPosition `json:"assignedJobs"`
}

type EmploymentStatusEnum string

const (
	FullTime  EmploymentStatusEnum = "full-time"
	PartTime  EmploymentStatusEnum = "part-time"
	Contract  EmploymentStatusEnum = "contract"
	Temporary EmploymentStatusEnum = "temporary"
	Seasonal  EmploymentStatusEnum = "seasonal"
	Intern    EmploymentStatusEnum = "intern"
)

type UserPosition struct {
	BaseModel
	CompanyID uint `json:"company_id" gorm:"constraint:OnDelete:CASCADE;not null"`
	// ------------------------------------------------------------------------------------------------
	UserID uint `json:"userId"`

	PositionID uint     `json:"positionId"`
	Position   Position `json:"position" gorm:"foreignKey:PositionID"`

	LocationID uint     `json:"locationId"`
	Location   Location `json:"location" gorm:"foreignKey:LocationID"`
	// ------------------------------------------------------------------------------------------------
	EmploymentStatus EmploymentStatusEnum `json:"employmentStatus" sql:"type:enum('full-time','part-time','contract','temporary','seasonal','intern');not null"`
	StartDate        time.Time            `json:"startDate"`
	EndDate          time.Time            `json:"endDate"`
}
