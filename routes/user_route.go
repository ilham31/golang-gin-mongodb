package routes

import (
	"gin-mongodb/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	router.POST("/user", controllers.CreateUser)
	router.GET("/user/:userId", controllers.GetUser)
	router.GET("/users", controllers.GetAllUser)
	router.PUT("/user/:userId", controllers.EditUser)
	router.DELETE("/user/:userId", controllers.DeleteUser)
}
