package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"

	"store/pkg/store/app"
	"store/pkg/store/domain"
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
	var data struct {
		Username  string `json:"username"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := s.userService.AddUser(data.Username, data.Firstname, data.Lastname, data.Email, data.Phone)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	io.WriteString(w, fmt.Sprintf(`{"id": "%s"}`, id))
	w.WriteHeader(http.StatusOK)
}

func (s *server) removeUserEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if len(id) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := uuid.FromString(id)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = s.userService.RemoveUser(userID)
	if errors.Cause(err) == domain.ErrUserNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) updateUserEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if len(id) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := uuid.FromString(id)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var data struct {
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.userService.UpdateUser(userID, data.Firstname, data.Lastname, data.Email, data.Phone)
	if errors.Cause(err) == domain.ErrUserNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) getUserEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if len(id) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := uuid.FromString(id)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userData, err := s.userQueryService.FindUser(userID)
	if errors.Cause(err) == app.ErrUserNotExists {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userDataJson, err := json.Marshal(userData)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(userDataJson))
	w.WriteHeader(http.StatusOK)
}
