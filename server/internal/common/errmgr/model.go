package errmgr

type APIError struct {
	TraceID string `json:"traceId"` // A unique identifier for the request
	Code    string `json:"code"`    // e.g., "ERR_CODE_INTERNAL"
	Status  int    `json:"status"`  // e.g., 500
}
