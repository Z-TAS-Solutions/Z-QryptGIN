package dto

// CreateUserRequest is what the client sends
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	PhoneNo  string `json:"phone_no" binding:"required"`
	Nic      string `json:"nic" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// UserResponse is what the client receives (no passwords!)
type UserResponse struct {
	Success  bool   `json:"success"`
	CustomID string `json:"custom_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}
