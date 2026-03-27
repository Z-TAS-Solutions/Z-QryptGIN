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
	PasswordHash  string
	Status        UserStatus        `gorm:"default:Active"`
	SecurityLevel UserSecurityLevel `gorm:"default:Low"`
	MFAStatus     bool
	//Relationships
	Passkeys      []Passkey            `gorm:"foreignKey:UserID"`
	Notifications []Notification       `gorm:"foreignKey:UserID"`
	MfaChallenges []MfaChallenge       `gorm:"foreignKey:UserID"`
	ActivityLogs  []ActivityLog        `gorm:"foreignKey:UserID"`
	Sessions      []Session            `gorm:"foreignKey:UserID"`
	Credentials   []WebAuthnCredential `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
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
	UserID uint
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

type Passkey struct {
	gorm.Model
	UserID    uint
	PassID    PasskeyID `gorm:"uniqueIndex"`
	Name      string
	PublicKey string
	BackedUp  *string
	Transport string
}

type Session struct {
	gorm.Model
	UserID     uint
	SessionNo  SessionID `gorm:"uniqueIndex"`
	DeviceName string
	Location   string
	IpAddress  IPV4
	LastActive time.Time
}

type WebAuthnCredential struct {
	gorm.Model
	UserID uint `gorm:"index;not null"`

	// FIDO2 / WebAuthn Core Data
	CredentialID    []byte `gorm:"type:bytea;uniqueIndex;not null"` // The 'kid'
	PublicKey       []byte `gorm:"type:bytea;not null"`
	AttestationType string `gorm:"size:64"`

	// Transport is stored as JSON (e.g., ["usb", "ble", "nfc", "internal"])
	// This matches your preference for "Security-key" vs "client-device"
	Transport []byte `gorm:"type:jsonb"`

	// FIDO2 Flags & Backup Status
	UserPresent    bool
	UserVerified   bool
	BackupEligible bool // Can this passkey be synced?
	BackupState    bool // Is it currently synced/backed up?

	// Authenticator Metadata
	AAGUID       []byte `gorm:"type:bytea"` // Identifies the device model
	SignCount    uint32 `gorm:"not null;default:0"`
	CloneWarning bool   `gorm:"default:false"`

	// Friendly Name (consistent with your Passkey struct)
	AuthenticatorName string `gorm:"size:100"`
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

func (p *Passkey) BeforeSave(tx *gorm.DB) error {
	return p.PassID.Validate()
}

func (s *Session) BeforeSave(tx *gorm.DB) error {
	if err := s.SessionNo.Validate(); err != nil {
		return err
	}
	if err := s.IpAddress.Validate(); err != nil {
		return err
	}
	return nil
}

// WebAuthn.User Interface Implementations
func (u *User) WebAuthnID() []byte          { return []byte(u.ID) }
func (u *User) WebAuthnName() string        { return u.Username }
func (u *User) WebAuthnDisplayName() string { return u.Username }
func (u *User) WebAuthnCredentials() []webauthn.Credential {
	var creds []webauthn.Credential
	for _, c := range u.Credentials {
		// Deserialize the stored JSON transports
		var transports []protocol.AuthenticatorTransport
		_ = json.Unmarshal(c.Transports, &transports)

		creds = append(creds, webauthn.Credential{
			ID:              c.CredentialID,
			PublicKey:       c.PublicKey,
			AttestationType: c.AttestationType,
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
