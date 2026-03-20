package database

import {
	"gorm.io/gorm"
	"time"
}

type User struct {
	gorm.Model
	Name string
	CustomID  UserCustomID `gorm:"uniqueIndex;size:10"`
	Email Email `gorm:"uniqueIndex"`
	PhoneNo PhoneNumber `gorm:"uniqueIndex"`
	Nic NIC `gorm:"uniqueIndex"`
	Role UserRole
	PasswordHash string
	Status UserStatus `gorm:"default:Active"`
	SecurityLevel UserSecurityLevel `gorm:"default:Low"`
	PasskeyStatus bool
	//Relationships
	Passkeys []Passkey `gorm:"foreignKey:UserID"`
	Notifications []Notification `gorm:"foreignKey:UserID"`
    MfaChallenges []MfaChallenge `gorm:"foreignKey:UserID"`
    ActivityLogs  []ActivityLog  `gorm:"foreignKey:UserID"`
    Sessions      []Session      `gorm:"foreignKey:UserID"`
}

type Notification struct {
	gorm.Model
	UserID uint
	NotifiID  NotificationID `gorm:"uniqueIndex"`
	Title string
	Message string
	NotifyType NotificationType `gorm:"unique"`
	IsRead bool
}

type MfaChallenge struct {
	gorm.Model
	UserID uint
	MfaID MfaChallengeID `gorm:"uniqueIndex"`
	DeviceName string
	Location string
	IpAddress IPV4
	Status MfaStatus
	Decision MfaDecision
	RespondedAt time.Time
}

type ActivityLog struct {
	gorm.Model
	UserID uint
	ActivityNo ActivityID `gorm:"uniqueIndex"`
	Title string
	Device string
	TimeLabel time.Time
	IsCritical bool
	Type ActivityLogType
}

type Passkey struct {
	gorm.Model
	UserID uint
	PassID PasskeyID `gorm:"uniqueIndex"`
	Name string
	PublicKey string
	BackedUp *string
	Transport string
}

type Session struct {
	gorm.Model
	UserID uint
	SessionNo SessionID `gorm:"uniqueIndex"`
	DeviceName string
	Location string
	IpAddress IPV4
	LastActive time.Time
}