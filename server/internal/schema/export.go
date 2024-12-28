package schema

type ExportStateConst string

const (
	ExportStateRequested ExportStateConst = "requested" // curator has submitted an export request
	ExportStateRunning   ExportStateConst = "running"   // export request is pending
	ExportStateDeploying ExportStateConst = "deploying" // export request is deploying
	ExportStateComplete  ExportStateConst = "completed" // export request is complete
	ExportStateFailed    ExportStateConst = "failed"    // export request has failed
)

type Export struct {
	BaseModelWithTime
	CompletedAt *string          `json:"completedAt"`
	Retries     int              `json:"retries"`
	FSPID       uint             `json:"fspId"`
	MemorialID  uint             `json:"memorialId"`
	State       ExportStateConst `json:"state" sql:"type:ENUM('requested', 'running', 'deploying', 'completed', 'failed')"`
}
