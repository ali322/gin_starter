package lib

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func DecodeJWTToken(tokenStr string, secret string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return token.Claims.(jwt.MapClaims), nil
}

func GenerateJWTToken(secret string, auth map[string]interface{}) (string, error) {
	expired := time.Now().Add(time.Hour * 24 * 30).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"auth": auth, "expire": expired,
	})
	tokenStr, err := token.SignedString([]byte(secret))
	return tokenStr, err
}
