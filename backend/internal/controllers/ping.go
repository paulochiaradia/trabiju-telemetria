package controllers

import "github.com/gin-gonic/gin"

// PingHandler responds with a simple JSON message to indicate that the server is running.
func PingHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}
