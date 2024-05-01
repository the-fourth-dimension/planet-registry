package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/the_fourth_dimension/planet_registry/pkg/errors/HttpError"
)

func SuperuserMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, exists := ctx.Get("tokenClaims")
		if !exists {
			ctx.Error(HttpError.NewHttpError("invalid claims", "", http.StatusForbidden))
			return
		}
		typedClaims := claims.(jwt.MapClaims)
		role, ok := typedClaims["role"].(string)
		if !ok {
			ctx.Error(HttpError.NewHttpError("missing claim", "role", http.StatusForbidden))
			return
		}
		if role != "superuser" {
			ctx.Error(HttpError.NewHttpError("invalid role", role, http.StatusUnauthorized))
			return
		}
		ctx.Next()
	}
}
