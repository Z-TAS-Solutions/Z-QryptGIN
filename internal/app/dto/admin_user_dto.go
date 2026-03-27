package dto

import "time"

// -- List Users --
type ListUsersRequest struct {
	Limit         int    `form:"limit" binding:"omitempty,min=1"`
	Offset        int    `form:"offset" binding:"omitempty,min=0"`
	Search        string `form:"search"`
	Status        string `form:"status"`
	MfaEnabled    *bool  `form:"mfaEnabled"`
	LastLoginFrom string `form:"lastLoginFrom"`
	LastLoginTo   string `form:"lastLoginTo"`
	SortBy        string `form:"sortBy"`
	Order         string `form:"order"`
}

type UserListItem struct {
	UserID     string     `json:"userId"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	Phone      string     `json:"phone"`
	MfaEnabled bool       `json:"mfaEnabled"`
	LastLogin  *time.Time `json:"lastLogin"`
	Status     string     `json:"status"`
}

type ListUsersResponse struct {
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
	Total  int64          `json:"total"`
	Data   []UserListItem `json:"data"`
}

// -- Get User Details --
type UserDetails struct {
	UserID       string    `json:"userId"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	MfaEnabled   bool      `json:"mfaEnabled"`
	MfaMethods   []string  `json:"mfaMethods"`
	Status       string    `json:"status"`
	RegisteredAt time.Time `json:"registeredAt"`
}

type UserDetailsResponse struct {
	Message string      `json:"message"`
	Data    UserDetails `json:"data"`
}

// -- Lock/Unlock User --
type LockUserRequest struct {
	Locked bool `json:"locked"`
}

type LockUserResponse struct {
	Message string `json:"message"`
	Data    struct {
		UserID string `json:"userId"`
		Locked bool   `json:"locked"`
	} `json:"data"`
}

// -- Delete User --
type DeleteUserResponse struct {
	Message string `json:"message"`
	Data    struct {
		UserID  string `json:"userId"`
		Deleted bool   `json:"deleted"`
	} `json:"data"`
}
