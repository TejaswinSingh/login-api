package user

import (
	"strings"

	"github.com/TejaswinSingh/login-api/internal/constants"
	common_errors "github.com/TejaswinSingh/login-api/internal/errors"
)

type createNewUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type createNewUserResponse struct {
	Status int `json:"-"`
}

func (r createNewUserResponse) getStatus() int { return r.Status }
func (r createNewUserRequest) clean() createNewUserRequest {
	return createNewUserRequest{
		Username: strings.TrimSpace(r.Username),
		Password: strings.TrimSpace(r.Password),
	}
}
func (r createNewUserRequest) validate() error {
	if r.Username == "" || r.Password == "" {
		return common_errors.ErrReqMissingFields
	}
	if len(r.Username) < constants.MIN_USER_USERNAME_LEN || len(r.Username) > constants.MAX_USER_USERNAME_LEN {
		return common_errors.ErrInvalidUsernameLength
	}
	length := len([]byte(r.Password))
	if length < constants.MIN_USER_PASSWORD_LEN_BYTES || length > constants.MAX_USER_PASSWORD_LEN_BYTES {
		return common_errors.ErrInvalidPasswordLength
	}
	return nil
}
