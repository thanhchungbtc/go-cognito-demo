package usecase

type Auth interface {
	Login(*LoginInput) (*LoginOutput, error)
	Register(*RegisterInput) (*RegisterOutput, error)
	Confirm(*ConfirmInput) (*ConfirmOutput, error)
	ParseAndVerifyJWT(string) (*Claims, error)
}

type RegisterInput struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
}

type RegisterOutput struct {
	UserSub       *string
	UserConfirmed *bool
}

type LoginInput struct {
	AuthFlow string `json:"auth_flow,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginOutput struct {
	AccessToken *string `json:"access_token,omitempty"`

	// The expiration period of the authentication result in seconds.
	ExpiresIn *int64 `json:"expires_in,omitempty"`

	// The ID token.
	IdToken *string `json:"id_token,omitempty"`

	// The refresh token.
	RefreshToken *string `json:"refresh_token,omitempty"`

	// The token type.
	TokenType *string `json:"token_type,omitempty"`
}

type ConfirmInput struct {
	Username string `json:"username,omitempty"`
	Code     string `json:"code,omitempty"`
}

type ConfirmOutput struct {
	Success bool `json:"success,omitempty"`
}

type Claims struct {
	UserName      string
	Email         string
	Sub           string
	EmailVerified bool
}
