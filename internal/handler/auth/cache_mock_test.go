package auth_test

import (
	"errors"
	"time"
)

type MockCacheRepository struct {
	GetFunc    func(key string) (string, error)
	SetFunc    func(key string, value any, expiration time.Duration) error
	DeleteFunc func(key string) error
}

func (m *MockCacheRepository) Get(key string) (string, error) {
	if m.GetFunc != nil {
		return m.GetFunc(key)
	}
	return "", errors.New("not implemented")
}

func (m *MockCacheRepository) Delete(key string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(key)
	}
	return errors.New("not implemented")
}

func (m *MockCacheRepository) Set(key string, value any, expiration time.Duration) error {
	if m.SetFunc != nil {
		return m.SetFunc(key, value, expiration)
	}
	return nil
}
