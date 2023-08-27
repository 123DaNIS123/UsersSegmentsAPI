package main

import (
	"log"

	"github.com/123DaNIS123/UsersSegments/config"
	"github.com/123DaNIS123/UsersSegments/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables %s", err.Error())
	}
	config.Connect()
	router := gin.New()
	routes.Route(router)
	router.Run(":8080")
}
