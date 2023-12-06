package shared

type APIResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
