package database

import (
	"errors"
	"regexp"
)

// --- Custom Error Definitions ---
var (
	ErrInvalidName           = errors.New("name cannot be empty or have special characters")
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
	ErrInvalidTemplateID     = errors.New("template ID cannot be empty")
	ErrInvalidNonce          = errors.New("invalid nonce length (expected 12 bytes)")
	ErrInvalidCrypticData    = errors.New("cryptic data (DEK/Ciphertext) cannot be empty")
	ErrInvalidCredentialID   = errors.New("credential ID cannot be empty")
	ErrInvalidPublicKey      = errors.New("public key cannot be empty")
	ErrInvalidAttestationType = errors.New("invalid attestation type")
	ErrInvalidTransport      = errors.New("transport data cannot be empty")
	ErrInvalidAAGUID         = errors.New("invalid AAGUID length (expected 16 bytes)")
	ErrInvalidAuthenticatorName = errors.New("authenticator name cannot be empty or too long (max 100 chars)")
)

// --- Regex Pre-compilation ---
var (
	nameRegex           = regexp.MustCompile(`^[a-zA-Z\s]+$`)
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
	attestationTypeRegex = regexp.MustCompile(`^[a-z0-9\-]+$`)
)

// --- String-Based Custom Types ---

type Name string

func (n Name) Validate() error {
	if !nameRegex.MatchString(string(n)) {
		return ErrInvalidName
	}
	return nil
}

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

type TemplateID string

func (t TemplateID) Validate() error {
	if len(t) == 0 {
		return ErrInvalidTemplateID
	}
	return nil
}

type Nonce12 []byte

func (n Nonce12) Validate() error {
	if len(n) != 12 {
		return ErrInvalidNonce
	}
	return nil
}

type CrypticData []byte

func (cd CrypticData) Validate() error {
	if len(cd) == 0 {
		return ErrInvalidCrypticData
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

// --- WebAuthnCredential Custom Types ---

type CredentialID []byte

func (ci CredentialID) Validate() error {
	if len(ci) == 0 {
		return ErrInvalidCredentialID
	}
	return nil
}

type WebAuthnPublicKey []byte

func (pk WebAuthnPublicKey) Validate() error {
	if len(pk) == 0 {
		return ErrInvalidPublicKey
	}
	return nil
}

type AttestationType string

const (
	AttestationNone   AttestationType = "none"
	AttestationFidoU2F AttestationType = "fido-u2f"
	AttestationPacked  AttestationType = "packed"
	AttestationAndroidKey AttestationType = "android-key"
	AttestationAndroidSafetyNet AttestationType = "android-safetynet"
	AttestationTPM AttestationType = "tpm"
	AttestationAppleAnonymousCA AttestationType = "apple-anonymous-ca"
)

func (at AttestationType) Validate() error {
	switch at {
	case AttestationNone, AttestationFidoU2F, AttestationPacked, AttestationAndroidKey, AttestationAndroidSafetyNet, AttestationTPM, AttestationAppleAnonymousCA, "":
		return nil
	}
	return ErrInvalidAttestationType
}

type WebAuthnTransport []byte

func (t WebAuthnTransport) Validate() error {
	if len(t) == 0 {
		return ErrInvalidTransport
	}
	return nil
}

type AAGUID []byte

func (a AAGUID) Validate() error {
	if len(a) != 16 && len(a) != 0 {  // Allow 0 for optional, 16 for valid UUID
		return ErrInvalidAAGUID
	}
	return nil
}

type AuthenticatorName string

func (an AuthenticatorName) Validate() error {
	if len(an) == 0 || len(an) > 100 {
		return ErrInvalidAuthenticatorName
	}
	return nil
}
