package cache

type CredentialCache struct {
	UserID          uint     `json:"u_id"`
	CredentialID    []byte   `json:"c_id"`
	PublicKey       []byte   `json:"pk"`
	AttestationType string   `json:"at"`
	Transport       []string `json:"tr"`
	
	// FIDO2 State & Validation
	UserPresent    bool   `json:"up"`
	UserVerified   bool   `json:"uv"`
	BackupEligible bool   `json:"be"`
	BackupState    bool   `json:"bs"`
	
	// Security & Metadata (Crucial for Z-TAS)
	AAGUID       []byte `json:"aa"` // Needed to identify the device type (Yubikey vs Apple)
	SignCount    uint32 `json:"sc"` // MUST be cached to detect cloned keys
	CloneWarning bool   `json:"cw"` // Persistent warning state
	
	// Optional: Metadata for UI
	AuthenticatorName string `json:"n"`
}