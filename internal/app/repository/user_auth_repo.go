package repository

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"gorm.io/gorm"
)

type UserAuthRepository interface {
	RegisterOptions(req dto.PasskeyRegisterOptionsRequest) (*dto.PasskeyRegisterOptionsResponse, error)
	RegisterVerify(req dto.PasskeyRegisterVerifyRequest) (*dto.PasskeyRegisterVerifyResponse, error)
	LoginOptions(req dto.PasskeyLoginOptionsRequest) (*dto.PasskeyLoginOptionsResponse, error)
	LoginVerify(req dto.PasskeyLoginVerifyRequest) (*dto.PasskeyLoginVerifyResponse, error)
}

type userAuthRepository struct {
	db *gorm.DB
}

func NewUserAuthRepository(db *gorm.DB) UserAuthRepository {
	return &userAuthRepository{db: db}
}

func (r *userAuthRepository) RegisterOptions(req dto.PasskeyRegisterOptionsRequest) (*dto.PasskeyRegisterOptionsResponse, error) {
	var res dto.PasskeyRegisterOptionsResponse
	res.Message = "Passkey registration options generated"
	res.Data.Challenge = "base64url_string_mock"
	res.Data.RP.Name = "Z-TAS"
	res.Data.RP.ID = "z-tas.com"
	res.Data.User.ID = "mock_user_id"
	res.Data.User.Name = req.Email
	res.Data.User.DisplayName = req.Email
	res.Data.PubKeyCredParams = []map[string]interface{}{}
	res.Data.Timeout = 60000
	return &res, nil
}

func (r *userAuthRepository) RegisterVerify(req dto.PasskeyRegisterVerifyRequest) (*dto.PasskeyRegisterVerifyResponse, error) {
	var res dto.PasskeyRegisterVerifyResponse
	res.Message = "Registration successful"
	res.Data.UserID = "user_123"
	res.Data.AccessToken = "eyJ..."
	return &res, nil
}

func (r *userAuthRepository) LoginOptions(req dto.PasskeyLoginOptionsRequest) (*dto.PasskeyLoginOptionsResponse, error) {
	var res dto.PasskeyLoginOptionsResponse
	res.Message = "Passkey login options generated"
	res.Data.Challenge = "base64url_string_mock"
	res.Data.RPID = "z-tas.com"
	res.Data.Timeout = 60000
	return &res, nil
}

func (r *userAuthRepository) LoginVerify(req dto.PasskeyLoginVerifyRequest) (*dto.PasskeyLoginVerifyResponse, error) {
	var res dto.PasskeyLoginVerifyResponse
	res.Message = "Login successful"
	res.Data.UserID = "user_123"
	res.Data.AccessToken = "eyJ..."
	res.Data.RefreshToken = "def502..."
	return &res, nil
}
