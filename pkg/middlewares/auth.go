package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/the_fourth_dimension/planet_registry/pkg/errors/HttpError"
	"github.com/the_fourth_dimension/planet_registry/pkg/lib/jwt"
	"github.com/the_fourth_dimension/planet_registry/pkg/roles"
)

func AuthMiddleware() func(*gin.Context) {
	return func(ctx *gin.Context) {
		tokenString := ctx.Request.Header.Get("Authorization")
		if tokenString == "" {
			ctx.Error(HttpError.NewHttpError("missing header", "Authorization", http.StatusBadRequest))
			ctx.Abort()
			return
		}
		if len(tokenString) < 8 {
			ctx.Error(HttpError.NewHttpError("Invalid authorization header", tokenString, http.StatusForbidden))
			ctx.Abort()
			return
		}
		tokenString = tokenString[len("Bearer "):]
		claims, err := jwt.VerifyJwt(tokenString)
		if err != nil {
			ctx.Error(HttpError.NewHttpError("Invalid jwt token", tokenString, http.StatusForbidden))
			ctx.Abort()
			return
		}
		role, ok := claims["role"].(string)
		if !ok {
			ctx.Error(HttpError.NewHttpError("missing claim", "role", http.StatusBadRequest))
			ctx.Abort()
			return
		}
		if !roles.IsValidRole(role) {
			ctx.Error(HttpError.NewHttpError("invalid role", role, http.StatusUnauthorized))
			ctx.Abort()
			return
		}
		ctx.Set("tokenClaims", claims)
		ctx.Next()
	}
}
