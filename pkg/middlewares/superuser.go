package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/the_fourth_dimension/planet_registry/pkg/errors/HttpError"
)

func SuperuserMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role := ctx.GetInt("role")
		if role > 0 {
			ctx.Error(HttpError.NewHttpError("invalid role", fmt.Sprintf("%d", role), http.StatusUnauthorized))
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
