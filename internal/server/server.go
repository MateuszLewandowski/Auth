package server

import (
	"Auth/config"
	"Auth/internal/handler/auth"
	"Auth/pkg"

	"github.com/gin-gonic/gin"
)

func StartServer(db *pkg.UserGormRepository, redis *pkg.RedisRepository, cfg *config.Config) *gin.Engine {
	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	router.POST("/login", func(c *gin.Context) { auth.LoginHandler(db, cfg.JWT) })
	router.POST("/register", func(c *gin.Context) { auth.RegisterHandler(db, cfg.JWT) })
	router.POST("/auth", func(c *gin.Context) { auth.AuthHandler(cfg.JWT.Secret) })

	return router
}
