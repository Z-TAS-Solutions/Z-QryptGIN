package dto

// -- Register Options --
type PasskeyRegisterOptionsRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type PasskeyRegisterOptionsResponse struct {
	Message string `json:"message"`
	Data    struct {
		Challenge string `json:"challenge"`
		RP        struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"rp"`
		User struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
		} `json:"user"`
		PubKeyCredParams []map[string]interface{} `json:"pubKeyCredParams"`
		Timeout          int                      `json:"timeout"`
	} `json:"data"`
}

// -- Register Verify --
type PasskeyRegisterVerifyRequest struct {
	ID       string `json:"id" binding:"required"`
	RawID    string `json:"rawId" binding:"required"`
	Response struct {
		ClientDataJSON    string `json:"clientDataJSON" binding:"required"`
		AttestationObject string `json:"attestationObject" binding:"required"`
	} `json:"response" binding:"required"`
	Type string `json:"type" binding:"required"`
}

type PasskeyRegisterVerifyResponse struct {
	Message string `json:"message"`
	Data    struct {
		UserID      string `json:"userId"`
		AccessToken string `json:"accessToken"`
	} `json:"data"`
}

// -- Login Options --
type PasskeyLoginOptionsRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type PasskeyLoginOptionsResponse struct {
	Message string `json:"message"`
	Data    struct {
		Challenge string `json:"challenge"`
		RPID      string `json:"rpId"`
		Timeout   int    `json:"timeout"`
	} `json:"data"`
}

// -- Login Verify --
type PasskeyLoginVerifyRequest struct {
	ID       string `json:"id" binding:"required"`
	RawID    string `json:"rawId" binding:"required"`
	Response struct {
		ClientDataJSON    string `json:"clientDataJSON" binding:"required"`
		AuthenticatorData string `json:"authenticatorData" binding:"required"`
		Signature         string `json:"signature" binding:"required"`
		UserHandle        string `json:"userHandle"`
	} `json:"response" binding:"required"`
	Type string `json:"type" binding:"required"`
}

type PasskeyLoginVerifyResponse struct {
	Message string `json:"message"`
	Data    struct {
		UserID       string `json:"userId"`
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	} `json:"data"`
}
