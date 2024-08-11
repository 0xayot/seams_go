package utils

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/api/idtoken"
)

func ValidateGoogleIdToken(email string, id_token string) bool {
	audience := os.Getenv("GOOGLE_CLIENT_ID")

	// if !audience || !id_token{

	// }

	payload, err := idtoken.Validate(context.Background(), id_token, audience)
	if err != nil {
		panic(err)
	}
	fmt.Print(payload.Claims)
	gmail := payload.Claims["email"]
	return email == gmail
}
