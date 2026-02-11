package service

import (
	"net/http"
)

type Service interface {
	RegisterHandlers(m *http.ServeMux)
}
