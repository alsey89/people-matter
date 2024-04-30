package schema

import "time"

type Position struct {
	BaseModel
	CompanyID uint `json:"company_id" gorm:"onUpdate:CASCADE onDelete:CASCADE;not null"`
	// ------------------------------------------------------------------------------------------------
	Title          string `json:"title"`
	Description    string `json:"description"`
	Duties         string `json:"duties"`
	Qualifications string `json:"qualifications"`
	Experience     string `json:"experience"`
	MinSalary      int    `json:"minSalary"`
	MaxSalary      int    `json:"maxSalary"`
	// ------------------------------------------------------------------------------------------------
	DepartmentID uint        `json:"departmentId"`
	Department   *Department `json:"department" gorm:"foreignKey:DepartmentID"`
	LocationID   uint        `json:"locationId"`
	Location     *Location   `json:"location" gorm:"foreignKey:LocationID"`
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
	CompanyID uint `json:"company_id" gorm:"onUpdate:CASCADE onDelete:CASCADE;not null"`
	// ------------------------------------------------------------------------------------------------
	UserID    uint        `json:"userId"`
	Positions []*Position `json:"job" gorm:"foreignKey:PositionID"`
	// ------------------------------------------------------------------------------------------------
	EmploymentStatus EmploymentStatusEnum `json:"employmentStatus" gorm:"type:enum('full-time','part-time','contract','temporary','seasonal','intern');not null"`
	StartDate        *time.Time           `json:"startDate"`
	EndDate          *time.Time           `json:"endDate"`
}
