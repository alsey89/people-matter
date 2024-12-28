package schema

type User struct {
	BaseModel
	FSPID uint   `json:"fspId" gorm:"not null;uniqueIndex:idx_fsp_email"`
	Email string `json:"email" gorm:"not null;uniqueIndex:idx_fsp_email"`
	// User Information ---------------------
	FirstName  string  `json:"firstName"`
	MiddleName *string `json:"middleName"`
	LastName   string  `json:"lastName"`
	AvatarURL  string  `json:"avatarUrl"`
	// Account Information ---------------------
	PasswordHash      string `json:"-"`
	EmailConfirmed    bool   `json:"emailConfirmed"`
	LockOutEnabled    bool   `json:"lockOutEnabled"`
	LockOutTime       int    `json:"lockOutTime"`
	AccessFailedCount int    `json:"accessFailedCount"`
	// Associations ---------------------
	UserFSPRole       *UserFSPRole        `json:"userFspRole"`
	UserMemorialRoles *[]UserMemorialRole `json:"userMemorialRoles"`
}
