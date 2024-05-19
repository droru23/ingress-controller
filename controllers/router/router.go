package router

import (
	"github.com/gin-gonic/gin"
)

func SetupHealth(handler *Handler) *gin.Engine {
	r := gin.Default()
	r.GET("/health", handler.Health)
	return r
}

func SetupRouter(handler *Handler) *gin.Engine {
	r := gin.Default()
	r.GET("/", handler.Route)
	return r
}
