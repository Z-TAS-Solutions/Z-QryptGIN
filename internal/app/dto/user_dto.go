package dto

import "time"

// ListUsersRequest handles query parameters for listing users
type ListUsersRequest struct {
	Limit         int        `form:"limit,default=20" binding:"min=1,max=100"`
	Offset        int        `form:"offset,default=0" binding:"min=0"`
	Search        string     `form:"search"`                                      // Used for querying nae, email, phone
	Status        string     `form:"status"`                                      // e.g. "active", "inactive", "suspended"
	MfaEnabled    *bool      `form:"mfaEnabled"`                                  // Optional pointer to allow strict true/false testing
	LastLoginFrom *time.Time `form:"lastLoginFrom" time_format:"2006-01-02T15:04:05Z07:00"` // ISO 8601 formatting
	LastLoginTo   *time.Time `form:"lastLoginTo" time_format:"2006-01-02T15:04:05Z07:00"`
	SortBy        string     `form:"sortBy,default=lastLogin"`
	Order         string     `form:"order,default=desc" binding:"oneof=asc desc"`
}

// UserDashboardResponse is the shape of a single user entry inside the ListUsers response list
type UserDashboardResponse struct {
	UserID     string     `json:"userId"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	Phone      string     `json:"phone"`
	MfaEnabled bool       `json:"mfaEnabled"`
	LastLogin  *time.Time `json:"lastLogin"` // Omitted if null
	Status     string     `json:"status"`
}

// ListUsersResponse wraps multiple users with pagination metadata
type ListUsersResponse struct {
	Limit  int                     `json:"limit"`
	Offset int                     `json:"offset"`
	Total  int64                   `json:"total"`
	Data   []UserDashboardResponse `json:"data"`
}
