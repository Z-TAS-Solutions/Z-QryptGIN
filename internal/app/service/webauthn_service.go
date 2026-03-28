package service

import (
	"context"
	"encoding/json"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/protocol/webauthncose"
	"github.com/go-webauthn/webauthn/webauthn"
)

func InitWebAuthn() (*webauthn.WebAuthn, error) {
	config := &webauthn.Config{
		RPDisplayName: "Z-QryptGIN", // Display Name for your site
		// RPID:          "api.z-tas.com",
		RPID:          "localhost",
		// RPOrigins:     []string{"https://z-tas.com"},
		RPOrigins:     []string{"http://localhost:5173"},

		// Registration Settings -> Attestation: Direct
		AttestationPreference: protocol.PreferDirectAttestation,

		AuthenticatorSelection: protocol.AuthenticatorSelection{
			// Registration/Auth Settings -> User Verification: Required
			UserVerification: protocol.VerificationRequired,

			// Registration Settings -> Discoverable Credential: Required
			ResidentKey:        protocol.ResidentKeyRequirementRequired,
			RequireResidentKey: protocol.ResidentKeyRequired(),

			// Registration Settings -> Attachment: All Supported
			// Leaving this empty permits both platform (TouchID) and cross-platform (Yubikey)
			AuthenticatorAttachment: protocol.AuthenticatorAttachment(""),
		},
	}
	return webauthn.New(config)
}

// WebAuthnService handles WebAuthn registration and authentication operations
type WebAuthnService struct {
	wa *webauthn.WebAuthn
}

// NewWebAuthnService creates a new WebAuthnService instance
func NewWebAuthnService(wa *webauthn.WebAuthn) *WebAuthnService {
	return &WebAuthnService{wa: wa}
}

// getRegistrationOptions returns the array of registration options for BeginRegistration
func (s *WebAuthnService) getRegistrationOptions() []webauthn.RegistrationOption {
	return []webauthn.RegistrationOption{
		func(cco *protocol.PublicKeyCredentialCreationOptions) {
			// Registration Settings -> Supported Public Key Algorithms: Ed25519, ES256
			cco.Parameters = []protocol.CredentialParameter{
				{Type: protocol.PublicKeyCredentialType, Algorithm: webauthncose.AlgEdDSA},
				{Type: protocol.PublicKeyCredentialType, Algorithm: webauthncose.AlgES256},
			}

			// Registration Settings -> Hints
			cco.Hints = []protocol.PublicKeyCredentialHints{
				protocol.PublicKeyCredentialHints("security-key"),
				protocol.PublicKeyCredentialHints("client-device"),
			}
		},
	}
}

// BeginRegistration initiates the WebAuthn registration ceremony for a user
// Accepts any webauthn.User (database.User, PendingUser, etc.)
func (s *WebAuthnService) BeginRegistration(ctx context.Context, user webauthn.User) (*webauthn.SessionData, *protocol.CredentialCreation, error) {
	// Begin the registration ceremony
	creationData, sessionData, err := s.wa.BeginRegistration(user, s.getRegistrationOptions()...)
	if err != nil {
		return nil, nil, err
	}

	// Serialize session data for caching
	_ = sessionData // Keep reference for cache storage

	return sessionData, creationData, nil
}

// SerializeSessionData converts WebAuthn session data to JSON for caching
func SerializeSessionData(sessionData *webauthn.SessionData) ([]byte, error) {
	return json.Marshal(sessionData)
}

// DeserializeSessionData converts JSON back to WebAuthn session data
func DeserializeSessionData(data []byte) (*webauthn.SessionData, error) {
	var sessionData webauthn.SessionData
	err := json.Unmarshal(data, &sessionData)
	if err != nil {
		return nil, err
	}
	return &sessionData, nil
}

// FinishRegistration verifies the credential creation response and returns the credential
// This is called after the user completes the passkey registration ceremony
// The credentialResponse should be parsed via protocol.ParseCredentialCreationResponseBody
func (s *WebAuthnService) FinishRegistration(ctx context.Context, user webauthn.User, parsedResponse *protocol.ParsedCredentialCreationData, sessionData *webauthn.SessionData) (*webauthn.Credential, error) {
	// Verify the parsed credential response against the session data
	// Note: CreateCredential expects sessionData as a value, not a pointer
	credential, err := s.wa.CreateCredential(user, *sessionData, parsedResponse)
	if err != nil {
		return nil, err
	}
	return credential, nil
}
