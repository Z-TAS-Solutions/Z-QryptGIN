package service

import (
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

func InitWebAuthn() (*webauthn.WebAuthn, error) {
	config := &webauthn.Config{
		RPDisplayName: "Z-QryptGIN", // Display Name for your site
		RPID:          "api.z-tas.com",
		RPOrigins:     []string{"https://z-tas.com"},

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
