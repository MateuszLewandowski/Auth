package auth

import (
	"Auth/internal/model"
	"context"
)

type UserCreateRepository interface {
	Create(user *model.User) error
}

type UserFindByUsernameRepository interface {
	FindByUsername(ctx context.Context, username string) (*model.User, error)
}

type UserDeleteRepository interface {
	Delete(username string) error
}
