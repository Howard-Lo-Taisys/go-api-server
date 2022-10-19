package routes

import (
	"go-api-server/controller"

	"github.com/gin-gonic/gin"
)

func ServiceRoute(router *gin.Engine) {
	router.GET("/healthz", controller.GetHealthz)
	router.GET("/api/:service-name", controller.GetServerName)
	router.GET("/api/:service-name/:version-number", controller.GetServerVersion)
	router.POST("/api/:service-name", controller.PostServerName)
	router.PUT("/api/:service-name/:version-number", controller.PutServerVersion)
	router.DELETE("/api/:service-name/:version-number", controller.DeleteServerVersion)
}
