package schema

// type Position struct {
// 	BaseModel

// 	// Basic Position Details
// 	Title          string `json:"title"`
// 	Description    string `json:"description"`
// 	Duties         string `json:"duties"`
// 	Qualifications string `json:"qualifications"`
// 	Experience     string `json:"experience"`
// 	MinSalary      int    `json:"minSalary"`
// 	MaxSalary      int    `json:"maxSalary"`

// 	// Foreign Keys
// 	DepartmentID uint `json:"departmentId"`
// 	LocationID   uint `json:"locationId"`
// 	CompanyID    uint `json:"companyId"`

// 	// Associated Structs
// 	Department *Department `json:"department" gorm:"foreignKey:DepartmentID"`
// 	Location   *Location   `json:"location" gorm:"foreignKey:LocationID"`

// 	// Hierarchical Relationship
// 	ManagerID    *uint       `json:"managerId"`                                //* if not reference, will cause foreign key constraint error when null
// 	Subordinates []*Position `json:"subordinates" gorm:"foreignKey:ManagerID"` // Jobs where this job is the manager

// 	// Other Related Data
// 	UserPositions []*UserPosition `json:"assignedJobs"`
// }

// type UserPosition struct {
// 	BaseModel
// 	PositionID uint `json:"positionId"` // Foreign key
// 	UserID     uint `json:"userId"`     // Foreign key

// 	Position Position `json:"job" gorm:"foreignKey:JobID"`

// 	IsActive  bool       `json:"isActive"`
// 	StartDate time.Time  `json:"startDate"`
// 	EndDate   *time.Time `json:"endDate"`
// }
