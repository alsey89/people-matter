package schema

type LeaveDocument struct {
	LeaveID    uint `json:"leaveId"    gorm:"primaryKey;not null"`
	DocumentID uint `json:"documentId" gorm:"primaryKey;not null"`
}
