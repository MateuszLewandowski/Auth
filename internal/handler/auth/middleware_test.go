package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const testSecret = "test-secret"

func generateToken(t *testing.T, claims jwt.MapClaims, secret string) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	return tokenString
}

func performRequest(_ *testing.T, handler gin.HandlerFunc, token string) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(handler)
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "authorized"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	if token != "" {
		req.Header.Set("Authorization", token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestAuthHandler_MissingAuthorizationHeader(t *testing.T) {
	responseRecorder := performRequest(t, AuthHandler(testSecret), "")
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
	assert.Contains(t, responseRecorder.Body.String(), "authorization header required")
}

func TestAuthHandler_InvalidTokenFormat(t *testing.T) {
	responseRecorder := performRequest(t, AuthHandler(testSecret), "Token abc.def.ghi")
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
	assert.Contains(t, responseRecorder.Body.String(), "invalid token format")
}

func TestAuthHandler_InvalidSignature(t *testing.T) {
	claims := jwt.MapClaims{
		"username": "testuser",
		"exp":      time.Now().Add(5 * time.Minute).Unix(),
	}
	token := generateToken(t, claims, "wrong-secret")
	responseRecorder := performRequest(t, AuthHandler(testSecret), "Bearer "+token)
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
	assert.Contains(t, responseRecorder.Body.String(), "signature is invalid")
}

func TestAuthHandler_ExpiredToken(t *testing.T) {
	claims := jwt.MapClaims{
		"username": "testuser",
		"exp":      time.Now().Add(-5 * time.Minute).Unix(),
	}
	token := generateToken(t, claims, testSecret)
	responseRecorder := performRequest(t, AuthHandler(testSecret), "Bearer "+token)
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
	assert.Contains(t, responseRecorder.Body.String(), "token is expired")
}

func TestAuthHandler_NotBeforeInFuture(t *testing.T) {
	claims := jwt.MapClaims{
		"username": "testuser",
		"nbf":      time.Now().Add(10 * time.Minute).Unix(),
		"exp":      time.Now().Add(20 * time.Minute).Unix(),
	}
	token := generateToken(t, claims, testSecret)
	responseRecorder := performRequest(t, AuthHandler(testSecret), "Bearer "+token)
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
	assert.Contains(t, responseRecorder.Body.String(), "token has invalid claims: token is not valid yet")
}

func TestAuthHandler_MissingUsernameClaim(t *testing.T) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(5 * time.Minute).Unix(),
	}
	token := generateToken(t, claims, testSecret)
	responseRecorder := performRequest(t, AuthHandler(testSecret), "Bearer "+token)
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
	assert.Contains(t, responseRecorder.Body.String(), "username claim required")
}

func TestAuthHandler_ValidToken(t *testing.T) {
	claims := jwt.MapClaims{
		"username": "validuser",
		"exp":      time.Now().Add(5 * time.Minute).Unix(),
		"nbf":      time.Now().Add(-5 * time.Minute).Unix(),
	}
	token := generateToken(t, claims, testSecret)
	responseRecorder := performRequest(t, AuthHandler(testSecret), "Bearer "+token)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Contains(t, responseRecorder.Body.String(), "authorized")
}
