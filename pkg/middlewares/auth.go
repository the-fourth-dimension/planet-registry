package middlewares

import (
	"fmt"
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
		floatingRole, ok := claims["role"].(float64)
		if !ok {
			ctx.Error(HttpError.NewHttpError("missing claim", "role", http.StatusBadRequest))
			ctx.Abort()
			return
		}
		role := int(floatingRole)
		if !roles.IsValidRole(role) {
			ctx.Error(HttpError.NewHttpError("invalid role", fmt.Sprintf("%d", role), http.StatusUnauthorized))
			ctx.Abort()
			return
		}
		if role == 1 {
			adminUsername, ok := claims["username"].(string)
			if !ok {
				ctx.Error(HttpError.NewHttpError("missing claim", "username", http.StatusForbidden))
				ctx.Abort()
				return
			}
			ctx.Set("username", adminUsername)
		}
		ctx.Set("role", role)
		ctx.Next()
	}
}
