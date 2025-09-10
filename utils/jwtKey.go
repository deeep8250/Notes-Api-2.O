package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var Secret_key = []byte("my_secret_key")

func Jwt(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(Secret_key)

}
