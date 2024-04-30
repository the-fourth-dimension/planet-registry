package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/the_fourth_dimension/planet_registry/pkg/env"
)

func VerifyJwt(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(env.GetEnv(env.JWT_SECRET)), nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}
	return nil
}
