package schema

type FSPTypeConst string

const (
	FSPTypeFuneralHome FSPTypeConst = "funeral home"
	FSPTypeCemetery    FSPTypeConst = "cemetery"
	FSPTypeMonument    FSPTypeConst = "monument"
	FSPTypeReseller    FSPTypeConst = "reseller"
	FSPTypeDistributor FSPTypeConst = "distributor"
	FSPTypeOther       FSPTypeConst = "other"
)

type BusinessTypeConst string

const (
	BusinessTypeCorporation BusinessTypeConst = "corporation"
	BusinessTypeLLC         BusinessTypeConst = "llc"
	BusinessTypeSoleProp    BusinessTypeConst = "sole proprietorship"
	BusinessTypePartnership BusinessTypeConst = "partnership"
	BusinessTypeOther       BusinessTypeConst = "other"
)

type FSP struct {
	BaseModel
	TenantIdentifier string `json:"-" gorm:"not null;index"`
	// Company Information ---------------------
	Name           string            `json:"name"`
	LogoURL        string            `json:"logoUrl"`
	FSPType        FSPTypeConst      `json:"fspType" sql:"type:ENUM('funeral home', 'cemetery', 'monument', 'reseller', 'distributor', 'other')" gorm:"not null"`
	BusinessTypeID int               `json:"businessTypeId"`
	BusinessType   BusinessTypeConst `json:"businessType" sql:"type:ENUM('corporation', 'llc', 'sole proprietorship', 'partnership', 'other')" gorm:"not null"`
	CRN            string            `json:"crn"`
	EIN            string            `json:"ein"`
	Established    string            `json:"established"`
	EmployeeCount  string            `json:"employeeCount"`
	ParentCompany  string            `json:"parentCompany"`
	Subsidiaries   string            `json:"subsidiaries"`
	// Contact Information ---------------------
	Email           string        `json:"email"`
	Phone           string        `json:"phone"`
	Website         string        `json:"website"`
	Address         string        `json:"address"`
	PostalCode      string        `json:"postalCode"`
	StateProvinceID int           `json:"stateProvinceId"`
	StateProvince   StateProvince `json:"stateProvince" gorm:"foreignKey:StateProvinceID" `
	CountryID       int           `json:"countryId"`
	Country         Country       `json:"country" gorm:"foreignKey:CountryID"`
	// Billing Information ---------------------
	BillingAddress string `json:"billingAddress"`
	// Account Information ---------------------
	MemorialQuota     int     `json:"memorialQuota"`
	MemorialQuotaUsed int     `json:"memorialQuotaUsed"`
	StorageQuota      float64 `json:"storageQuota"`
	StorageQuotaUsed  float64 `json:"storageQuotaUsed"`
	// Associations
	Users     *[]User     `json:"users" gorm:"foreignKey:FSPID"`
	Memorials *[]Memorial `json:"memorials" gorm:"foreignKey:FSPID"`
}

type Country struct {
	BaseModel
	Name string `json:"name"`
	Code string `json:"code"`
}

type StateProvince struct {
	BaseModel
	Name      string   `json:"name"`
	Code      string   `json:"code"`
	CountryID int      `json:"countryId"`
	Country   *Country `json:"country" gorm:"foreignKey:CountryID"`
}
