package server

import (
	"Auth/pkg"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	router.POST("/login", func(c *gin.Context) { pkg.Login(c) })
	router.POST("/register", func(c *gin.Context) { pkg.Register(c) })
	router.POST("/auth", func(c *gin.Context) { pkg.Auth(c) })

	return router
}
