package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Router struct {
	Engine *gin.Engine
	DB     *gorm.DB
}

func NewRouter(db *gorm.DB) *Router {
	return &Router{
		Engine: gin.Default(),
		DB:     db,
	}
}
