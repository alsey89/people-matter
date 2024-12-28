package response

import (
	"github.com/alsey89/people-matter/internal/common/errmgr"
)

type APIResponse struct {
	Message string             `json:"message"`
	Data    interface{}        `json:"data"`
	Error   errmgr.ErrResponse `json:"error"`
}

type PayloadDebug struct {
	TraceId  string      `json:"traceId"`
	Expected interface{} `json:"expected"`
	Actual   interface{} `json:"actual"`
}
