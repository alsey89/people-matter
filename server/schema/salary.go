package schema

import (
	"time"

	"gorm.io/gorm"
)

type SalaryPaymentIntervalEnum string

const (
	Monthly   SalaryPaymentIntervalEnum = "monthly"
	Weekly    SalaryPaymentIntervalEnum = "weekly"
	Quarterly SalaryPaymentIntervalEnum = "quarterly"
)

type Salary struct {
	BaseModel
	CompanyID uint `json:"company_id" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	// ------------------------------------------------------------------------------------------------
	UserID          uint                      `json:"userId"          gorm:"onUpdate:CASCADE;onDelete:CASCADE"`
	Amount          float64                   `json:"amount"          gorm:"type:decimal(10,2)"`
	Currency        string                    `json:"currency"        gorm:"type:varchar(5)"`
	PaymentInterval SalaryPaymentIntervalEnum `json:"paymentInterval" sql:"type:enum('monthly','weekly','quarterly')"`
	EffectiveDate   time.Time                 `json:"effectiveDate"   gorm:"type:date"`
	// ------------------------------------------------------------------------------------------------
	ApprovalStatus ApprovalStatusEnum `json:"approvalStatus"    sql:"type:enum('pending','approved','rejected');not null"`
}

type PaymentMethodEnum string

const (
	Cash   PaymentMethodEnum = "cash"
	Bank   PaymentMethodEnum = "bank"
	Cheque PaymentMethodEnum = "cheque"
)

type PaymentStatusEnum string

const (
	Approving  PaymentStatusEnum = "approving"
	Processing PaymentStatusEnum = "processing"
	Completed  PaymentStatusEnum = "completed"
)

type Payment struct {
	BaseModel
	CompanyID uint `json:"company_id" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	// ------------------------------------------------------------------------------------------------
	UserID        uint              `json:"userId"         gorm:"onUpdate:CASCADE;onDelete:CASCADE"`
	SalaryID      uint              `json:"salaryId"       gorm:"onUpdate:CASCADE;onDelete:SET NULL"`
	PaymentDate   time.Time         `json:"paymentDate"    gorm:"type:date"`
	Amount        float64           `json:"amount"         gorm:"type:decimal(10,2)"`
	PaymentMethod PaymentMethodEnum `json:"paymentMethod"  sql:"type:enum('cash','bank','cheque')"`
	PeriodStart   time.Time         `json:"periodStart"    gorm:"type:date"`
	PeriodEnd     time.Time         `json:"periodEnd"      gorm:"type:date"`
	Adjustments   []*Adjustments    `json:"adjustments"    gorm:"foreignKey:PaymentID"`
	Notes         string            `json:"notes"          gorm:"type:text"`
	//todo: figure out relationship Documents     []*Document       `json:"documents"`
	// ------------------------------------------------------------------------------------------------
	ApprovalStatus ApprovalStatusEnum `json:"approvalStatus" sql:"type:enum('pending','approved','rejected')"`
}

type AdjustmentTypeEnum string

const (
	Bonus     AdjustmentTypeEnum = "bonus"
	Deduction AdjustmentTypeEnum = "deduction"
)

type Adjustments struct {
	BaseModel
	CompanyID uint `json:"company_id" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	// ------------------------------------------------------------------------------------------------
	PaymentID      uint               `json:"paymentId" gorm:"onUpdate:CASCADE;onDelete:CASCADE"`
	Amount         float64            `json:"amount"    gorm:"type:decimal(10,2)"`
	AdjustmentType AdjustmentTypeEnum `json:"type"      sql:"type:enum('bonus','deduction')"`
	Notes          string             `json:"notes"     gorm:"type:text"`
	// ------------------------------------------------------------------------------------------------
}

func (salary *Salary) AfterCreate(tx *gorm.DB) (err error) {
	// Update User's SalaryID field
	err = tx.Model(&User{}).Where("id = ?", salary.UserID).Update("current_salary_info_id", salary.ID).Error
	if err != nil {
		return err
	}
	return nil
}
