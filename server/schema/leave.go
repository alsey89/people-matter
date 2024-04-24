package schema

import (
	"time"

	"gorm.io/gorm"
)

type Leave struct {
	BaseModel
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
