package routes

import (
	"github.com/123DaNIS123/UsersSegments/controller"
	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	router.GET("/", controller.UserController)
}
