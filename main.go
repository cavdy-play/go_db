package main

import (
	"log"

	"github.com/gin-gonic/gin"

	config "github.com/cavdy-play/go_db/config"
	routes "github.com/cavdy-play/go_db/routes"
)

func main() {
	// Connect DB
	config.Connect()

	// Init Router
	router := gin.Default()

	// Route Handlers / Endpoints
	routes.Routes(router)

	log.Fatal(router.Run(":4747"))
}
