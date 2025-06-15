package main

import (
	"Auth/config"
	"Auth/internal/server"
	"Auth/pkg"
)

// func authMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		tokenString := c.GetHeader("Authorization")
// 		if tokenString == "" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "brak tokena"})
// 			return
// 		}

// 		claims := &Claims{}
// 		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
// 			return jwtSecret, nil
// 		})
// 		if err != nil || !token.Valid {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "nieprawidłowy token"})
// 			return
// 		}

// 		c.Set("username", claims.Username)
// 		c.Next()
// 	}
// }

func main() {

	cfg := config.LoadConfig()

	// Init database
	pkg.InitializeDatabase(cfg)
	pkg.InitializeRedis(cfg)
	pkg.InitializeHandler(cfg)

	server := server.StartServer()

	server.Run(":" + cfg.Server.Port)

	// 	authGroup := router.Group("/protected")
	// 	authGroup.Use(authMiddleware())
	// 	authGroup.GET("/profile", func(c *gin.Context) {
	// 		username, _ := c.Get("username")
	// 		c.JSON(http.StatusOK, gin.H{"message": "access granted", "user": username})
	// 	})

	//		router.Run(":8080")
	//	}
}
