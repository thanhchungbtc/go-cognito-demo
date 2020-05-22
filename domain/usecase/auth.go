package usecase

import (
	"github.com/lestrrat-go/jwx/jwk"
)

type Auth interface {
	Login(*LoginInput) (*LoginOutput, error)
	Register(*RegisterInput) (*RegisterOutput, error)
	Confirm(*ConfirmInput) (*ConfirmOutput, error)
	ParseAndVerifyJWT(string) (*Claims, error)
}

func NewAuth(cognito *Cognito, keys *jwk.Set) Auth {
	return &auth{cognito: cognito, keys: keys}
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
	AuthFlow     string `json:"auth_flow,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
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
	UserName      string `json:"user_name,omitempty"`
	Email         string `json:"email,omitempty"`
	Sub           string `json:"sub,omitempty"`
	EmailVerified bool   `json:"email_verified,omitempty"`
}

type LogoutInput struct {
	AccessToken string `json:"access_token,omitempty"`
}
type LogoutOutput struct {
}
