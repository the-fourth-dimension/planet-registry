package middlewares

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/gorm"
	"github.com/the_fourth_dimension/planet_registry/pkg/errors/HttpError"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
	"github.com/the_fourth_dimension/planet_registry/pkg/repositories"
)

type AdminInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AdminMiddleware(a *repositories.AdminRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, exists := ctx.Get("tokenClaims")
		if !exists {
			ctx.Error(HttpError.NewHttpError("missing claims", "", http.StatusBadRequest))
			ctx.Abort()
			return
		}
		typedClaims := claims.(jwt.MapClaims)
		role := typedClaims["role"].(string)
		if role != "admin" {
			ctx.Error(HttpError.NewHttpError("invalid role", role, http.StatusUnauthorized))
			ctx.Abort()
			return
		}
		username, ok := typedClaims["username"].(string)
		if !ok {
			ctx.Error(HttpError.NewHttpError("missing claim", "username", http.StatusForbidden))
			ctx.Abort()
			return
		}
		println(username)
		findQuery := models.Admin{Username: username}
		findAdminResult := a.FindFirst(&findQuery)
		if findAdminResult.Error != nil {
			if errors.Is(findAdminResult.Error, gorm.ErrRecordNotFound) {
				ctx.Error(HttpError.NewHttpError("invalid username", findAdminResult.Error.Error(), http.StatusForbidden))
				ctx.Abort()
				return
			}
			ctx.AbortWithError(http.StatusInternalServerError, findAdminResult.Error)
			return
		}
		ctx.Next()
	}
}
