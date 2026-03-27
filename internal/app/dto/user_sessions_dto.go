package dto

// -- Fetch Active Sessions --
type FetchSessionsRequest struct {
	Limit  int `form:"limit" binding:"omitempty,min=1"`
	Offset int `form:"offset" binding:"omitempty,min=0"`
}

type SessionItem struct {
	SessionID  string `json:"session_id"`
	DeviceID   string `json:"device_id"`
	DeviceName string `json:"device_name"`
	Location   string `json:"location"`
	IpAddress  string `json:"ip_address"`
	LastActive int64  `json:"last_active"`
	Current    bool   `json:"current"`
}

type FetchSessionsResponse struct {
	Message string `json:"message"`
	Data    struct {
		Sessions   []SessionItem `json:"sessions"`
		Pagination struct {
			Limit    int  `json:"limit"`
			Offset   int  `json:"offset"`
			Returned int  `json:"returned"`
			HasMore  bool `json:"has_more"`
		} `json:"pagination"`
	} `json:"data"`
}

// -- Logout Others --
type LogoutOtherSessionsResponse struct {
	Message string `json:"message"`
	Data    struct {
		SessionsTerminated int `json:"sessions_terminated"`
	} `json:"data"`
}
