package auth

import (
	"Auth/internal/model"
	"context"
)

type UserCreateRepository interface {
	Create(user *model.User) error
}

type UserFindByUsernameRepository interface {
	FindUserByUsername(ctx context.Context, username string) (*model.User, error)
}
