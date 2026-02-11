package auth

import (
	"log/slog"
	"net/http"

	"github.com/TejaswinSingh/login-api/internal/config"
	"github.com/TejaswinSingh/login-api/internal/repository/user"
	"github.com/TejaswinSingh/login-api/internal/service"
	"github.com/TejaswinSingh/login-api/internal/service/auth/jwt"
)

type authService struct {
	logger         *slog.Logger
	config         config.Config
	userRepository user.UserRepository
	jwtModule      *jwt.JwtModule
}

func NewAuthService(logger *slog.Logger, config config.Config, jwtModule *jwt.JwtModule, userRepo user.UserRepository) service.Service {
	return &authService{
		logger:         logger,
		config:         config,
		jwtModule:      jwtModule,
		userRepository: userRepo,
	}
}

func (s *authService) RegisterHandlers(m *http.ServeMux) {
	m.HandleFunc("POST /login", s.loginManager)
}
