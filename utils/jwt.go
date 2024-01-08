package utils

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserId string
	Role   string
}

func VerifyToken(token string) (*TokenClaims, error) {

	SECRET_KEY := os.Getenv("SECRET_KEY")

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Unexpected signin method")
		}
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		return nil, errors.New("Could not parse token")
	}

	tokenIsValid := parsedToken.Valid

	if !tokenIsValid {
		return nil, errors.New("Invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return nil, errors.New("Invalid token claims")
	}

	var tokenClaims TokenClaims

	tokenClaims.UserId = claims["id"].(string)
	tokenClaims.Role = claims["role"].(string)

	return &tokenClaims, nil
}
