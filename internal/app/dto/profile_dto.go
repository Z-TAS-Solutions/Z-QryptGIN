package dto

type UserProfileResponse struct {
	Message string `json:"message"`
	Data    struct {
		UserID              string `json:"userId"`
		Name                string `json:"name"`
		Email               string `json:"email"`
		Phone               string `json:"phone"`
		MFAEnabled          bool   `json:"mfaEnabled"`
		LinkedPasskeysCount int    `json:"linkedPasskeysCount"`
	} `json:"data"`
}

type UpdateProfileRequest struct {
	Name  string `json:"name" binding:"omitempty"`
	Phone string `json:"phone" binding:"omitempty"`
}

type UpdateProfileResponse struct {
	Message string `json:"message"`
	Data    struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	} `json:"data"`
}
