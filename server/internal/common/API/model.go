package API

import (
	"github.com/alsey89/people-matter/internal/common/errmgr"
)

type Response struct {
	Message    string          `json:"message"`
	Data       interface{}     `json:"data"`
	Pagination *Pagination     `json:"pagination"`
	Error      errmgr.APIError `json:"error"`
}

type Pagination struct {
	Page    int `json:"page"`
	PerPage int `json:"perPage"`
	Total   int `json:"total"`
}

type PayloadDebug struct {
	TraceId  string      `json:"traceId"`
	Expected interface{} `json:"expected"`
	Actual   interface{} `json:"actual"`
}
