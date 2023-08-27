package routes

import (
	"github.com/123DaNIS123/UsersSegments/controller"
	"github.com/gin-gonic/gin"
)

func Route(router *gin.Engine) {
	router.GET("/users", controller.GetUsers)
	router.POST("/user", controller.CreateUser)
	router.DELETE("/user/:id", controller.DeleteUser)
	router.PUT("/user/:id", controller.UpdateUser)

	router.GET("/segments", controller.GetSegments)
	router.GET("/segment/:id", controller.GetSegment)
	router.POST("/segment", controller.CreateSegment)
	router.DELETE("/segment", controller.DeleteSegment)
	router.PUT("/segment/:id", controller.UpdateSegment)

	router.POST("/bind", controller.Bind)
	router.GET("/binds", controller.GetBinds)
}
