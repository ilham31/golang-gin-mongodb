package main

import (
	"gin-mongodb/configs"
	"gin-mongodb/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// run database
	configs.ConnectDB()

	routes.UserRoute(router)
	router.Run("localhost:6000")
}
