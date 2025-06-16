package auth

import "Auth/internal/model"

type UserRepository interface {
	Create(user *model.User) error
	FindUserByUsername(username string) (*model.User, error)
}
