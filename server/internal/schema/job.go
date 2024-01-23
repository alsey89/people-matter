package schema

import (
	"time"

	"gorm.io/gorm"
)

type Job struct {
	gorm.Model

	// Basic Job Details
	Title          string `json:"title"`
	Description    string `json:"description"`
	Duties         string `json:"duties"`
	Qualifications string `json:"qualifications"`
	Experience     string `json:"experience"`
	MinSalary      int    `json:"minSalary"`
	MaxSalary      int    `json:"maxSalary"`

	// Foreign Keys
	DepartmentID uint `json:"departmentId"`
	LocationID   uint `json:"locationId"`
	CompanyID    uint `json:"companyId"`

	// Associated Structs
	Department *Department `json:"department" gorm:"foreignKey:DepartmentID"`
	Location   *Location   `json:"location" gorm:"foreignKey:LocationID"`

	// Hierarchical Relationship
	ManagerID    uint   `json:"managerId"`
	Subordinates []*Job `json:"subordinates" gorm:"foreignKey:ManagerID"` // Jobs where this job is the manager

	// Other Related Data
	AssignedJobs []*AssignedJob `json:"assignedJobs"`
}

type AssignedJob struct {
	gorm.Model
	JobID  uint `json:"jobId"`  // Foreign key
	UserID uint `json:"userId"` // Foreign key

	Job Job `json:"job" gorm:"foreignKey:JobID"` // Relationship to Job

	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}
