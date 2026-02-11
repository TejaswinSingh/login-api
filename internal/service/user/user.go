package user

import (
	"log/slog"
	"net/http"

	"github.com/TejaswinSingh/login-api/internal/config"
	"github.com/TejaswinSingh/login-api/internal/repository/user"
	"github.com/TejaswinSingh/login-api/internal/service"
)

type userService struct {
	logger         *slog.Logger
	config         config.Config
	userRepository user.UserRepository
}

func NewUserService(logger *slog.Logger, config config.Config, userRepo user.UserRepository) service.Service {
	return &userService{
		logger:         logger,
		config:         config,
		userRepository: userRepo,
	}
}

func (s *userService) RegisterHandlers(m *http.ServeMux) {
	m.HandleFunc("POST /users/new", s.createNewUserManager)
}
