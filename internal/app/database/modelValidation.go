package database

import "regexp"

type Email string

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func (e Email) IsValid() bool {
	return emailRegex.MatchString(string(e))
}

type PhoneNumber string

var phoneRegex = regexp.MustCompile(`^(\+94|0)?7\d{8}$`)

func (p PhoneNumber) IsValid() bool {
	return phoneRegex.MatchString(string(p))
}

type UserStatus string

const (
	StatusActive   UserStatus = "Active"
	StatusInactive UserStatus = "Inactive"
)

func (us UserStatus) IsValid() bool {
	switch us {
	case StatusActive, StatusInactive:
		return true
	}
	return false
}

type UserSecurityLevel string

const (
	SecurityLow    UserSecurityLevel = "Low"
	SecurityMedium UserSecurityLevel = "Medium"
	SecurityHigh   UserSecurityLevel = "High"
)

func (sl UserSecurityLevel) IsValid() bool {
	switch sl {
	case SecurityLow, SecurityMedium, SecurityHigh:
		return true
	}
	return false
}

type NotificationType string

const (
	NotifySecurity NotificationType = "Security"
	NotifyInfo     NotificationType = "Info"
	NotifySuccess  NotificationType = "Success"
	NotifyError    NotificationType = "Error"
	NotifyWarning  NotificationType = "Warning"
)

func (nt NotificationType) IsValid() bool {
	switch nt {
	case NotifySecurity, NotifyInfo, NotifySuccess, NotifyError, NotifyWarning:
		return true
	}
	return false
}

type ActivityLogType string

const (
	ActivityFailedLogin      ActivityLogType = "Failed_Login"
	ActivityMFAApproved      ActivityLogType = "MFA_Approved"
	ActivityMFADenied        ActivityLogType = "MFA_Denied"
	ActivityLoginSuccess     ActivityLogType = "Login_Success"
	ActivityLogout           ActivityLogType = "Logout"
	ActivityPassKeyActivated ActivityLogType = "Passkey_Activated"
	ActivitySessionRevoked   ActivityLogType = "Session_Revoked"
)

func (alt ActivityLogType) IsValid() bool {
	switch alt {
	case ActivityFailedLogin, ActivityLoginSuccess, ActivityLogout, ActivityMFAApproved, ActivityMFADenied, ActivityPassKeyActivated, ActivitySessionRevoked:
		return true
	}
	return false
}
