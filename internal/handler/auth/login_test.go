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
	resp := httptest.NewRecorder()

	r := gin.New()
	r.POST("/login", auth.LoginHandler(repo, tokenCfg, cache))
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "token")
}

func TestLoginHandler_CacheMissAndDBHit(t *testing.T) {
	password := "testpass"
	hashed := hashPassword(password)

	cache := &MockCacheRepository{
		GetFunc: func(key string) (string, error) {
			return "", errors.New("not found")
		},
		SetFunc: func(key string, value any, expiration time.Duration) error {
			return nil
		},
	}

	repo := &MockUserRepository{
		FindByUsernameFunc: func(ctx context.Context, username string) (*model.User, error) {
			return &model.User{
				Username: "dbuser",
				Password: hashed,
			}, nil
		},
	}

	tokenCfg := config.JWTConfig{
		Secret:            "s3cr3t",
		ExpirationMinutes: 10,
	}

	input := auth.AuthInput{Username: "dbuser", Password: password}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r := gin.New()
	r.POST("/login", auth.LoginHandler(repo, tokenCfg, cache))
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "token")
}

func TestLoginHandler_CacheWrongPassword(t *testing.T) {
	cache := &MockCacheRepository{
		GetFunc: func(key string) (string, error) {
			return hashPassword("correctpassword"), nil
		},
	}

	repo := &MockUserRepository{}

	tokenCfg := config.JWTConfig{
		Secret:            "topsecret",
		ExpirationMinutes: 5,
	}

	input := auth.AuthInput{Username: "user", Password: "wrong"}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r := gin.New()
	r.POST("/login", auth.LoginHandler(repo, tokenCfg, cache))
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "invalid credentials")
}

func TestLoginHandler_InvalidJSON(t *testing.T) {
	cache := &MockCacheRepository{}
	repo := &MockUserRepository{}
	tokenCfg := config.JWTConfig{}

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte("not-json")))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r := gin.New()
	r.POST("/login", auth.LoginHandler(repo, tokenCfg, cache))
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "corrupted input")
}

func TestLoginHandler_UserNotFoundInCacheAndDB(t *testing.T) {
	cache := &MockCacheRepository{
		GetFunc: func(key string) (string, error) {
			return "", errors.New("not found")
		},
	}

	repo := &MockUserRepository{
		FindByUsernameFunc: func(ctx context.Context, username string) (*model.User, error) {
			return nil, errors.New("not found")
		},
	}

	tokenCfg := config.JWTConfig{}

	input := auth.AuthInput{Username: "nonexistent", Password: "any"}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r := gin.New()
	r.POST("/login", auth.LoginHandler(repo, tokenCfg, cache))
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "invalid credentials")
}
