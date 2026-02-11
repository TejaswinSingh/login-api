package user

import (
	"context"
	"net/http"

	user "github.com/TejaswinSingh/login-api/internal/repository/user"
)

func (s *userService) createNewUser(ctx context.Context, r createNewUserRequest) (*createNewUserResponse, error) {
	passhash, err := s.userRepository.GeneratePasswordHash(ctx, r.Password)
	if err != nil {
		return nil, err
	}

	user := user.User{
		Username: r.Username,
		Passhash: passhash,
	}

	if err := s.userRepository.CreateNewUser(ctx, &user); err != nil {
		return nil, err
	}

	return &createNewUserResponse{Status: http.StatusCreated}, nil
}
