package database

import "fmt"

type UserStatus string

const (
	StatusActive   UserStatus = "Active"
	StatusInactive UserStatus = "Inactive"
)

func (us UserStatus) Validate() error {
	switch us {
	case StatusActive, StatusInactive:
		return nil
	}
	return fmt.Errorf("invalid status: %s", us)
}

type UserRole string

const (
	RoleAdmin  UserRole = "Admin"
	RoleClient UserRole = "Client"
)

func (ur UserRole) Validate() error {
	switch ur {
	case RoleAdmin, RoleClient:
		return nil
	}
	return fmt.Errorf("invalid role: %s", ur)
}
