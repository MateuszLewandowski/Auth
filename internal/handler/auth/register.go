package auth

import (
	"Auth/config"
	"Auth/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(repo UserRepository, tokenConfig config.JWTConfig) gin.HandlerFunc {
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

		ctx.JSON(http.StatusCreated, gin.H{"message": "user created"})
	}
}
