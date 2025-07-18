package main

import (
	"github.com/gin-gonic/gin"
	"github.com/paulochiaradia/trabiju-telemetria/internal/routes"
)

// main initializes the Gin router and registers the application routes.
// It then starts the server on port 8080.
func main() {
	r := gin.Default()
	routes.RegisterRoutes(r)
	r.Run(":8080")
}
