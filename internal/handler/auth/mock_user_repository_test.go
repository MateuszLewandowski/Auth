package auth_test

import (
	"Auth/internal/model"
	"context"
	"errors"
)

type MockUserRepository struct {
	CreateFunc         func(user *model.User) error
	FindByUsernameFunc func(ctx context.Context, username string) (*model.User, error)
	DeleteFunc         func(username string) error
	User               *model.User
}

func (m *MockUserRepository) Create(user *model.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(user)
	}
	return nil
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	if m.FindByUsernameFunc != nil {
		return m.FindByUsernameFunc(ctx, username)
	}
	if m.User != nil && m.User.Username == username {
		return m.User, nil
	}
	return nil, errors.New("not found")
}

func (m *MockUserRepository) Delete(username string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(username)
	}
	return nil
}
