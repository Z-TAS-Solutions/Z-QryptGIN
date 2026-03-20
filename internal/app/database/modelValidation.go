package database

import "fmt"

import (
	"errors"
	"regexp"
)

// --- Domain Types ---

type Email string

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func (e Email) Validate() error {
	if !emailRegex.MatchString(string(e)) {
		return errors.New("invalid email format")
	}
	return nil
}

type PhoneNumber string

var phoneRegex = regexp.MustCompile(`^(\+94|0)?7\d{8}$`)

func (p PhoneNumber) Validate() error {
	if !phoneRegex.MatchString(string(p)) {
		return errors.New("invalid phone number format: must be Sri Lankan standard")
	}
	return nil
}

type IPV4 string

var ipv4Regex = regexp.MustCompile(`^((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)\.){3}(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)$`)

func (i IPV4) Validate() error {
	if !ipv4Regex.MatchString(string(i)) {
		return errors.New("invalid IPv4 address")
	}
	return nil
}

type NIC string

var nicRegex = regexp.MustCompile(`^([0-9]{9}[vVxX]|[0-9]{12})$`)

func (n NIC) Validate() error {
	if !nicRegex.MatchString(string(n)) {
		return errors.New("invalid NIC format")
	}
	return nil
}

// --- ID Types with Specific Error Messages ---

type UserCustomID string

var userCustomIDRegex = regexp.MustCompile(`^USR-[a-zA-Z0-9]{6}$`)

func (id UserCustomID) Validate() error {
	if !userCustomIDRegex.MatchString(string(id)) {
		return errors.New("invalid User ID format (expected USR-XXXXXX)")
	}
	return nil
}

type NotificationID string

var notificationIDRegex = regexp.MustCompile(`^NOT-[a-zA-Z0-9]{8}$`)

func (id NotificationID) Validate() error {
	if !notificationIDRegex.MatchString(string(id)) {
		return errors.New("invalid Notification ID format")
	}
	return nil
}

type SessionID string

var sessionIDRegex = regexp.MustCompile(`^SESS-[a-zA-Z0-9]{8}$`)

func (id SessionID) Validate() error {
	if !sessionIDRegex.MatchString(string(id)) {
		return errors.New("invalid Session ID format")
	}
	return nil
}

// ... Follow this pattern for MfaChallengeID, ActivityID, and PasskeyID ...
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
