package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthHandler(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}

		if !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			if secret == "" {
				return nil, fmt.Errorf("secret cannot be empty")
			}
			return []byte(secret), nil
		})

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		if err := validateTokenClaims(claims); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token validation failed: " + err.Error()})
			return
		}

		if username, exists := claims["username"]; exists {
			ctx.Set("username", username)
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "username claim required"})
			return
		}

		ctx.Next()
	}
}

func validateTokenClaims(claims jwt.MapClaims) error {
	now := time.Now()

	shouldReturn, err := verifyExpExpiryTime(claims, now)
	if shouldReturn {
		return err
	}

	shouldReturn, err = verifyNbfExpiryDate(claims, now)
	if shouldReturn {
		return err
	}

	return nil
}

func verifyExpExpiryTime(claims jwt.MapClaims, now time.Time) (bool, error) {
	if exp, ok := claims["exp"].(float64); ok {
		if now.After(time.Unix(int64(exp), 0)) {
			return true, fmt.Errorf("token expired")
		}
	}
	return false, nil
}

func verifyNbfExpiryDate(claims jwt.MapClaims, now time.Time) (bool, error) {
	if nbf, ok := claims["nbf"].(float64); ok {
		if now.Before(time.Unix(int64(nbf), 0)) {
			return true, fmt.Errorf("token not valid yet")
		}
	}
	return false, nil
}
