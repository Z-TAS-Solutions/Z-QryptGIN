package dto

// -- List Devices --
type ListDevicesRequest struct {
	Limit          int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset         int    `form:"offset" binding:"omitempty,min=0"`
	Search         string `form:"search"`
	Location       string `form:"location"`
	LastActiveFrom int64  `form:"last_active_from"`
	LastActiveTo   int64  `form:"last_active_to"`
	SortBy         string `form:"sort_by"`
	SortOrder      string `form:"sort_order"`
}

type DeviceItem struct {
	DeviceID   string `json:"device_id"`
	DeviceName string `json:"device_name"`
	Location   string `json:"location"`
	LastActive int64  `json:"last_active"`
}

type ListDevicesResponse struct {
	Message string `json:"message"`
	Data    struct {
		Devices    []DeviceItem `json:"devices"`
		Pagination struct {
			Limit    int  `json:"limit"`
			Offset   int  `json:"offset"`
			Returned int  `json:"returned"`
			HasMore  bool `json:"has_more"`
		} `json:"pagination"`
	} `json:"data"`
}

// -- Force Logout Device --
type ForceLogoutDeviceResponse struct {
	Message string `json:"message"`
	Data    struct {
		DeviceID  string `json:"device_id"`
		LoggedOut bool   `json:"logged_out"`
	} `json:"data"`
}
