package auth_test

import (
	"Auth/config"
	"Auth/internal/handler/auth"
	"Auth/internal/model"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	input := auth.AuthInput{
		Username: "testuser",
		Password: "testpass",
	}
	body, _ := json.Marshal(input)

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	repo := &MockUserRepository{
		CreateFunc: func(user *model.User) error {
			return nil
		},
	}

	cache := &MockCacheRepository{
		SetFunc: func(key string, value any, expiration time.Duration) error {
			return nil
		},
	}

	router := gin.Default()
	router.POST("/register", auth.RegisterHandler(repo, config.JWTConfig{}, cache))
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Contains(t, resp.Body.String(), "user created")
}

func TestRegisterHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer([]byte("invalid-json")))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router := gin.Default()
	router.POST("/register", auth.RegisterHandler(&MockUserRepository{}, config.JWTConfig{}, &MockCacheRepository{}))
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "invalid input")
}

func TestRegisterHandler_UserAlreadyExists(t *testing.T) {
	gin.SetMode(gin.TestMode)

	input := auth.AuthInput{
		Username: "existinguser",
		Password: "password123",
	}
	body, _ := json.Marshal(input)

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	repo := &MockUserRepository{
		CreateFunc: func(user *model.User) error {
			return errors.New("user already exists")
		},
	}

	router := gin.Default()
	router.POST("/register", auth.RegisterHandler(repo, config.JWTConfig{}, &MockCacheRepository{}))
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusConflict, resp.Code)
	assert.Contains(t, resp.Body.String(), "user already exists")
}

func TestRegisterHandler_CacheError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	input := auth.AuthInput{
		Username: "testuser",
		Password: "testpass",
	}
	body, _ := json.Marshal(input)

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	repo := &MockUserRepository{
		CreateFunc: func(user *model.User) error {
			return nil
		},
	}

	cache := &MockCacheRepository{
		SetFunc: func(key string, value any, expiration time.Duration) error {
			return errors.New("redis error")
		},
	}

	router := gin.Default()
	router.POST("/register", auth.RegisterHandler(repo, config.JWTConfig{}, cache))
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Contains(t, resp.Body.String(), "could not save to redis")
}
