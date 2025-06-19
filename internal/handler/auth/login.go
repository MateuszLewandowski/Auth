package auth

import (
	"Auth/config"
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(repo UserFindByUsernameRepository, tokenConfig config.JWTConfig, cache Cache) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input AuthInput
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "corrupted input payload"})
			return
		}

		var hashedPassword string
		foundInCache := false

		if val, err := cache.Get(input.Username); err == nil {
			hashedPassword = val
			foundInCache = true
		}

		if !foundInCache {
			reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), time.Second)
			defer cancel()

			user, err := repo.FindUserByUsername(reqCtx, input.Username)
			if err != nil {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
				return
			}

			hashedPassword = user.Password

			_ = cache.Set(user.Username, user.Password, 5*time.Minute)
		}

		if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input.Password)) != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		claims := jwt.MapClaims{
			"username": input.Username,
			"exp":      time.Now().Add(time.Duration(tokenConfig.ExpirationMinutes) * time.Minute).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(tokenConfig.Secret))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
}
