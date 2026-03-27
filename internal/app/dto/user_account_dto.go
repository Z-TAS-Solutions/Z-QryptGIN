package dto

// -- Force Logout Devices (Admin action per User) --
type ForceLogoutUserDevicesRequest struct {
	UserID string `json:"userId"`
}

type ForceLogoutUserDevicesResponse struct {
	Message string `json:"message"`
	Data    struct {
		UserID            string `json:"user_id"`
		DevicesTerminated int    `json:"devices_terminated"`
	} `json:"data"`
}

// -- Delete Account --
type DeleteAccountRequest struct {
	Password string `json:"password" binding:"required"`
}

type DeleteAccountResponse struct {
	Message string `json:"message"`
}
