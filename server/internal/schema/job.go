package schema

import (
	"time"

	"gorm.io/gorm"
)

type Job struct {
	gorm.Model

	TitleID      uint `json:"titleId"`
	DepartmentID uint `json:"departmentId"`
	LocationID   uint `json:"locationId"`

	Title      Title      `gorm:"foreignKey:TitleID"`      // Relationship to Title
	Department Department `gorm:"foreignKey:DepartmentID"` // Relationship to Department
	Location   Location   `gorm:"foreignKey:LocationID"`   // Relationship to Location

	Duties         string `json:"duties"`
	Qualifications string `json:"qualifications"`

	//* hierarchical relationship
	ManagerID    uint  `json:"managerId"`
	Subordinates []Job `gorm:"foreignKey:ManagerID"` // Jobs where this job is the manager

	AssignedJobs []AssignedJob `json:"assignedJobs"`
}

type AssignedJob struct {
	gorm.Model
	JobID  uint `json:"jobId"`  // Foreign key
	UserID uint `json:"userId"` // Foreign key

	Job Job `gorm:"foreignKey:JobID"` // Relationship to Job

	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}
