package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		findQuery := models.Admin{Username: input.Username, Password: input.Password}

		findAdminResult := a.FindFirst(&findQuery)
		if findAdminResult.Error != nil {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Invalid username or password",
			})
			return
		}
		ctx.Next()
	}
}
