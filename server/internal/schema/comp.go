package schema

import "time"

type compIntervalConst string

const (
	compIntervalYearly    compIntervalConst = "yearly"
	compIntervalQuarterly compIntervalConst = "quarterly"
	compIntervalMonthly   compIntervalConst = "monthly"
	compIntervalWeekly    compIntervalConst = "weekly"
	compIntervalDaily     compIntervalConst = "daily"
	compIntervalOther     compIntervalConst = "other"
)

type Comp struct {
	BaseModel
	TenantIdentifier string            `json:"-" gorm:"not null;index"`
	Amount           float64           `json:"amount"`
	CompInterval     compIntervalConst `json:"compInterval" sql:"type:compIntervalConst['yearly','quarterly','monthly','weekly','daily','other']"`
	UserID           uint              `json:"userId" gorm:"not null"`
	// Associations
	User User `json:"user" gorm:"foreignKey:UserID"`
}

type PaymentChannelTypeConst string

const (
	PaymentRecipientBank  PaymentChannelTypeConst = "bank"
	PaymentRecipientOther PaymentChannelTypeConst = "other"
)

type PaymentAccount struct {
	BaseModel
	TenantIdentifier string `json:"-" gorm:"not null;index"`
	// Payment Information ---------------------
	UserID         uint                    `json:"userId" gorm:"not null"`
	PaymentChannel PaymentChannelTypeConst `json:"paymentChannel" sql:"type:PaymentRecipientConst['bank','other']"`
	ProviderID     *string                 `json:"providerId"`
	AccountNumber  string                  `json:"accountNumber"`
	AccountName    *string                 `json:"accountName"`
	AccountEmail   *string                 `json:"accountEmail"`
}

type Payment struct {
	BaseModel
	TenantIdentifier string `json:"-" gorm:"not null;index"`
	// Payment Information ---------------------
	CompID           uint       `json:"compId" gorm:"not null"`
	PaymentAccountID uint       `json:"paymentAccountId" gorm:"not null"`
	Amount           float64    `json:"amount"`
	Bonus            float64    `json:"bonus"`
	PaidAt           *time.Time `json:"paidAt"`
	// Associations
	PaymentAccount PaymentAccount `json:"paymentAccount" gorm:"foreignKey:PaymentAccountID"`
	Comp           Comp           `json:"comp" gorm:"foreignKey:CompID"`
}
