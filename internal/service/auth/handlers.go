package auth

import (
	"context"
	"net/http"
)

func (s *authService) login(ctx context.Context, r loginRequest) (*loginResponse, error) {

	if err := s.userRepository.ValidatePassword(ctx, r.Username, r.Password); err != nil {
		return nil, err
	}

	token, err := s.jwtModule.CreateNewJWT(ctx, r.Username)
	if err != nil {
		return nil, err
	}

	return &loginResponse{Token: token, Status: http.StatusOK}, nil
}
