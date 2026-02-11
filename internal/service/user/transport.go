package user

import (
	"encoding/json"
	"net/http"

	"github.com/TejaswinSingh/login-api/internal/constants"
	common_errors "github.com/TejaswinSingh/login-api/internal/errors"
)

func (s *userService) createNewUserDecoder(_ http.ResponseWriter, r *http.Request) (createNewUserRequest, error) {
	var req createNewUserRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)
	if err != nil {
		return req, err
	}
	if err := req.clean().validate(); err != nil {
		return req, err
	}
	return req, nil
}

func (s *userService) encoder(w http.ResponseWriter, r *http.Request, resp any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	status := http.StatusOK
	if resp, ok := resp.(interface{ getStatus() int }); ok {
		status = resp.getStatus()
	}
	w.WriteHeader(status)

	if status == http.StatusNoContent {
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		s.logger.ErrorContext(r.Context(), "unable to encode response", "error", err)
	}
}

func (s *userService) errorEncoder(w http.ResponseWriter, r *http.Request, err *ErrorWithStatus) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err == nil || err.Err == nil || err.Status == 0 {
		err = &ErrorWithStatus{
			Err:    common_errors.ErrInternalServer,
			Status: http.StatusInternalServerError,
		}
	}

	if r.Context().Err() != nil { // if client closed the request
		err.Err = common_errors.ErrClientCloseRequest
		err.Status = constants.HTTP_STATUS_CLIENT_CLOSE_REQUEST
	}

	s.logger.ErrorContext(r.Context(), "", "error", err.Err.Error())

	w.WriteHeader(err.Status)

	if errEncode := json.NewEncoder(w).Encode(map[string]string{
		"error": err.Err.Error(),
	}); errEncode != nil {
		s.logger.ErrorContext(r.Context(), "unable to encode error response", "error", errEncode)
	}
}

type ErrorWithStatus struct {
	Err    error
	Status int
}
