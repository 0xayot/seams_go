package utils

import (
	"context"
	"fmt"
	"os"
	"seams_go/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
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

type Claims struct {
	Id string `json:"user_id"`
	jwt.StandardClaims
}

func GenerateJWT(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Id: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30 * 24 * time.Hour).Unix(),
			Issuer:    "seams",
		},
	})

	secretKey := []byte(os.Getenv("SECRET_KEY"))

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Printf("Error signing token: %v\n", err)
		return "", fmt.Errorf("error generating jwt: %w", err)
	}
	return tokenString, nil
}

func EnsureAuthurised(jwtToken string) (*models.User, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	token, err := jwt.ParseWithClaims(jwtToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		fmt.Printf("Claims: %+v\n", claims)
		var user models.User

		result := DB.Where("id = ?", claims.Id).First(&user)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				return nil, fmt.Errorf("User not found")
			} else {
				return nil, fmt.Errorf("Error occurred: %v\n", result.Error)
			}

		}
		return &user, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
