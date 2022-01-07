package transport

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"store/pkg/common/infrastructure/jwt"

	"store/pkg/order/app"
)

const PathPrefix = "/api/v1/"
const authTokenHeader = "X-Auth-Token"

const (
	createOrderEndpoint    = PathPrefix + "create"
	cancelOrderEndpoint    = PathPrefix + "{id}/cancel"
	getOrderStatusEndpoint = PathPrefix + "{id}/status"
	listOrdersEndpoint     = PathPrefix + "list"
)

const (
	errorCodeUnknown       = 0
	errorCodeOrderNotFound = 1
)

var errUnauthorized = errors.New("not authorized")
var errForbidden = errors.New("access denied")

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

type errorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (s *server) Start() {
	s.router.HandleFunc(createOrderEndpoint, s.createOrderEndpoint).Methods(http.MethodPost)
	s.router.HandleFunc(cancelOrderEndpoint, s.cancelOrderEndpoint).Methods(http.MethodPost)
	s.router.HandleFunc(getOrderStatusEndpoint, s.getOrderStatusEndpoint).Methods(http.MethodGet)
	s.router.HandleFunc(listOrdersEndpoint, s.listOrdersEndpoint).Methods(http.MethodGet)
}

func (s *server) createOrderEndpoint(w http.ResponseWriter, r *http.Request) {
	tokenData, err := s.extractAuthorizationData(r)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	_ = tokenData
}

func (s *server) cancelOrderEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	tokenData, err := s.extractAuthorizationData(r)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}
	_ = id
	_ = tokenData
}

func (s *server) getOrderStatusEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	tokenData, err := s.extractAuthorizationData(r)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	_ = id
	_ = tokenData
}

func (s *server) listOrdersEndpoint(w http.ResponseWriter, r *http.Request) {
	tokenData, err := s.extractAuthorizationData(r)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	_ = tokenData
}

func writeErrorResponse(w http.ResponseWriter, err error) {
	info := errorInfo{Code: errorCodeUnknown, Message: err.Error()}
	switch errors.Cause(err) {
	case app.ErrUserOrderNotExists:
		info.Code = errorCodeOrderNotFound
		w.WriteHeader(http.StatusNotFound)
	case errForbidden:
		w.WriteHeader(http.StatusForbidden)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	js, _ := json.Marshal(info)
	_, _ = w.Write(js)
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
