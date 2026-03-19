package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func getting(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"method": "GET",
	})
}

func pingStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ping time: " + time.Now().Format("2006-01-02 15:04:05"),
	})
}

func main() {
	router := gin.Default()

	router.GET("/ping", pingStatus)
	router.GET("pong", getting)

	router.Run() // listens on 0.0.0.0:8080 by default
}
