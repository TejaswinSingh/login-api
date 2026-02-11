package common_errors

import (
	"errors"
	"fmt"

	"github.com/TejaswinSingh/login-api/internal/constants"
)

var (
	ErrReqMissingFields      = errors.New("required fields were missing from request body")
	ErrInvalidCreds          = errors.New("invalid credentials")
	ErrInvalidPasswordLength = fmt.Errorf("password length must be between %d-%d bytes", constants.MIN_USER_PASSWORD_LEN_BYTES, constants.MAX_USER_PASSWORD_LEN_BYTES)
	ErrInvalidUsernameLength = fmt.Errorf("username length must be between %d-%d characters", constants.MIN_USER_USERNAME_LEN, constants.MAX_USER_USERNAME_LEN)
	ErrInternalServer        = errors.New("internal server error")
	ErrClientCloseRequest    = errors.New("client closed the request")
	ErrNotFound              = errors.New("resource not found")
	ErrResourceAlreadyExists = errors.New("resource already exists")
)
