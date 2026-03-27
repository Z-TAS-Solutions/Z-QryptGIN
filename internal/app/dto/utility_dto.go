package dto

type PingResponse struct {
	Message string `json:"message"`
	Data    struct {
		ServerTimestamp int64 `json:"server_timestamp"`
	} `json:"data"`
}
