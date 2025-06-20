package auth_test

import (
	"Auth/internal/handler/auth"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUnregisterHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		username           any
		usernameExists     bool
		setupCache         func() *MockCacheRepository
		setupRepo          func() *MockUserRepository
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "missing username in context",
			username:           nil,
			usernameExists:     false,
			setupCache:         func() *MockCacheRepository { return &MockCacheRepository{} },
			setupRepo:          func() *MockUserRepository { return &MockUserRepository{} },
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `{"error":"unauthenticated"}`,
		},
		{
			name:               "invalid username type",
			username:           123,
			usernameExists:     true,
			setupCache:         func() *MockCacheRepository { return &MockCacheRepository{} },
			setupRepo:          func() *MockUserRepository { return &MockUserRepository{} },
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `{"error":"invalid username type"}`,
		},
		{
			name:           "cache delete error",
			username:       "user1",
			usernameExists: true,
			setupCache: func() *MockCacheRepository {
				return &MockCacheRepository{
					DeleteFunc: func(key string) error {
						assert.Equal(t, "user:user1", key)
						return errors.New("cache error")
					},
				}
			},
			setupRepo: func() *MockUserRepository {
				return &MockUserRepository{}
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedBody:       `{"error":"user not found"}`,
		},
		{
			name:           "repo delete error",
			username:       "user1",
			usernameExists: true,
			setupCache: func() *MockCacheRepository {
				return &MockCacheRepository{
					DeleteFunc: func(key string) error {
						assert.Equal(t, "user:user1", key)
						return nil
					},
				}
			},
			setupRepo: func() *MockUserRepository {
				return &MockUserRepository{
					DeleteFunc: func(username string) error {
						assert.Equal(t, "user1", username)
						return errors.New("repo error")
					},
				}
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedBody:       `{"error":"user not found"}`,
		},
		{
			name:           "successful delete",
			username:       "user1",
			usernameExists: true,
			setupCache: func() *MockCacheRepository {
				return &MockCacheRepository{
					DeleteFunc: func(key string) error {
						assert.Equal(t, "user:user1", key)
						return nil
					},
				}
			},
			setupRepo: func() *MockUserRepository {
				return &MockUserRepository{
					DeleteFunc: func(username string) error {
						assert.Equal(t, "user1", username)
						return nil
					},
				}
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message":"user deleted"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := tt.setupCache()
			repo := tt.setupRepo()

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			if tt.usernameExists {
				ctx.Set("username", tt.username)
			}

			handler := auth.UnregisterHandler(repo, cache)

			handler(ctx)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}
