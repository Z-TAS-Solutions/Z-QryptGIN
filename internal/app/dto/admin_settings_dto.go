package dto

// -- Get System Settings --
type SystemSettingsResponse struct {
	Message string `json:"message"`
	Data    struct {
		SessionTimeoutMinutes  int      `json:"sessionTimeoutMinutes"`
		MaxFailedLoginAttempts int      `json:"maxFailedLoginAttempts"`
		AllowedAdminIpRanges   []string `json:"allowedAdminIpRanges"`
	} `json:"data"`
}
