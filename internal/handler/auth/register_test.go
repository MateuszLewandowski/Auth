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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// mock implementation
type mockUserRepo struct {
	createFunc             func(user *model.User) error
	findUserByUsernameFunc func(ctx context.Context, username string) (*model.User, error)
}

func (m *mockUserRepo) Create(user *model.User) error {
	return m.createFunc(user)
}

func (m *mockUserRepo) FindUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return m.findUserByUsernameFunc(ctx, username)
}

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

	repo := &mockUserRepo{
		createFunc: func(user *model.User) error {
			return nil
		},
	}

	r := gin.Default()
	r.POST("/register", auth.RegisterHandler(repo, config.JWTConfig{}))
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Contains(t, resp.Body.String(), "user created")
}

func TestRegisterHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer([]byte(`invalid-json`)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r := gin.Default()
	r.POST("/register", auth.RegisterHandler(nil, config.JWTConfig{}))
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "invalid input")
}

func TestRegisterHandler_RepoUserExistsError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	input := auth.AuthInput{
		Username: "exists",
		Password: "password",
	}
	body, _ := json.Marshal(input)

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	repo := &mockUserRepo{
		createFunc: func(user *model.User) error {
			return errors.New("user already exists")
		},
	}

	r := gin.Default()
	r.POST("/register", auth.RegisterHandler(repo, config.JWTConfig{}))
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusConflict, resp.Code)
	assert.Contains(t, resp.Body.String(), "user already exists")
}
