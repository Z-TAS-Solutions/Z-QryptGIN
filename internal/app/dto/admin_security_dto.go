package dto

// -- Enforce MFA --
type EnforceMfaRequest struct {
	Enabled bool `json:"enabled"`
}

type EnforceMfaResponse struct {
	Message string `json:"message"`
	Data    struct {
		Enabled bool `json:"enabled"`
	} `json:"data"`
}
