package main

import (
	"Auth/config"
	"Auth/internal/server"
	"Auth/pkg"
	"os"
)

func main() {
	cfg := config.LoadConfig(os.Getenv("ENV_FILE"))

	db := pkg.InitializeDatabase(cfg)
	redis := pkg.InitializeRedis(cfg)

	server := server.StartServer(db, redis, cfg)

	err := server.Run(":" + cfg.Server.Port)

	if err != nil {
		panic(err)
	}
}
