package dto

// -- Admin Login --
type AdminLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AdminLoginResponse struct {
	Message string `json:"message"`
	Data    struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		AdminID      string `json:"adminId"`
		Role         string `json:"role"`
	} `json:"data"`
}

// -- Admin Token Refresh --
type AdminRefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type AdminRefreshResponse struct {
	Message string `json:"message"`
	Data    struct {
		AccessToken string `json:"accessToken"`
	} `json:"data"`
}
