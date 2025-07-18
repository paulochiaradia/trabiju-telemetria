package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/paulochiaradia/trabiju-telemetria/internal/controllers"
)

// RegisterRoutes sets up the routes for the application.
func RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", controllers.PingHandler)
}
