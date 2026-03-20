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

type IPV4 string

var ipv4Regex = regexp.MustCompile(`^((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)\.){3}(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)$`)

func (i IPV4) IsValid() bool {
	return ipv4Regex.MatchString((string(i)))
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

type MfaStatus string

const (
	MfaPending  MfaStatus = "Pending"
	MfaApproved MfaStatus = "Approved"
	MfaDenied   MfaStatus = "Denied"
	MfaExpired  MfaStatus = "Expired"
)

func (ms MfaStatus) IsValid() bool {
	switch ms {
	case MfaPending, MfaApproved, MfaDenied, MfaExpired:
		return true
	}
	return false
}

type MfaDecision string

const (
	MfaDecApproved MfaDecision = "Approved"
	MfaDecDenied   MfaDecision = "Denied"
)

func (md MfaDecision) IsValid() bool {
	switch md {
	case MfaDecApproved, MfaDecDenied:
		return true
	}
	return false
}

type UserRole string

const (
	RoleAdmin  UserRole = "Admin"
	RoleClient UserRole = "Client"
)

func (ur UserRole) IsValid() bool {
	switch ur {
	case RoleAdmin, RoleClient:
		return true
	}
	return false
}
