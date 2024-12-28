package fsp

import "github.com/alsey89/people-matter/internal/schema"

type postTeamRequest struct {
	Email string              `json:"email" validate:"required,email"`
	Role  schema.FSPRoleConst `json:"role" validate:"required"`
}

type putTeamRequest struct {
	Role schema.FSPRoleConst `json:"role" validate:"required"`
}

type deleteTeamRequest struct {
	Email      string `json:"email" validate:"required,email"`
	NotifyUser bool   `json:"userIsNotified"`
}

type postMemorialRequest struct {
	FirstName           string                   `json:"firstName" validate:"required"`
	LastName            string                   `json:"lastName" validate:"required"`
	DOBString           *string                  `json:"dobString" validate:"required"`
	DODString           *string                  `json:"dodString"`
	CuratorEmail        string                   `json:"curatorEmail" validate:"required,email"`
	CuratorRelationship schema.RelationshipConst `json:"curatorRelationship" validate:"required"`
}

type putMemorialRequest struct {
	Title     *string `json:"title"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	DOBString *string `json:"dobString"`
	DODString *string `json:"dodString"`
}
