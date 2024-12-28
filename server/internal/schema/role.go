package schema

import "time"

type FSPRoleConst string

const (
	RoleFSPSuperAdmin FSPRoleConst = "superadmin"
	RoleFSPAdmin      FSPRoleConst = "admin"
	RoleFSPUser       FSPRoleConst = "user"
)

type MemorialRoleConst string

const (
	RoleMemSelf                 MemorialRoleConst = "self"
	RoleMemCurator              MemorialRoleConst = "curator"
	RoleMemContributor          MemorialRoleConst = "contributor"
	RoleMemContributorApplicant MemorialRoleConst = "contributor_applicant"
	RoleMemInsitu               MemorialRoleConst = "insitu"
)

type FSPRole struct {
	BaseModel
	Name        FSPRoleConst `json:"name" gorm:"default:'user';not null;unique;index" sql:"type:ENUM('superadmin', 'admin', 'user')"`
	Description string       `json:"description"`

	UserFSPRoles []UserFSPRole `json:"userRoles" gorm:"foreignKey:FSPRoleID"`
}

type MemorialRole struct {
	BaseModel
	Name        MemorialRoleConst `json:"name" gorm:"default:'contributor';not null;unique;index" sql:"type:ENUM('self', 'curator', 'contributor', 'insitu')"`
	Description string            `json:"description"`

	UserMemorialRoles []UserMemorialRole `json:"userRoles" gorm:"foreignKey:MemorialRoleID"`
}

type UserFSPRole struct {
	BaseModel
	FSPID     uint     `json:"fspID" gorm:"not null;uniqueIndex:idx_user_fsp"`
	FSP       *FSP     `json:"fsp" gorm:"foreignKey:FSPID"`
	UserID    uint     `json:"userID" gorm:"not null;uniqueIndex:idx_user_fsp"`
	User      *User    `json:"user" gorm:"foreignKey:UserID"`
	FSPRoleID uint     `json:"fspRoleID" gorm:"not null;index"`
	FSPRole   *FSPRole `json:"fspRole" gorm:"foreignKey:FSPRoleID"`
}

// todo: we might need to reconsider this for different cultures and languages
// todo: for example, in chinese paternal and maternal relatives have different titles
type RelationshipConst string

const (
	RelationshipSelf             RelationshipConst = "self"
	RelationshipFriend           RelationshipConst = "friend"
	RelationshipMother           RelationshipConst = "mother"
	RelationshipFather           RelationshipConst = "father"
	RelationshipBrother          RelationshipConst = "brother"
	RelationshipSister           RelationshipConst = "sister"
	RelationshipSon              RelationshipConst = "son"
	RelationshipDaughter         RelationshipConst = "daughter"
	RelationshipAunt             RelationshipConst = "aunt"
	RelationshipUncle            RelationshipConst = "uncle"
	RelationshipGrandmother      RelationshipConst = "grandmother"
	RelationshipGrandfather      RelationshipConst = "grandfather"
	RelationshipGreatGrandmother RelationshipConst = "great-grandmother"
	RelationshipGreatGrandfather RelationshipConst = "great-grandfather"
	RelationshipNiece            RelationshipConst = "niece"
	RelationshipNephew           RelationshipConst = "nephew"
	RelationshipCousin           RelationshipConst = "cousin"
	RelationshipSecondCousin     RelationshipConst = "second cousin"
	RelationshipStepMother       RelationshipConst = "step-mother"
	RelationshipStepFather       RelationshipConst = "step-father"
	RelationshipStepSon          RelationshipConst = "step-son"
	RelationshipStepdaughter     RelationshipConst = "step-daughter"
)

