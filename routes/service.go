package routes

import (
	"go-api-server/controller"
	"go-api-server/middleware"

	"github.com/gin-gonic/gin"
)

func ServiceRoute(router *gin.Engine) {

	router.GET("/healthz", controller.GetHealthz)

	visit := router.Group("")
	{
		visit.POST("/register", controller.Register)
		visit.POST("/login", controller.Login)
	}

	secure := router.Group("/api")
	secure.Use(middleware.JWTAuth())
	{
		secure.GET("/:service-name", controller.GetServerName)
		secure.GET("/:service-name/:version-id", controller.GetServerVersion)
		secure.POST("/:service-name", controller.PostServerName)
		secure.PUT("/:service-name/:version-id", controller.PutServerVersion)
		secure.DELETE("/:service-name/:version-id", controller.DeleteServerVersion)
	}
}
