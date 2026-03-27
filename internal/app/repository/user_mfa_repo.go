package repository

import (
	"fmt"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserMfaRepository interface {
	Send(req dto.MfaSendRequest) (*dto.MfaSendResponse, error)
	Respond(req dto.MfaRespondRequest) (*dto.MfaRespondResponse, error)
}

type userMfaRepository struct {
	db *gorm.DB
}

func NewUserMfaRepository(db *gorm.DB) UserMfaRepository {
	return &userMfaRepository{db: db}
}

func (r *userMfaRepository) Send(req dto.MfaSendRequest) (*dto.MfaSendResponse, error) {
	// Usually there's a User ID in this context, assuming we use a generic placeholder or fetch current user.
	// For now, we will create the challenge assigning it to an unknown or dummy userID unless available in req.
	mfaChallenge := database.MfaChallenge{
		MfaID:       database.MfaChallengeID(fmt.Sprintf("Mfa_%s", uuid.New().String()[:8])),
		DeviceName:  req.DeviceID, // Just as a mapped representation
		Status:      "pending",
		RespondedAt: time.Now(),
	}

	err := r.db.Create(&mfaChallenge).Error
	if err != nil {
		return nil, err
	}

	return &dto.MfaSendResponse{Message: "MFA prompt sent successfully"}, nil
}

func (r *userMfaRepository) Respond(req dto.MfaRespondRequest) (*dto.MfaRespondResponse, error) {
	decision := database.MfaDecision(req.Action)
	status := database.MfaStatus("approved")
	if req.Action == "denied" {
		status = database.MfaStatus("denied")
	}

	err := r.db.Model(&database.MfaChallenge{}).
		Where("mfa_id = ?", req.NotificationID).
		Updates(map[string]interface{}{
			"decision":     decision,
			"status":       status,
			"responded_at": time.Now(),
		}).Error

	if err != nil {
		return nil, err
	}

	return &dto.MfaRespondResponse{Message: "MFA action logged successfully"}, nil
}
