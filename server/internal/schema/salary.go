package schema

import (
	"time"

	"gorm.io/gorm"
)

type Salary struct {
	gorm.Model
	UserID          uint      `json:"userId"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	PaymentInterval string    `json:"paymentInterval"`
	EffectiveDate   time.Time `json:"effectiveDate"`
	IsActive        bool      `json:"isActive"`
}

type Payment struct {
	gorm.Model
	UserID        uint      `json:"userId"`
	SalaryID      uint      `json:"salaryId"`
	PaymentDate   string    `json:"paymentDate"`
	Amount        float64   `json:"amount"`
	PaymentMethod string    `json:"paymentMethod"`
	Status        string    `json:"status"`
	PeriodStart   time.Time `json:"periodStart"`
	PeriodEnd     time.Time `json:"periodEnd"`
	Deductions    float64   `json:"deductions"`
	Bonuses       float64   `json:"bonuses"`
	Notes         string    `json:"notes"`
}

func (salary *Salary) AfterCreate(tx *gorm.DB) (err error) {
	// Update User's SalaryID field
	err = tx.Model(&User{}).Where("id = ?", salary.UserID).Update("current_salary_info_id", salary.ID).Error
	if err != nil {
		return err
	}
	return nil
}
