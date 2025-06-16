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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	User  *model.User
	Error error
}

// Create implements auth.UserRepository.
func (m *MockUserRepository) Create(user *model.User) error {
	panic("unimplemented")
}

func (m *MockUserRepository) FindUserByUsername(username string) (*model.User, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.User, nil
}

func hashPassword(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed)
}

func TestLoginHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	password := "secret123"
	username := "testuser"
	user := &model.User{
		Username: username,
		Password: hashPassword(password),
	}

	repo := &MockUserRepository{User: user}
	tokenConfig := config.JWTConfig{
		Secret:            "testsecret",
		ExpirationMinutes: 60,
	}

	input := auth.AuthInput{
		Username: username,
		Password: password,
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()

	r := gin.New()
	r.POST("/login", auth.LoginHandler(repo, tokenConfig))

	r.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("unexpected http response %d", responseRecorder.Code)
	}

	var resp map[string]string
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("corrupted json payload %v", err)
	}

	token, ok := resp["token"]

	assert.True(t, ok)
	assert.NotEmpty(t, token, "token should not be empty")
}

func TestLoginHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &MockUserRepository{}
	tokenConfig := config.JWTConfig{}

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()

	r := gin.New()
	r.POST("/login", auth.LoginHandler(repo, tokenConfig))

	r.ServeHTTP(responseRecorder, req)

	assert.Equal(t, `{"error":"corrupted input payload"}`, responseRecorder.Body.String())
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}

func TestLoginHandler_UserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &MockUserRepository{Error: errors.New("user not found")}
	tokenConfig := config.JWTConfig{}

	input := auth.AuthInput{
		Username: "unknown",
		Password: "pass",
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()

	r := gin.New()
	r.POST("/login", auth.LoginHandler(repo, tokenConfig))

	r.ServeHTTP(responseRecorder, req)

	assert.Equal(t, `{"error":"invalid credentials"}`, responseRecorder.Body.String())
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
}

func TestLoginHandler_WrongPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	user := &model.User{
		Username: "testuser",
		Password: hashPassword("correctpassword"),
	}
	repo := &MockUserRepository{User: user}
	tokenConfig := config.JWTConfig{}

	input := auth.AuthInput{
		Username: "testuser",
		Password: "wrongpassword",
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()

	r := gin.New()
	r.POST("/login", auth.LoginHandler(repo, tokenConfig))

	r.ServeHTTP(responseRecorder, req)

	assert.Equal(t, `{"error":"invalid credentials"}`, responseRecorder.Body.String())
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
}
