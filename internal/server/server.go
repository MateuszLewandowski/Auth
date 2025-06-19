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
	router.POST("/login", auth.LoginHandler(db, cfg.JWT))
	router.POST("/register", auth.RegisterHandler(db, cfg.JWT, redis))
	router.GET("/auth", auth.AuthHandler(cfg.JWT.Secret)) // traefik sends get req

	return router
}
