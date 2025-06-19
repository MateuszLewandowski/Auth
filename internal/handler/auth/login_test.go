package auth_test

import (
	"Auth/config"
	"Auth/internal/handler/auth"
	"Auth/internal/model"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	User  *model.User
	Error error
}

func (m *MockUserRepository) FindUserByUsername(ctx context.Context, username string) (*model.User, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.User, nil
}

func hashPassword(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed)
}

func TestLoginHandler_CacheHit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	password := "test123"
	hashed := hashPassword(password)

	cache := &MockCacheRepository{
		GetFunc: func(key string) (string, error) {
			return hashed, nil
		},
	}

	repo := &MockUserRepository{}
	tokenCfg := config.JWTConfig{
		Secret:            "mysecret",
		ExpirationMinutes: 15,
	}

	input := auth.AuthInput{Username: "cacheduser", Password: password}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r := gin.New()
	r.POST("/login", auth.LoginHandler(repo, tokenCfg, cache))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")
}

func TestLoginHandler_CacheMissAndDBHit(t *testing.T) {
	password := "testpass"
	hashed := hashPassword(password)

	cache := &MockCacheRepository{
		GetFunc: func(key string) (string, error) {
			return "", errors.New("not found")
		},
		SetFunc: func(key string, value interface{}, expiration time.Duration) error {
			return nil
		},
	}

	repo := &MockUserRepository{
		User: &model.User{Username: "dbuser", Password: hashed},
	}

	tokenCfg := config.JWTConfig{Secret: "s3cr3t", ExpirationMinutes: 10}

	input := auth.AuthInput{Username: "dbuser", Password: password}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r := gin.New()
	r.POST("/login", auth.LoginHandler(repo, tokenCfg, cache))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")
}

func TestLoginHandler_CacheWrongPassword(t *testing.T) {
	cache := &MockCacheRepository{
		GetFunc: func(key string) (string, error) {
			return hashPassword("correctpassword"), nil
		},
	}

	repo := &MockUserRepository{} // won't be used
	tokenCfg := config.JWTConfig{Secret: "xxx", ExpirationMinutes: 1}

	input := auth.AuthInput{Username: "x", Password: "wrong"}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r := gin.New()
	r.POST("/login", auth.LoginHandler(repo, tokenCfg, cache))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid credentials")
}

func TestLoginHandler_InvalidJSON(t *testing.T) {
	cache := &MockCacheRepository{}
	repo := &MockUserRepository{}
	tokenCfg := config.JWTConfig{}

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte("not-json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r := gin.New()
	r.POST("/login", auth.LoginHandler(repo, tokenCfg, cache))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "corrupted input")
}

func TestLoginHandler_UserNotFoundInCacheAndDB(t *testing.T) {
	cache := &MockCacheRepository{
		GetFunc: func(key string) (string, error) {
			return "", errors.New("not found")
		},
	}
	repo := &MockUserRepository{
		Error: errors.New("not found"),
	}
	tokenCfg := config.JWTConfig{}

	input := auth.AuthInput{Username: "none", Password: "any"}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r := gin.New()
	r.POST("/login", auth.LoginHandler(repo, tokenCfg, cache))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid credentials")
}
