package transport

import (
	"github.com/gorilla/mux"

	"store/pkg/common/infrastructure/jwt"
)

type Server interface {
	Start()
}

func NewServer(router *mux.Router, tokenParser jwt.TokenParser) Server {
	return &server{
		router:      router,
		tokenParser: tokenParser,
	}
}

type server struct {
	router      *mux.Router
	tokenParser jwt.TokenParser
}

func (s *server) Start() {
}
