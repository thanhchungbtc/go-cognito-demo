package usecase

import (
	"crypto/rsa"
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
)

type auth struct {
	cognito *Cognito
	keys    *jwk.Set
}

func (a *auth) Login(input *LoginInput) (*LoginOutput, error) {
	var params map[string]*string
	if input.AuthFlow == "REFRESH_TOKEN_AUTH" {
		params = map[string]*string{
			"REFRESH_TOKEN": aws.String(input.RefreshToken),
		}
	} else if input.AuthFlow == "USER_PASSWORD_AUTH" {
		params = map[string]*string{
			"USERNAME": aws.String(input.Username),
			"PASSWORD": aws.String(input.Password),
		}
	} else {
		return nil, errors.New("unsupport auth flow")
	}

	output, err := a.cognito.Client.InitiateAuth(&cognitoidentityprovider.InitiateAuthInput{
		ClientId:       &a.cognito.AppClientID,
		AuthFlow:       aws.String(input.AuthFlow),
		AuthParameters: params,
	})

	if err != nil {
		return nil, err
	}

	authResult := output.AuthenticationResult
	return &LoginOutput{
		AccessToken:  authResult.RefreshToken,
		TokenType:    authResult.TokenType,
		IdToken:      authResult.IdToken,
		RefreshToken: authResult.RefreshToken,
		ExpiresIn:    authResult.ExpiresIn,
	}, nil
}

func (a *auth) Register(input *RegisterInput) (*RegisterOutput, error) {
	output, err := a.cognito.Client.SignUp(&cognitoidentityprovider.SignUpInput{
		Username: &input.Username,
		Password: &input.Password,
		ClientId: &a.cognito.AppClientID,
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: &input.Email,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &RegisterOutput{
		UserSub:       output.UserSub,
		UserConfirmed: output.UserConfirmed,
	}, nil
}

func (a *auth) Confirm(input *ConfirmInput) (*ConfirmOutput, error) {
	_, err := a.cognito.Client.ConfirmSignUp(&cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         &a.cognito.AppClientID,
		ConfirmationCode: aws.String(input.Code),
		Username:         &input.Username,
	})
	return &ConfirmOutput{Success: err == nil}, err
}

func (a *auth) ParseAndVerifyJWT(tokenStr string) (*Claims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		kid := token.Header["kid"].(string)
		jsonWebKey := a.keys.LookupKeyID(kid)[0]

		publicKey := &rsa.PublicKey{}

		if err := jsonWebKey.Raw(publicKey); err != nil {
			return nil, err
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	// the token has been verified, now we can safely use the unverfied claimes
	claims := token.Claims.(jwt.MapClaims)

	// https://docs.aws.amazon.com/cognito/latest/developerguide/amazon-cognito-user-pools-using-tokens-verifying-a-jwt.html
	// Verify that the token is not expired.
	// The audience (aud) claim should match the app client ID created in the Amazon Cognito user pool.
	// The issuer (iss) claim should match your user pool. For example, a user pool created in the us-east-1 region will have an iss value of:
	// https://cognito-idp.us-east-1.amazonaws.com/<userpoolID>.
	// Check the token_use claim.
	// If you are only accepting the access token in your web APIs, its value must be access.
	// If you are only using the ID token, its value must be id.
	// If you are using both ID and access tokens, the token_use claim must be either id or access.

	// exp
	exp := claims["exp"].(float64)
	if float64(time.Now().Unix()) > exp {
		return nil, errors.New("Token expired")
	}

	// aud must match the client_id
	if claims["aud"] != a.cognito.AppClientID {
		return nil, errors.New("Token was not issued for this audience")
	}

	// all good
	var results Claims
	results.UserName = claims["cognito:username"].(string)
	results.Email = claims["email"].(string)
	results.Sub = claims["sub"].(string)
	results.EmailVerified = claims["email_verified"].(bool)

	return &results, nil
}

type Cognito struct {
	Client      *cognitoidentityprovider.CognitoIdentityProvider
	AppClientID string
	UserPoolID  string
}

func NewCognito(appClientID, userPoolID string) *Cognito {
	conf := &aws.Config{
		Region: aws.String("us-east-1"),
	}

	sess, err := session.NewSession(conf)
	if err != nil {
		log.Fatal(err)
	}
	client := cognitoidentityprovider.New(sess)
	return &Cognito{
		Client:      client,
		AppClientID: appClientID,
		UserPoolID:  userPoolID,
	}
}
