package auth

type (
	LoginRequest struct {
		Email    string
		Password string
	}

	AuthResponse struct {
		Token        string      `json:"token"`
		RefreshToken string      `json:"refresh_token"`
		Data         interface{} `json:"data"`
	}

	RenewalTokenRequest struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
)
