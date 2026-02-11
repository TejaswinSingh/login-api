package jwt

import (
	"context"
	"log/slog"
	"time"

	"github.com/TejaswinSingh/login-api/internal/config"
	common_errors "github.com/TejaswinSingh/login-api/internal/errors"
	"github.com/golang-jwt/jwt/v5"
)

type JWT string

const (
	TOKEN_EXPIRY = time.Hour * 1
)

type JwtModule struct {
	logger *slog.Logger
	config config.Config
}

func NewJwtModule(logger *slog.Logger, config config.Config) *JwtModule {
	return &JwtModule{logger: logger, config: config}
}

func (m *JwtModule) CreateNewJWT(ctx context.Context, sub string) (JWT, error) {
	now := time.Now()
	exp := now.Add(TOKEN_EXPIRY)

	claims := jwt.MapClaims{
		"sub": sub,
		"iss": "loginapi.com", // issuing authority
		"aud": "example.com",  // target resource servers where the token will be sent
		"iat": now.Unix(),
		"exp": exp.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(m.config.SecretKey)
	if err != nil {
		m.logger.ErrorContext(ctx, "unable to create new JWT", "error", err)
		return "", common_errors.ErrInternalServer
	}
	return JWT(tokenString), nil
}
