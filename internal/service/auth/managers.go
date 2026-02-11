package auth

import (
	"errors"
	"net/http"

	common_errors "github.com/TejaswinSingh/login-api/internal/errors"
)

func (s *authService) loginManager(w http.ResponseWriter, r *http.Request) {
	req, err := s.loginDecoder(w, r)
	if err != nil {
		s.errorEncoder(w, r, &ErrorWithStatus{Err: err, Status: http.StatusBadRequest})
		return
	}

	resp, err := s.login(r.Context(), req)
	if err != nil {
		if errors.Is(err, common_errors.ErrNotFound) {
			s.errorEncoder(w, r, &ErrorWithStatus{Err: err, Status: http.StatusNotFound})
		} else if errors.Is(err, common_errors.ErrInvalidCreds) {
			s.errorEncoder(w, r, &ErrorWithStatus{Err: err, Status: http.StatusUnauthorized})
		} else {
			s.errorEncoder(w, r, &ErrorWithStatus{Err: common_errors.ErrInternalServer, Status: http.StatusInternalServerError})
		}
		return
	}

	s.encoder(w, r, resp)
}
