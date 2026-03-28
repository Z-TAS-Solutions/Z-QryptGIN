package repository

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
	"gorm.io/gorm"
)

type WebAuthnCredentialRepository interface {
	CreateCredential(credential *database.WebAuthnCredential) error
	FindCredentialByID(credentialID []byte) (*database.WebAuthnCredential, error)
	FindCredentialsByUserID(userID uint) ([]database.WebAuthnCredential, error)
	UpdateSignCount(credentialID []byte, signCount uint32) error
	UpdateCloneWarning(credentialID []byte, cloneWarning bool) error
}

type webAuthnCredentialRepository struct {
	db *gorm.DB
}

func NewWebAuthnCredentialRepository(db *gorm.DB) WebAuthnCredentialRepository {
	return &webAuthnCredentialRepository{db: db}
}

// CreateCredential inserts a new WebAuthn credential for a user
func (r *webAuthnCredentialRepository) CreateCredential(credential *database.WebAuthnCredential) error {
	return r.db.Create(credential).Error
}

// FindCredentialByID retrieves a credential by its ID
func (r *webAuthnCredentialRepository) FindCredentialByID(credentialID []byte) (*database.WebAuthnCredential, error) {
	var credential database.WebAuthnCredential
	err := r.db.Where("credential_id = ?", credentialID).First(&credential).Error
	return &credential, err
}

// FindCredentialsByUserID retrieves all credentials for a specific user
func (r *webAuthnCredentialRepository) FindCredentialsByUserID(userID uint) ([]database.WebAuthnCredential, error) {
	var credentials []database.WebAuthnCredential
	err := r.db.Where("user_id = ?", userID).Find(&credentials).Error
	return credentials, err
}

// UpdateSignCount updates the signature count for a credential (for cloning detection)
func (r *webAuthnCredentialRepository) UpdateSignCount(credentialID []byte, signCount uint32) error {
	return r.db.Model(&database.WebAuthnCredential{}).
		Where("credential_id = ?", credentialID).
		Update("sign_count", signCount).Error
}

// UpdateCloneWarning sets the clone warning flag for a credential
func (r *webAuthnCredentialRepository) UpdateCloneWarning(credentialID []byte, cloneWarning bool) error {
	return r.db.Model(&database.WebAuthnCredential{}).
		Where("credential_id = ?", credentialID).
		Update("clone_warning", cloneWarning).Error
}
