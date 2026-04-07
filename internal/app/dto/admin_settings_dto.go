package dto

// Admin Settings DTOs
type SecurityPolicy struct {
	MFARequired              bool `json:"mfa_required"`
	WebAuthnRequired         bool `json:"webauthn_required"`
	PasswordExpiryDays       int  `json:"password_expiry_days"`
	MaxLoginAttempts         int  `json:"max_login_attempts"`
	LockoutDurationMinutes   int  `json:"lockout_duration_minutes"`
	SessionTimeoutMinutes    int  `json:"session_timeout_minutes"`
	AllowWeakPasswords       bool `json:"allow_weak_passwords"`
	require2FAForAdmin       bool `json:"require_2fa_for_admin"`
}

type NotificationSettings struct {
	EmailOnNewLogin        bool `json:"email_on_new_login"`
	EmailOnSecurityChange  bool `json:"email_on_security_change"`
	EmailOnAccountActivity bool `json:"email_on_account_activity"`
	SMSNotifications       bool `json:"sms_notifications"`
}

type AdminSettingsResponse struct {
	Success              bool                  `json:"success"`
	SecurityPolicy       SecurityPolicy        `json:"security_policy"`
	NotificationSettings NotificationSettings  `json:"notification_settings"`
	AppVersion           string                `json:"app_version"`
	LastUpdated          string                `json:"last_updated"`
}

type UpdateSettingsRequest struct {
	SecurityPolicy       *SecurityPolicy       `json:"security_policy,omitempty"`
	NotificationSettings *NotificationSettings `json:"notification_settings,omitempty"`
}

type UpdateSettingsResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type AdminMFAEnforcementRequest struct {
	Enabled bool   `json:"enabled"`
	Reason  string `json:"reason,omitempty"`
}

type AdminMFAEnforcementResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Status  bool   `json:"status"`
}
