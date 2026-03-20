package database

import (
	"errors"
	"regexp"
)

// --- Custom Error Definitions ---
var (
	ErrInvalidEmail          = errors.New("invalid email format")
	ErrInvalidPhone          = errors.New("invalid phone number format (use +94 or 07...)")
	ErrInvalidIPv4           = errors.New("invalid IPv4 address format")
	ErrInvalidNIC            = errors.New("invalid NIC format")
	ErrInvalidUserID         = errors.New("invalid User ID format (USR-XXXXXX)")
	ErrInvalidNotificationID = errors.New("invalid Notification ID format (NOT-XXXXXXXX)")
	ErrInvalidMfaID          = errors.New("invalid MFA ID format (MFA-XXXXXXXX)")
	ErrInvalidActivityID     = errors.New("invalid Activity ID format (ACT-XXXXXXXX)")
	ErrInvalidPasskeyID      = errors.New("invalid Passkey ID format (PassK-XXXXXXXX)")
	ErrInvalidSessionID      = errors.New("invalid Session ID format (SESS-XXXXXXXX)")
	ErrInvalidStatus         = errors.New("invalid user status")
	ErrInvalidSecurityLevel  = errors.New("invalid security level")
	ErrInvalidNotifyType     = errors.New("invalid notification type")
	ErrInvalidActivityType   = errors.New("invalid activity log type")
	ErrInvalidMfaStatus      = errors.New("invalid MFA status")
	ErrInvalidMfaDecision    = errors.New("invalid MFA decision")
	ErrInvalidRole           = errors.New("invalid user role")
)

// --- Regex Pre-compilation ---
var (
	emailRegex          = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	phoneRegex          = regexp.MustCompile(`^(\+94|0)?7\d{8}$`)
	ipv4Regex           = regexp.MustCompile(`^((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)\.){3}(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)$`)
	userCustomIDRegex   = regexp.MustCompile(`^USR-[a-zA-Z0-9]{6}$`)
	notificationIDRegex = regexp.MustCompile(`^NOT-[a-zA-Z0-9]{8}$`)
	mfaChallengeIDRegex = regexp.MustCompile(`^MFA-[a-zA-Z0-9]{8}$`)
	activityLogIDRegex  = regexp.MustCompile(`^ACT-[a-zA-Z0-9]{8}$`)
	passkeyIDRegex      = regexp.MustCompile(`^PassK-[a-zA-Z0-9]{8}$`)
	sessionIDRegex      = regexp.MustCompile(`^SESS-[a-zA-Z0-9]{8}$`)
	nicRegex            = regexp.MustCompile(`^([0-9]{9}[vVxX]|[0-9]{12})$`)
)

// --- String-Based Custom Types ---

type Email string

func (e Email) Validate() error {
	if !emailRegex.MatchString(string(e)) {
		return ErrInvalidEmail
	}
	return nil
}

type PhoneNumber string

func (p PhoneNumber) Validate() error {
	if !phoneRegex.MatchString(string(p)) {
		return ErrInvalidPhone
	}
	return nil
}

type IPV4 string

func (i IPV4) Validate() error {
	if !ipv4Regex.MatchString(string(i)) {
		return ErrInvalidIPv4
	}
	return nil
}

type NIC string

func (n NIC) Validate() error {
	if !nicRegex.MatchString(string(n)) {
		return ErrInvalidNIC
	}
	return nil
}

type UserCustomID string

func (uci UserCustomID) Validate() error {
	if !userCustomIDRegex.MatchString(string(uci)) {
		return ErrInvalidUserID
	}
	return nil
}

type NotificationID string

func (ni NotificationID) Validate() error {
	if !notificationIDRegex.MatchString(string(ni)) {
		return ErrInvalidNotificationID
	}
	return nil
}

type MfaChallengeID string

func (mci MfaChallengeID) Validate() error {
	if !mfaChallengeIDRegex.MatchString(string(mci)) {
		return ErrInvalidMfaID
	}
	return nil
}

type ActivityID string

func (aI ActivityID) Validate() error {
	if !activityLogIDRegex.MatchString(string(aI)) {
		return ErrInvalidActivityID
	}
	return nil
}

type PasskeyID string

func (pi PasskeyID) Validate() error {
	if !passkeyIDRegex.MatchString(string(pi)) {
		return ErrInvalidPasskeyID
	}
	return nil
}

type SessionID string

func (si SessionID) Validate() error {
	if !sessionIDRegex.MatchString(string(si)) {
		return ErrInvalidSessionID
	}
	return nil
}

// --- Enum-Based Custom Types ---

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
	return ErrInvalidStatus
}

type UserSecurityLevel string

const (
	SecurityLow    UserSecurityLevel = "Low"
	SecurityMedium UserSecurityLevel = "Medium"
	SecurityHigh   UserSecurityLevel = "High"
)

func (sl UserSecurityLevel) Validate() error {
	switch sl {
	case SecurityLow, SecurityMedium, SecurityHigh:
		return nil
	}
	return ErrInvalidSecurityLevel
}

type NotificationType string

const (
	NotifySecurity NotificationType = "Security"
	NotifyInfo     NotificationType = "Info"
	NotifySuccess  NotificationType = "Success"
	NotifyError    NotificationType = "Error"
	NotifyWarning  NotificationType = "Warning"
)

func (nt NotificationType) Validate() error {
	switch nt {
	case NotifySecurity, NotifyInfo, NotifySuccess, NotifyError, NotifyWarning:
		return nil
	}
	return ErrInvalidNotifyType
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

func (alt ActivityLogType) Validate() error {
	switch alt {
	case ActivityFailedLogin, ActivityLoginSuccess, ActivityLogout, ActivityMFAApproved, ActivityMFADenied, ActivityPassKeyActivated, ActivitySessionRevoked:
		return nil
	}
	return ErrInvalidActivityType
}

type MfaStatus string

const (
	MfaPending  MfaStatus = "Pending"
	MfaApproved MfaStatus = "Approved"
	MfaDenied   MfaStatus = "Denied"
	MfaExpired  MfaStatus = "Expired"
)

func (ms MfaStatus) Validate() error {
	switch ms {
	case MfaPending, MfaApproved, MfaDenied, MfaExpired:
		return nil
	}
	return ErrInvalidMfaStatus
}

type MfaDecision string

const (
	MfaDecApproved MfaDecision = "Approved"
	MfaDecDenied   MfaDecision = "Denied"
)

func (md MfaDecision) Validate() error {
	switch md {
	case MfaDecApproved, MfaDecDenied:
		return nil
	}
	return ErrInvalidMfaDecision
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
	return ErrInvalidRole
}
