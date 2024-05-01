package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/the_fourth_dimension/planet_registry/pkg/errors/HttpError"
	"github.com/the_fourth_dimension/planet_registry/pkg/lib/jwt"
)

func AuthMiddleware() func(*gin.Context) {
	return func(ctx *gin.Context) {
		tokenString := ctx.Request.Header.Get("Authorization")
		if tokenString == "" {
			ctx.Error(HttpError.NewHttpError("missing header", "Authorization", http.StatusBadRequest))
			return
		}
		if len(tokenString) < 8 {
			ctx.Error(HttpError.NewHttpError("Invalid authorization header", tokenString, http.StatusForbidden))
			return
		}
		tokenString = tokenString[len("Bearer "):]
		claims, err := jwt.VerifyJwt(tokenString)
		if err != nil {
			ctx.Error(HttpError.NewHttpError("Invalid jwt token", tokenString, http.StatusForbidden))
			return
		}
		ctx.Set("tokenClaims", claims)
		ctx.Next()
	}
}
