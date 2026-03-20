package database

import {
	"gorm.io/gorm"
	"time"
}

type User struct {
	gorm.Model
	Name string
	Email Email `gorm:"unique"`
	PhoneNo PhoneNumber `gorm:"unique"`
	Role UserRole
	PasswordHash string
	Status UserStatus `gorm:"default:Active"`
	SecurityLevel UserSecurityLevel `gorm:"default:Low"`
	PasskeyStatus bool
}

type Passkey struct {
	gorm.Model

}

type Notification struct {
	gorm.Model
	NotID 
	Title string
	Message string
	NotifyType NotificationType
	IsRead bool
}

type MfaChallenges struct {
	gorm.Model
	DeviceName string
	Location string
	IpAddress IPV4
	Status MfaStatus
	Decision MfaDecision
	RespondedAt time.Time
}

type ActivityLog struct {
	gorm.Model
	Title string
	Device string
	TimeLabel time.Time
	IsCritical bool
	Type ActivityLogType
}

type Passkey struct {
	gorm.Model
	Name string
	PublicKey string
	BackedUp *string
	Transport string
}

type Session struct {
	gorm.Model
	DeviceName string
	Location string
	IpAddress IPV4
	LastActive time.Time
}