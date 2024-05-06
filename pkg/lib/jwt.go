package lib

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/the_fourth_dimension/planet_registry/pkg/env"
)

func SignJwt(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(env.GetEnv(env.JWT_SECRET)))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyJwt(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(env.GetEnv(env.JWT_SECRET)), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)

	if !token.Valid || !ok {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
