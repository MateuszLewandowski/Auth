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

func LoginHandler(repo UserRepository, tokenConfig config.JWTConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input AuthInput
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "corrupted input payload"})
			return
		}

		reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), time.Second)
		defer cancel()

		user, err := repo.FindUserByUsername(reqCtx, input.Username)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		claims := jwt.MapClaims{
			"username": user.Username,
			"exp":      time.Now().Add(time.Duration(tokenConfig.ExpirationMinutes) * time.Minute).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(tokenConfig.Secret))

		ctx.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
}