type UserMemorialRole struct {
	BaseModel
	FSPID          uint              `json:"fspID" gorm:"not null;index"`
	FSP            *FSP              `json:"fsp" gorm:"foreignKey:FSPID"`
	UserID         uint              `json:"userID" gorm:"not null;uniqueIndex:idx_user_memorial"`
	User           *User             `json:"user" gorm:"foreignKey:UserID"`
	MemorialRoleID uint              `json:"memorialRoleID" gorm:"not null;index"`
	MemorialRole   *MemorialRole     `json:"memorialRole" gorm:"foreignKey:MemorialRoleID"`
	MemorialID     uint              `json:"memorialID" gorm:"not null;index;uniqueIndex:idx_user_memorial"`
	Memorial       *Memorial         `json:"memorial" gorm:"foreignKey:MemorialID"`
	Relationship   RelationshipConst `json:"relationship" gorm:"not null" sql:"type:ENUM('self', 'friend', 'mother', 'father', 'brother', 'sister', 'son', 'daughter', 'aunt', 'uncle', 'grandmother', 'grandfather', 'great-grandmother', 'great-grandfather', 'niece', 'nephew', 'cousin', 'second cousin', 'step-mother', 'step-father', 'step-son', 'step-daughter')"`
}

// Invitation ---------------------------------------------------------------

type InvitationTypeConst string

const (
	InvitationTypeFSP      InvitationTypeConst = "fsp"
	InvitationTypeMemorial InvitationTypeConst = "memorial"
)

type InvitationStatusConst string

const (
	InvitationStatusPending  InvitationStatusConst = "pending"
	InvitationStatusAccepted InvitationStatusConst = "accepted"
	InvitationStatusDeclined InvitationStatusConst = "declined"
)

type Invitation struct {
	BaseModel
	FSPID uint `json:"fspID" gorm:"index"`

	// For FSP invitations
	FSPRole FSPRoleConst `json:"fspRole" gorm:"default:'user'"`

	// For Memorial contributor invitations
	MemorialID   *uint             `json:"memorialID" gorm:"index;uniqueIndex:idx_memorial_invitee"`
	MemorialRole MemorialRoleConst `json:"memorialRole" gorm:"default:'contributor'"`

	InvitationType InvitationTypeConst `json:"invitationType" gorm:"not null" sql:"type:ENUM('fsp', 'memorial')"`
	InviterID      uint                `json:"inviterID" gorm:"not null"`
	Inviter        *User               `json:"inviter" gorm:"foreignKey:InviterID"`
	InviteeEmail   string              `json:"inviteeEmail" gorm:"not null;index;uniqueIndex:idx_memorial_invitee"`
	// TODO: Default to self is wrong, need to set to lower roles
	// TODO: Requires proper data migration to before setting 'not null' column issue
	Relationship RelationshipConst     `json:"relationship" gorm:"type:string;not null;default:'self'"`
	Status       InvitationStatusConst `json:"status" gorm:"default:'pending'" sql:"type:ENUM('pending', 'accepted', 'declined')"`
	InvitedOn    time.Time             `json:"invitedOn" gorm:"not null"`
	ExpiresOn    time.Time             `json:"expiresOn" gorm:"not null"`
	Token        string                `json:"-" gorm:"not null;index;unique"`
	IsUsed       bool                  `json:"isUsed" gorm:"default:false"`
}

// Memorial Applications ---------------------------------------------------------------

type ApplicationStatusConst string

const (
	ApplicationStatusPending  ApplicationStatusConst = "pending"
	ApplicationStatusAccepted ApplicationStatusConst = "accepted"
	ApplicationStatusDeclined ApplicationStatusConst = "declined"
)

type ApplicationTypeConst string

const (
	ApplicationTypeContributor ApplicationTypeConst = "contributor"
)

type Application struct {
	BaseModel
	FSPID      uint `json:"fspID" gorm:"not null;uniqueIndex:idx_fsp_applicant"`
	MemorialID uint `json:"memorialID"`

	ApplicantID     uint                   `json:"applicantID" gorm:"not null;uniqueIndex:idx_fsp_applicant"`
	Applicant       *User                  `json:"applicant" gorm:"foreignKey:ApplicantID"`
	Relationship    RelationshipConst      `json:"relationship" gorm:"not null"`
	ApplicationType ApplicationTypeConst   `json:"applicationType" gorm:"not null" sql:"type:ENUM('contributor')"`
	Status          ApplicationStatusConst `json:"status" gorm:"default:'pending'" sql:"type:ENUM('pending', 'accepted', 'declined')"`
	AppliedOn       time.Time              `json:"appliedOn" gorm:"not null"`
}
