package database

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name          string
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
	Passkeys      []Passkey      `gorm:"foreignKey:UserID"`
	Notifications []Notification `gorm:"foreignKey:UserID"`
	MfaChallenges []MfaChallenge `gorm:"foreignKey:UserID"`
	ActivityLogs  []ActivityLog  `gorm:"foreignKey:UserID"`
	Sessions      []Session      `gorm:"foreignKey:UserID"`
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

//-- GORM Hooks to early fail validation logic--

// BeforeSave hooks for automatic validation in GORM

func (u *User) BeforeSave(tx *gorm.DB) error {
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
