package database

import {
	"gorm.io/gorm"
}

type User struct {
	gorm.Model
	Name string
	Email Email `gorm:"unique"`
	PhoneNo PhoneNumber `gorm:"unique"`
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
	title string
	message string
	NotifyType NotificationType
	IsRead bool
}