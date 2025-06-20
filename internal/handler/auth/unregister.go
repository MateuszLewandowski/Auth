package auth

import (
	"Auth/config"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func UnregisterHandler(
	repo UserDeleteRepository,
	cache Cache,
	tokenConfig config.JWTConfig,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username, exists := ctx.Get("username")

		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
			return
		}

		usernameStr, ok := username.(string)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid username type"})
			return
		}

		redisKey := fmt.Sprintf("user:%s", usernameStr)
		if err := cache.Delete(redisKey); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "user not found"})
			return
		}

		if err := repo.Delete(usernameStr); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "user not found"})
			return
		}

		authHeader := ctx.GetHeader("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		blacklistKey := fmt.Sprintf("jwt_blacklist:%s", tokenString)
		expiration := time.Duration(tokenConfig.ExpirationMinutes) * time.Minute
		_ = cache.Set(blacklistKey, "invalid", expiration)

		ctx.JSON(http.StatusOK, gin.H{"message": "user deleted"})
	}
}
