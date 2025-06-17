package auth

import (
	"Auth/internal/model"
	"context"
)

type UserRepository interface {
	Create(user *model.User) error
	FindUserByUsername(ctx context.Context, username string) (*model.User, error)
}
