package schema

import (
	"time"
)

type LeaveTypeEnum string

const (
	Sick        LeaveTypeEnum = "sick"
	Vacation    LeaveTypeEnum = "vacation"
	Maternity   LeaveTypeEnum = "maternity"
	Paternity   LeaveTypeEnum = "paternity"
	Bereavement LeaveTypeEnum = "bereavement"
	Other       LeaveTypeEnum = "other"
)

type Leave struct {
	BaseModel
	CompanyID uint `json:"company_id" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	// ------------------------------------------------------------------------------------------------
	UserID    uint          `json:"user_id"           gorm:"onUpdate:CASCADE onDelete:CASCADE"`
	Type      LeaveTypeEnum `json:"type"              sql:"type:enum('sick','vacation','maternity','paternity','bereavement','other'); not null"`
	StartDate time.Time     `json:"startDate"         gorm:"type:date;not null"`
	EndDate   time.Time     `json:"endDate"           gorm:"type:date;not null"`
	Reason    string        `json:"reason"            gorm:"type:text"`
	//todo: figure out relationship Documents []Document    `json:"documents"`
	// ------------------------------------------------------------------------------------------------
	ApprovalStatus ApprovalStatusEnum `json:"approvalStatus"    sql:"type:enum('pending','approved','rejected');not null"`
}
