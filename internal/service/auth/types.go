package auth

import (
	"strings"

	"github.com/TejaswinSingh/login-api/internal/constants"
	common_errors "github.com/TejaswinSingh/login-api/internal/errors"
	"github.com/TejaswinSingh/login-api/internal/service/auth/jwt"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token  jwt.JWT `json:"token"`
	Status int     `json:"-"`
}

func (r loginResponse) getStatus() int { return r.Status }
func (r loginRequest) clean() loginRequest {
	return loginRequest{
		Username: strings.TrimSpace(r.Username),
		Password: strings.TrimSpace(r.Password),
	}
}
func (r loginRequest) validate() error {
	if r.Username == "" || r.Password == "" {
		return common_errors.ErrReqMissingFields
	}
	length := len([]byte(r.Password))
	if length < constants.MIN_USER_PASSWORD_LEN_BYTES || length > constants.MAX_USER_PASSWORD_LEN_BYTES {
		return common_errors.ErrInvalidPasswordLength
	}
	return nil
}
