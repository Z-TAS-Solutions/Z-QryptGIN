package handlers

import (
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

func getRegistrationOptions() []webauthn.RegistrationOption {
	return []webauthn.RegistrationOption{
		func(cco *protocol.PublicKeyCredentialCreationOptions) {
			// Registration Settings -> Supported Public Key Algorithms: Ed25519, ES256
			cco.Parameters = []protocol.CredentialParameter{
				{Type: protocol.PublicKeyCredentialType, Algorithm: protocol.AlgEdDSA}, // Ed25519
				{Type: protocol.PublicKeyCredentialType, Algorithm: protocol.AlgES256}, // ES256
			}

			// Registration Settings -> Hints: ["security-key", "client-device"]
			cco.Hints = []string{"security-key", "client-device"}
		},
	}
}

// Usage in Handler:
// sessionData, creationData, err := wa.BeginRegistration(user, getRegistrationOptions()...)

func getLoginOptions() []webauthn.LoginOption {
	return []webauthn.LoginOption{
		func(cro *protocol.PublicKeyCredentialRequestOptions) {
			// Authentication Settings -> Hints: ["security-key", "client-device"]
			cro.Hints = []string{"security-key", "client-device"}

			// Authentication Settings -> User Verification: Required
			cro.UserVerification = protocol.VerificationRequired
		},
	}
}

// Usage in Handler:
// sessionData, assertionData, err := wa.BeginLogin(user, getLoginOptions()...)
// Or, if using a username-less flow via Discoverable Credentials:
// sessionData, assertionData, err := wa.BeginDiscoverableLogin(getLoginOptions()...)
