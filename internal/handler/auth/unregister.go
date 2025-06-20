package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UnregisterHandler(
	repo UserDeleteRepository,
	cache Cache,
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

		ctx.JSON(http.StatusOK, gin.H{"message": "user deleted"})
	}
}
