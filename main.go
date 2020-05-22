package main

import (
	"fmt"
	"log"
	"os"
	"thanhchungbtc/go-cognito-demo/app"
	"thanhchungbtc/go-cognito-demo/domain/usecase"

	"github.com/lestrrat-go/jwx/jwk"
)

func main() {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}
	clientId := os.Getenv("COGNITO_APP_CLIENT_ID")
	userPoolId := os.Getenv("COGNITO_USER_POOL_ID")

	keySet, err := downloadPublickeys(region, userPoolId)
	if err != nil {
		log.Fatal(err)
	}

	cognito := usecase.NewCognito(
		clientId,
		userPoolId)

	auth := usecase.NewAuth(cognito, keySet)

	a := app.New(auth)

	log.Fatal(a.Run(":3000"))

}

func downloadPublickeys(region, userPoolId string) (*jwk.Set, error) {
	url := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", region, userPoolId)
	keySet, err := jwk.FetchHTTP(url)
	if err != nil {
		return nil, err
	}
	return keySet, nil
}
