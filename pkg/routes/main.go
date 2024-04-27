package routes

import "github.com/gin-gonic/gin"

type Router struct {
	Engine *gin.Engine
}

func NewRouter() *Router {
	return &Router{
		Engine: gin.Default(),
	}
}
