package auth

import "time"

type Cache interface {
	Set(key string, value any, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
}