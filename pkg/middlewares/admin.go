package middlewares

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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
		role := ctx.GetInt("role")
		if role > 1 {
			ctx.Error(HttpError.NewHttpError("invalid role", fmt.Sprintf("%d", role), http.StatusUnauthorized))
			ctx.Abort()
			return
		}
		if role < 1 {
			ctx.Next()
			return
		}
		username := ctx.GetString("username")

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
