package main

import (
	"go-api-server/config"
	"go-api-server/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	config.Connect()
	routes.ServiceRoute(router)

	router.Run(":8080")
}
