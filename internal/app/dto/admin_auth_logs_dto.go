package dto

// -- Authentication Logs --
type AuthLogsRequest struct {
	Limit     int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset    int    `form:"offset" binding:"omitempty,min=0"`
	Search    string `form:"search"`
	Status    string `form:"status"`
	From      int64  `form:"from"`
	To        int64  `form:"to"`
	SortOrder string `form:"sort_order"`
}

type AuthLogItem struct {
	LogID     string `json:"log_id"`
	Timestamp int64  `json:"timestamp"`
	UserName  string `json:"user_name"`
	UserID    string `json:"user_id"`
	Status    string `json:"status"`
	Location  string `json:"location"`
	Device    string `json:"device"`
}

type AuthLogsResponse struct {
	Message string `json:"message"`
	Data    struct {
		Logs       []AuthLogItem `json:"logs"`
		Pagination struct {
			Limit    int  `json:"limit"`
			Offset   int  `json:"offset"`
			Returned int  `json:"returned"`
			HasMore  bool `json:"has_more"`
		} `json:"pagination"`
	} `json:"data"`
}

// -- Authentication Analytics --
type AuthAnalyticsRequest struct {
	From int64 `form:"from"`
	To   int64 `form:"to"`
}

type AuthAnalyticsResponse struct {
	Message string `json:"message"`
	Data    struct {
		TimeRange struct {
			From int64 `json:"from"`
			To   int64 `json:"to"`
		} `json:"time_range"`
		Metrics struct {
			SuccessfulLogins     int64 `json:"successful_logins"`
			FailedLogins         int64 `json:"failed_logins"`
			SuspiciousActivities int64 `json:"suspicious_activities"`
		} `json:"metrics"`
	} `json:"data"`
}
