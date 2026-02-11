package user

import "context"

type UserRepository interface {
	GetUserFromID(ctx context.Context, id int) (*User, error)
	CreateNewUser(ctx context.Context, user *User) error
	ValidatePassword(ctx context.Context, username, password string) error
	GeneratePasswordHash(ctx context.Context, password string) (string, error)
}
