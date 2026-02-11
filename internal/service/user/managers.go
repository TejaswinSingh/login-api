package user

import (
	"errors"
	"net/http"

	common_errors "github.com/TejaswinSingh/login-api/internal/errors"
)

func (s *userService) createNewUserManager(w http.ResponseWriter, r *http.Request) {
	req, err := s.createNewUserDecoder(w, r)
	if err != nil {
		s.errorEncoder(w, r, &ErrorWithStatus{Err: err, Status: http.StatusBadRequest})
		return
	}

	resp, err := s.createNewUser(r.Context(), req)
	if err != nil {
		if errors.Is(err, common_errors.ErrResourceAlreadyExists) {
			s.errorEncoder(w, r, &ErrorWithStatus{Err: err, Status: http.StatusConflict})
		} else {
			s.errorEncoder(w, r, &ErrorWithStatus{Err: common_errors.ErrInternalServer, Status: http.StatusInternalServerError})
		}
		return
	}

	s.encoder(w, r, resp)
}
