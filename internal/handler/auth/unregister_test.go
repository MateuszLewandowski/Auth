package auth_test

import (
	"Auth/config"
	"Auth/internal/handler/auth"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUnregisterHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		username           any
		usernameExists     bool
		authHeader         string
		setupCache         func() *MockCacheRepository
		setupRepo          func() *MockUserRepository
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "missing username in context",
			username:           nil,
			usernameExists:     false,
			authHeader:         "",
			setupCache:         func() *MockCacheRepository { return &MockCacheRepository{} },
			setupRepo:          func() *MockUserRepository { return &MockUserRepository{} },
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `{"error":"unauthenticated"}`,
		},
		{
			name:               "invalid username type",
			username:           123,
			usernameExists:     true,
			authHeader:         "",
			setupCache:         func() *MockCacheRepository { return &MockCacheRepository{} },
			setupRepo:          func() *MockUserRepository { return &MockUserRepository{} },
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `{"error":"invalid username type"}`,
		},
		{
			name:           "cache delete error",
			username:       "user1",
			usernameExists: true,
			authHeader:     "Bearer token123",
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
			authHeader:     "Bearer token123",
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
			name:           "successful delete without token",
			username:       "user1",
			usernameExists: true,
			authHeader:     "",
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
		{
			name:           "successful delete with token",
			username:       "user1",
			usernameExists: true,
			authHeader:     "Bearer token123",
			setupCache: func() *MockCacheRepository {
				return &MockCacheRepository{
					DeleteFunc: func(key string) error {
						if key == "user:user1" {
							return nil
						}
						return errors.New("unexpected key")
					},
					SetFunc: func(key string, value any, expiration time.Duration) error {
						assert.Equal(t, "jwt_blacklist:token123", key)
						assert.Equal(t, "invalid", value)
						assert.Equal(t, 15*time.Minute, expiration)
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

			ctx.Request = &http.Request{
				Header: make(http.Header),
			}

			if tt.usernameExists {
				ctx.Set("username", tt.username)
			}

			if tt.authHeader != "" {
				ctx.Request.Header.Set("Authorization", tt.authHeader)
			}

			handler := auth.UnregisterHandler(repo, cache, config.JWTConfig{
				ExpirationMinutes: 15,
			})

			handler(ctx)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}
