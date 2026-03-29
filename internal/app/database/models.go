package database

import (
	"encoding/json"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name          Name
	CustomID      UserCustomID `gorm:"uniqueIndex;size:10"`
	Email         Email        `gorm:"uniqueIndex"`
	PhoneNo       PhoneNumber  `gorm:"uniqueIndex"`
	Nic           NIC          `gorm:"uniqueIndex"`
	Role          UserRole
	Status        UserStatus        `gorm:"default:Active"`
	SecurityLevel UserSecurityLevel `gorm:"default:Low"`
	MFAStatus     bool
	//Relationships
	Notifications  []Notification       `gorm:"foreignKey:UserID"`
	MfaChallenges  []MfaChallenge       `gorm:"foreignKey:UserID"`
	ActivityLogs   []ActivityLog        `gorm:"foreignKey:UserID"`
	Sessions       []Session            `gorm:"foreignKey:UserID"`
	Credentials    []WebAuthnCredential `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CrypticRecords []CrypticRecord      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Notification struct {
	gorm.Model
	UserID     uint
	NotifiID   NotificationID `gorm:"uniqueIndex"`
	Title      string
	Message    string
	NotifyType NotificationType `gorm:"unique"`
	IsRead     bool
}

type CrypticRecord struct {
	gorm.Model
	UserID uint `gorm:"index;not null"`

	SchemaVersion uint16     `gorm:"not null"`
	TemplateID    TemplateID `gorm:"size:100;index"`
	TemplateType  string     `gorm:"size:50"`
	TemplateVer   uint16     `gorm:"not null"`

	TemplateNonce Nonce12     `gorm:"type:bytea;not null"`
	WrappedDek    CrypticData `gorm:"type:bytea;not null"`
	WrapNonce     Nonce12     `gorm:"type:bytea;not null"`
	Ciphertext    CrypticData `gorm:"type:bytea;not null"`
}

type MfaChallenge struct {
	gorm.Model
	UserID      uint
	MfaID       MfaChallengeID `gorm:"uniqueIndex"`
	DeviceName  string
	Location    string
	IpAddress   IPV4
	Status      MfaStatus
	Decision    MfaDecision
	RespondedAt time.Time
}

type ActivityLog struct {
	gorm.Model
	UserID     uint
	ActivityNo ActivityID `gorm:"uniqueIndex"`
	Title      string
	Device     string
	TimeLabel  time.Time
	IsCritical bool
	Type       ActivityLogType
}

type Session struct {
	gorm.Model
	UserID     uint
	SessionNo  SessionID `gorm:"uniqueIndex"`
	DeviceID   DeviceCustomID `gorm:"index"`
	DeviceName string
	Location   string
	IpAddress  IPV4
	LastActive time.Time
}

type WebAuthnCredential struct {
	gorm.Model
	UserID uint `gorm:"index;not null"`

	// FIDO2 / WebAuthn Core Data
	CredentialID    CredentialID      `gorm:"type:bytea;uniqueIndex;not null"` // The 'kid'
	PublicKey       WebAuthnPublicKey `gorm:"type:bytea;not null"`
	AttestationType AttestationType   `gorm:"size:64"`

	// Transport is stored as JSON (e.g., ["usb", "ble", "nfc", "internal"])
	// This matches your preference for "Security-key" vs "client-device"
	Transport WebAuthnTransport `gorm:"type:jsonb"`

	// FIDO2 Flags & Backup Status
	UserPresent    bool
	UserVerified   bool
	BackupEligible bool // Can this passkey be synced?
	BackupState    bool // Is it currently synced/backed up?

	// Authenticator Metadata
	AAGUID       AAGUID `gorm:"type:bytea"` // Identifies the device model
	SignCount    uint32 `gorm:"not null;default:0"`
	CloneWarning bool   `gorm:"default:false"`

	// Friendly Name (consistent with your Passkey struct)
	AuthenticatorName AuthenticatorName `gorm:"size:100"`
}

//-- GORM Hooks to early fail validation logic--

// BeforeSave hooks for automatic validation in GORM

func (u *User) BeforeSave(tx *gorm.DB) error {
	if err := u.Name.Validate(); err != nil {
		return err
	}
	if err := u.CustomID.Validate(); err != nil {
		return err
	}
	if err := u.Email.Validate(); err != nil {
		return err
	}
	if err := u.PhoneNo.Validate(); err != nil {
		return err
	}
	if err := u.Nic.Validate(); err != nil {
		return err
	}
	if err := u.Role.Validate(); err != nil {
		return err
	}
	if err := u.Status.Validate(); err != nil {
		return err
	}
	if err := u.SecurityLevel.Validate(); err != nil {
		return err
	}
	return nil
}

func (n *Notification) BeforeSave(tx *gorm.DB) error {
	if err := n.NotifiID.Validate(); err != nil {
		return err
	}
	if err := n.NotifyType.Validate(); err != nil {
		return err
	}
	return nil
}

func (m *MfaChallenge) BeforeSave(tx *gorm.DB) error {
	if err := m.MfaID.Validate(); err != nil {
		return err
	}
	if err := m.IpAddress.Validate(); err != nil {
		return err
	}
	if err := m.Status.Validate(); err != nil {
		return err
	}
	return nil
}

func (a *ActivityLog) BeforeSave(tx *gorm.DB) error {
	if err := a.ActivityNo.Validate(); err != nil {
		return err
	}
	if err := a.Type.Validate(); err != nil {
		return err
	}
	return nil
}

func (s *Session) BeforeSave(tx *gorm.DB) error {
	if err := s.SessionNo.Validate(); err != nil {
		return err
	}
	if err := s.DeviceID.Validate(); err != nil {
		return err
	}
	if err := s.IpAddress.Validate(); err != nil {
		return err
	}
	return nil
}

func (cr *CrypticRecord) BeforeSave(tx *gorm.DB) error {
	if err := cr.TemplateID.Validate(); err != nil {
		return err
	}
	if err := cr.TemplateNonce.Validate(); err != nil {
		return err
	}
	if err := cr.WrappedDek.Validate(); err != nil {
		return err
	}
	if err := cr.WrapNonce.Validate(); err != nil {
		return err
	}
	if err := cr.Ciphertext.Validate(); err != nil {
		return err
	}
	return nil
}

func (wc *WebAuthnCredential) BeforeSave(tx *gorm.DB) error {
	if err := wc.CredentialID.Validate(); err != nil {
		return err
	}
	if err := wc.PublicKey.Validate(); err != nil {
		return err
	}
	if err := wc.AttestationType.Validate(); err != nil {
		return err
	}
	if err := wc.Transport.Validate(); err != nil {
		return err
	}
	if err := wc.AAGUID.Validate(); err != nil {
		return err
	}
	if err := wc.AuthenticatorName.Validate(); err != nil {
		return err
	}
	return nil
}

func (u *User) WebAuthnID() []byte {
	// WebAuthn requires a stable, unique []byte.
	// We cast your UserCustomID to a string, then to a byte array.
	return []byte(string(u.CustomID))
}

func (u *User) WebAuthnName() string {
	// The unique identifier the browser uses to distinguish accounts.
	// Casting your custom Email type to a standard string.
	return string(u.Email)
}

func (u *User) WebAuthnDisplayName() string {
	// The friendly name displayed on the user's security key or phone.
	// Casting your custom Name type to a standard string.
	return string(u.Name)
}

func (u *User) WebAuthnCredentials() []webauthn.Credential {
	var creds []webauthn.Credential
	for _, c := range u.Credentials {
		// Deserialize the stored JSON transports
		var transports []protocol.AuthenticatorTransport
		_ = json.Unmarshal(c.Transport, &transports)

		creds = append(creds, webauthn.Credential{
			ID:              c.CredentialID,
			PublicKey:       c.PublicKey,
			AttestationType: string(c.AttestationType),
			Transport:       transports,
			Flags: webauthn.CredentialFlags{
				UserPresent:    c.UserPresent,
				UserVerified:   c.UserVerified,
				BackupEligible: c.BackupEligible,
				BackupState:    c.BackupState,
			},
			Authenticator: webauthn.Authenticator{
				AAGUID:       c.AAGUID,
				SignCount:    c.SignCount,
				CloneWarning: c.CloneWarning,
			},
		})
	}
	return creds
}
