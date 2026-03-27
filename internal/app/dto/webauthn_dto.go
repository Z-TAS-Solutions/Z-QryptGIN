package dto

import "github.com/go-webauthn/webauthn/protocol"

// Endpoint 1: Fetch Registration Args (/register/begin)
type RegisterBeginRequest struct {
	Username string `json:"username" binding:"required"`
}

type RegisterBeginResponse struct {
	// A temporary token to map the upcoming finish request to the session data
	SessionToken string                       `json:"session_token"`
	CreationData *protocol.CredentialCreation `json:"creation_data"`
}

// Endpoint 2: Send Passkey Data (/register/finish)
type RegisterFinishRequest struct {
	SessionToken string `header:"X-Session-Token" binding:"required"`
	// Note: The HTTP body is the raw JSON from the WebAuthn browser API.
	// You will parse it using `protocol.ParseCredentialCreationResponseponseBody(r.Body)`
}

// Endpoint 3: Fetch Authentication Args (/login/begin)
type LoginBeginRequest struct {
	// Optional: Since you require Discoverable Credentials (Responseident Keys),
	// users can technically authenticate without providing a username upfront (User-nameless flow).
	Username string `json:"username,omitempty"`
}

type LoginBeginResponse struct {
	SessionToken  string                        `json:"session_token"`
	AssertionData *protocol.CredentialAssertion `json:"assertion_data"`
}

// Endpoint 4: Send Signed Passkey Data (/login/finish)
type LoginFinishRequest struct {
	SessionToken string `header:"X-Session-Token" binding:"required"`
	// Note: The HTTP body is the raw JSON from the WebAuthn browser API.
	// You will parse it using `protocol.ParseCredentialRequestuestResponseponseBody(r.Body)`
}
