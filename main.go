package main

import (
	"log"

	_ "github.com/123DaNIS123/UsersSegments/docs"

	"github.com/123DaNIS123/UsersSegments/db"
	"github.com/123DaNIS123/UsersSegments/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title User Segments API
// @version 1.0
// @descriotion A service for managing segments and user belongings

// @contact.name   Danis Nizamutdinov
// @contact.email  danisnizamutdinov1@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host	localhost:8080
// @BasePath /
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables %s", err.Error())
	}
	db.Connect()
	router := gin.New()
	routes.Route(router)
	router.Run(":8080")
}
