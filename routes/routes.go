package routes

import (
	"github.com/123DaNIS123/UsersSegments/controller"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Route(router *gin.Engine) {
	//add swagger
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/", controller.Index)

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
	router.POST("/userbinds", controller.GetUserBinds)
	router.POST("/timedata", controller.GetTimeData)
}
