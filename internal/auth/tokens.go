package auth

import (
	"crypto/rand"
	"encoding/base64"
	"reflect"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(claims interface{}, expires int, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)

	jwtClaims := token.Claims.(jwt.MapClaims)

	val := reflect.ValueOf(claims)
	for i := 0; i < val.NumField(); i++ {
		fieldName := val.Type().Field(i).Name
		fieldValue := val.Field(i).Interface()
		jwtClaims[fieldName] = fieldValue
	}

	jwtClaims["exp"] = time.Now().Add(time.Duration(expires) * time.Hour).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken() (string, error) {
	byteSlice := make([]byte, 32)
	_, err := rand.Read(byteSlice)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(byteSlice), nil
}
