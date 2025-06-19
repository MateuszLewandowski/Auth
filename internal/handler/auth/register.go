package auth

import (
	"Auth/config"
	"Auth/internal/model"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(repo UserRepository, tokenConfig config.JWTConfig, cache Cache) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input AuthInput
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "hash error"})
			return
		}

		user := model.User{
			Username: input.Username,
			Password: string(hashedPassword),
		}

		if err := repo.Create(&user); err != nil {
			ctx.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}

		userJSON, err := json.Marshal(user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not serialize user"})
			return
		}

		redisKey := fmt.Sprintf("user:%s", user.Username)
		if err := cache.Set(redisKey, userJSON, time.Hour*24); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not save to redis"})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "user created"})
	}
}
