package transport

import (
	"encoding/json"
	"io"
	"net/http"
	"store/pkg/store/infrastructure/jwt"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"

	"store/pkg/store/app"
	"store/pkg/store/domain"
)

const (
	currentUserEndpoint = PathPrefix + "user"
	specUserEndpoint    = PathPrefix + "user/{id}"
)

const (
	errorCodeUnknown      = 0
	errorCodeUserNotFound = 1
)

const authTokenHeader = "X-Auth-Token"

type Server interface {
	Start()
}

func NewServer(router *mux.Router, tokenParser jwt.TokenParser, userService app.UserService, userQueryService app.UserQueryService) Server {
	return &server{
		router:           router,
		tokenParser:      tokenParser,
		userService:      userService,
		userQueryService: userQueryService,
	}
}

type server struct {
	router           *mux.Router
	tokenParser      jwt.TokenParser
	userService      app.UserService
	userQueryService app.UserQueryService
}

type currentUserInfo struct {
	UserID string `json:"id"`
}

type errorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (s *server) Start() {
	s.router.HandleFunc(currentUserEndpoint, s.getCurrentUserIDHandler).Methods(http.MethodGet)
	s.router.HandleFunc(specUserEndpoint, s.removeUserEndpoint).Methods(http.MethodDelete)
	s.router.HandleFunc(specUserEndpoint, s.updateUserEndpoint).Methods(http.MethodPut)
	s.router.HandleFunc(specUserEndpoint, s.getUserEndpoint).Methods(http.MethodGet)
}

func (s *server) getCurrentUserIDHandler(w http.ResponseWriter, r *http.Request) {
	tokenData, err := s.extractAuthorizationData(r)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	writeResponse(w, currentUserInfo{UserID: tokenData.UserID()})
}

func (s *server) removeUserEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	tokenData, err := s.extractAuthorizationData(r)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	if tokenData.UserID() != id {
		writeErrorResponse(w, errForbidden)
		return
	}

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

	tokenData, err := s.extractAuthorizationData(r)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	if tokenData.UserID() != id {
		writeErrorResponse(w, errForbidden)
		return
	}

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

	err = s.userService.UpdateUser(userID, tokenData.UserLogin(), data.Firstname, data.Lastname, data.Email, data.Phone)
	if errors.Cause(err) == domain.ErrUserNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, http.StatusText(http.StatusOK))
}

func (s *server) getUserEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	tokenData, err := s.extractAuthorizationData(r)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	if tokenData.UserID() != id {
		writeErrorResponse(w, errForbidden)
		return
	}

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
		userData = app.UserData{ID: userID}
	} else if err != nil {
		writeErrorResponse(w, err)
		return
	}
	userData.Username = tokenData.UserLogin()

	userDataJson, err := json.Marshal(userData)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(userDataJson))
	w.WriteHeader(http.StatusOK)
}

func (s *server) extractAuthorizationData(r *http.Request) (jwt.TokenData, error) {
	token := r.Header.Get(authTokenHeader)
	if token == "" {
		return nil, errForbidden
	}
	tokenData, err := s.tokenParser.ParseToken(token)
	if err != nil {
		return nil, errors.Wrap(errForbidden, err.Error())
	}
	return tokenData, nil
}

func writeResponse(w http.ResponseWriter, response interface{}) {
	js, err := json.Marshal(response)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(js)
}

func writeErrorResponse(w http.ResponseWriter, err error) {
	info := errorInfo{Code: errorCodeUnknown, Message: err.Error()}
	switch errors.Cause(err) {
	case app.ErrUserNotExists:
		info.Code = errorCodeUserNotFound
		w.WriteHeader(http.StatusNotFound)
	case errForbidden:
		w.WriteHeader(http.StatusForbidden)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	js, _ := json.Marshal(info)
	_, _ = w.Write(js)
}

var errForbidden = errors.New("access denied")
