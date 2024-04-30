package middlewares

import (
	"errors"
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
		var input AdminInput
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.Error(HttpError.NewHttpError("invalid input", err.Error(), http.StatusBadRequest))
			return
		}

		findQuery := models.Admin{Username: input.Username, Password: input.Password}

		findAdminResult := a.FindFirst(&findQuery)
		if findAdminResult.Error != nil {
			if errors.Is(findAdminResult.Error, gorm.ErrRecordNotFound) {
				ctx.Error(HttpError.NewHttpError("invalid username or password", findAdminResult.Error.Error(), http.StatusForbidden))
				return
			}
			ctx.AbortWithError(http.StatusInternalServerError, findAdminResult.Error)
			return
		}
		ctx.Next()
	}
}
