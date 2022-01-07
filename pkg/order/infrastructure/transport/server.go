package transport

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"

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
var errBadRequest = errors.New("bad request")

type userOrderStatusData struct {
	Status int `json:"status"`
}

type Server interface {
	Start()
}

func NewServer(router *mux.Router, tokenParser jwt.TokenParser, userOrderQueryService app.UserOrderQueryService) Server {
	return &server{
		router:                router,
		tokenParser:           tokenParser,
		userOrderQueryService: userOrderQueryService,
	}
}

type server struct {
	router                *mux.Router
	tokenParser           jwt.TokenParser
	userOrderQueryService app.UserOrderQueryService
}

type errorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (s *server) Start() {
	s.router.Methods(http.MethodPost).Path(createOrderEndpoint).Handler(s.makeHandlerFunc(s.createOrderEndpoint))
	s.router.Methods(http.MethodPost).Path(cancelOrderEndpoint).Handler(s.makeHandlerFunc(s.cancelOrderEndpoint))
	s.router.Methods(http.MethodGet).Path(getOrderStatusEndpoint).Handler(s.makeHandlerFunc(s.getOrderStatusEndpoint))
	s.router.Methods(http.MethodGet).Path(listOrdersEndpoint).Handler(s.makeHandlerFunc(s.listOrdersEndpoint))
}

func (s *server) makeHandlerFunc(handler func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_ = request.ParseForm()
		fields := logrus.Fields{
			"method": request.Method,
			"host":   request.Host,
			"path":   request.URL.Path,
		}
		if request.URL.RawQuery != "" {
			fields["query"] = request.URL.RawQuery
		}
		if request.PostForm != nil {
			fields["post"] = request.PostForm
		}

		err := handler(writer, request)

		if err != nil {
			writeErrorResponse(writer, err)

			fields["err"] = err
			logrus.WithFields(fields).Error(err)
		} else {
			logrus.WithFields(fields).Info("call")
		}
	}
}

func (s *server) createOrderEndpoint(w http.ResponseWriter, r *http.Request) error {
	tokenData, err := s.extractAuthorizationData(r)
	if err != nil {
		return err
	}

	_ = tokenData
	return nil
}

func (s *server) cancelOrderEndpoint(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	tokenData, err := s.extractAuthorizationData(r)
	if err != nil {
		return err
	}
	_ = id
	_ = tokenData

	return nil
}

func (s *server) getOrderStatusEndpoint(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	tokenData, err := s.extractAuthorizationData(r)
	if err != nil {
		return err
	}

	orderID, err := uuid.FromString(id)
	if err != nil {
		return errBadRequest
	}

	userOrderData, err := s.userOrderQueryService.FindUserOrderByOrderID(orderID)
	if err != nil {
		return err
	}

	if userOrderData.UserID.String() != tokenData.UserID() {
		return errForbidden
	}

	writeResponse(w, userOrderStatusData{Status: userOrderData.Status})
	return nil
}

func (s *server) listOrdersEndpoint(w http.ResponseWriter, r *http.Request) error {
	tokenData, err := s.extractAuthorizationData(r)
	if err != nil {
		return err
	}

	_ = tokenData

	return nil
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
