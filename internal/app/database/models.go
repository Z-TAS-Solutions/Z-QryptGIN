package database

import {
	"gorm.io/gorm"
	"time"
}

type User struct {
	gorm.Model
	Name string
	CustomID  UserCustomID
	Email Email `gorm:"unique"`
	PhoneNo PhoneNumber `gorm:"unique"`
	Nic NIC `gorm:"unique"`
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
	NotID 
	Title string
	Message string
	NotifyType NotificationType
	IsRead bool
}

type MfaChallenge struct {
	gorm.Model
	UserID uint
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
	Title string
	Device string
	TimeLabel time.Time
	IsCritical bool
	Type ActivityLogType
}

type Passkey struct {
	gorm.Model
	UserID uint
	Name string
	PublicKey string
	BackedUp *string
	Transport string
}

type Session struct {
	gorm.Model
	UserID uint
	DeviceName string
	Location string
	IpAddress IPV4
	LastActive time.Time
}