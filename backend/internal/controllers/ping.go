package controllers

import "github.com/gin-gonic/gin"

// PingHandler responde com uma mensagem JSON simples para indicar que o servidor está funcionando.
func PingHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}
