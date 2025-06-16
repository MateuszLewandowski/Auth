package main

import (
	"Auth/config"
	"Auth/internal/server"
	"Auth/pkg"
)

func main() {
	cfg := config.LoadConfig(".env.dev")

	db := pkg.InitializeDatabase(cfg)
	redis := pkg.InitializeRedis(cfg)

	server := server.StartServer(db, redis, cfg)

	err := server.Run(":" + cfg.Server.Port)

	if err != nil {
		panic(err)
	}
}
