package identity

import "github.com/alsey89/people-matter/internal/schema"

type applicationRequest struct {
	MemorialID   int64                    `json:"memorialId" validate:"required"`
	Relationship schema.RelationshipConst `json:"relationship" validate:"required"`
}

type applicationSignupRequest struct {
	MemorialID           int64                    `json:"memorialId" validate:"required"`
	FirstName            string                   `json:"firstName" validate:"required"`
	LastName             string                   `json:"lastName" validate:"required"`
	Relationship         schema.RelationshipConst `json:"relationship" validate:"required"`
	Email                string                   `json:"email" validate:"required,email"`
	Password             string                   `json:"password" validate:"required,min=6"`
	PasswordConfirmation string                   `json:"passwordConfirmation" validate:"required,min=6,eqfield=Password"`
}

type acceptInvitationRequest struct {
	MemorialID int64  `json:"memorialId" validate:"required"`
	Token      string `json:"token" validate:"required"`
}

type invitationSignupRequest struct {
	MemorialID           int64  `json:"memorialId" validate:"required"`
	Token                string `json:"token" validate:"required"`
	Email                string `json:"email" validate:"required,email"`
	FirstName            string `json:"firstName" validate:"required"`
	LastName             string `json:"lastName" validate:"required"`
	Password             string `json:"password" validate:"required,min=6"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"required,min=6,eqfield=Password"`
}

type signinRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type signinResponseData struct {
	User               *schema.User             `json:"user"`
	RolesByLevel       *rolesByLevel            `json:"rolesByLevel"`
	ActiveMemorialRole *schema.UserMemorialRole `json:"activeMemorialRole"`
}

type resetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type confirmResetPasswordRequest struct {
	Password             string `json:"password" validate:"required,min=6"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"required,min=6,eqfield=Password"`
}

// Internal
type rolesByLevel struct {
	Tenant   schema.UserFSPRole
	Memorial []schema.UserMemorialRole
}
