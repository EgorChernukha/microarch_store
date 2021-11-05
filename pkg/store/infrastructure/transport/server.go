package transport

import (
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"store/pkg/store/app"
)

type Server interface {
	Start()
}

func NewServer(router *mux.Router, userService app.UserService, userQueryService app.UserQueryService) Server {
	return &server{
		router:           router,
		userService:      userService,
		userQueryService: userQueryService,
	}
}

type server struct {
	router           *mux.Router
	userService      app.UserService
	userQueryService app.UserQueryService
}

func (s *server) Start() {
	s.router.HandleFunc("/api/v1/user", s.createUserEndpoint).Methods(http.MethodPost)
	s.router.HandleFunc("/api/v1/user/{id}", s.removeUserEndpoint).Methods(http.MethodDelete)
	s.router.HandleFunc("/api/v1/user/{id}", s.updateUserEndpoint).Methods(http.MethodPut)
	s.router.HandleFunc("/api/v1/user/{id}", s.getUserEndpoint).Methods(http.MethodGet)
}

func (s *server) createUserEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	firstname := vars["firstname"]
	lastname := vars["lastname"]
	email := vars["email"]
	phone := vars["phone"]

	_ = username
	_ = firstname
	_ = lastname
	_ = email
	_ = phone

	io.WriteString(w, `{"id": atata}`)
	w.WriteHeader(http.StatusOK)
}

func (s *server) removeUserEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_ = id

	io.WriteString(w, `{"id": atata}`)
	w.WriteHeader(http.StatusOK)
}

func (s *server) updateUserEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	firstname := vars["firstname"]
	lastname := vars["lastname"]
	email := vars["email"]
	phone := vars["phone"]

	_ = id
	_ = firstname
	_ = lastname
	_ = email
	_ = phone

	io.WriteString(w, `{"id": atata}`)
	w.WriteHeader(http.StatusOK)
}

func (s *server) getUserEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_ = id

	io.WriteString(w, `{"id": atata}`)
	w.WriteHeader(http.StatusOK)
}
