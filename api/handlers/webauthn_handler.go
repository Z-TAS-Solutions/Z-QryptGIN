package handlers

import (
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/protocol/webauthncose" // Import the COSE algorithms package
	"github.com/go-webauthn/webauthn/webauthn"
)

func getRegistrationOptions() []webauthn.RegistrationOption {
	return []webauthn.RegistrationOption{
		func(cco *protocol.PublicKeyCredentialCreationOptions) {
			// Registration Settings -> Supported Public Key Algorithms: Ed25519, ES256
			cco.Parameters = []protocol.CredentialParameter{
				// Use the webauthncose package for algorithm constants
				{Type: protocol.PublicKeyCredentialType, Algorithm: webauthncose.AlgEdDSA},
				{Type: protocol.PublicKeyCredentialType, Algorithm: webauthncose.AlgES256},
			}

			// Registration Settings -> Hints
			// Explicitly cast strings to the library's required custom type
			cco.Hints = []protocol.PublicKeyCredentialHints{
				protocol.PublicKeyCredentialHints("security-key"),
				protocol.PublicKeyCredentialHints("client-device"),
			}
		},
	}
}

// Usage in Handler:
// sessionData, creationData, err := wa.BeginRegistration(user, getRegistrationOptions()...)

func getLoginOptions() []webauthn.LoginOption {
	return []webauthn.LoginOption{
		func(cro *protocol.PublicKeyCredentialRequestOptions) {
			// Authentication Settings -> Hints
			// Explicitly cast strings to the library's required custom type
			cro.Hints = []protocol.PublicKeyCredentialHints{
				protocol.PublicKeyCredentialHints("security-key"),
				protocol.PublicKeyCredentialHints("client-device"),
			}

			// Authentication Settings -> User Verification: Required
			cro.UserVerification = protocol.VerificationRequired
		},
	}
}

// Usage in Handler:
// sessionData, assertionData, err := wa.BeginLogin(user, getLoginOptions()...)
