package transport

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"store/pkg/common/infrastructure/jwt"
)

const (
	currentDeliveryEndpoint = PathPrefix + "delivery"
	specDeliveryEndpoint    = PathPrefix + "delivery/{id}"
)

const (
	errorCodeUnknown  = 0
	errorCodeNotFound = 1
)

const authTokenHeader = "X-Auth-Token"

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
	case errForbidden:
		w.WriteHeader(http.StatusForbidden)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	js, _ := json.Marshal(info)
	_, _ = w.Write(js)
}

var errForbidden = errors.New("access denied")
