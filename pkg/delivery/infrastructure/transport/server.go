package transport

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"

	"store/pkg/common/infrastructure/jwt"
	"store/pkg/delivery/app"
)

const PathPrefix = "/api/v1/"

const (
	specOrderDeliveryEndpoint = PathPrefix + "order/{id}"
)

const (
	errorCodeUnknown          = 0
	errorCodeDeliveryNotFound = 1
)

const authTokenHeader = "X-Auth-Token"

var errUnauthorized = errors.New("not authorized")
var errForbidden = errors.New("access denied")
var errBadRequest = errors.New("bad request")

type Server interface {
	Start()
}

func NewServer(router *mux.Router, tokenParser jwt.TokenParser, trUnitFactory app.TransactionalUnitFactory, orderDeliveryQueryService app.OrderDeliveryQueryService) Server {
	return &server{
		router:                    router,
		tokenParser:               tokenParser,
		trUnitFactory:             trUnitFactory,
		orderDeliveryQueryService: orderDeliveryQueryService,
	}
}

type server struct {
	router                    *mux.Router
	tokenParser               jwt.TokenParser
	trUnitFactory             app.TransactionalUnitFactory
	orderDeliveryQueryService app.OrderDeliveryQueryService
}

type errorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (s *server) Start() {
	s.router.Methods(http.MethodGet).Path(specOrderDeliveryEndpoint).Handler(s.makeHandlerFunc(s.getOrderDeliveryEndpoint))
}

func (s *server) executeInTransaction(f func(provider app.RepositoryProvider) error) (err error) {
	var trUnit app.TransactionalUnit
	trUnit, err = s.trUnitFactory.NewTransactionalUnit()
	if err != nil {
		return err
	}
	defer func() {
		err = trUnit.Complete(err)
	}()
	err = f(trUnit)
	return err
}

func (s *server) getOrderDeliveryEndpoint(w http.ResponseWriter, r *http.Request, _ app.OrderDeliveryService) error {
	vars := mux.Vars(r)
	id := vars["id"]

	orderID, err := uuid.FromString(id)
	if err != nil {
		return errBadRequest
	}

	positionData, err := s.orderDeliveryQueryService.FindByOrderID(orderID)
	if err != nil {
		return err
	}

	writeResponse(w, positionData)
	return nil
}

func (s *server) makeHandlerFunc(handler func(http.ResponseWriter, *http.Request, app.OrderDeliveryService) error) http.HandlerFunc {
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

		err := s.executeInTransaction(func(provider app.RepositoryProvider) error {
			service := app.NewOrderDeliveryService(provider.OrderDeliveryRepository())
			return handler(writer, request, service)
		})

		if err != nil {
			writeErrorResponse(writer, err)

			fields["err"] = err
			logrus.WithFields(fields).Error(err)
		} else {
			logrus.WithFields(fields).Info("call")
		}
	}
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
	case app.ErrOrderDeliveryNotFound:
		info.Code = errorCodeDeliveryNotFound
		w.WriteHeader(http.StatusNotFound)
	case app.ErrOrderDeliveryNotExists:
		info.Code = errorCodeDeliveryNotFound
		w.WriteHeader(http.StatusNotFound)
	case errUnauthorized:
		w.WriteHeader(http.StatusUnauthorized)
	case errBadRequest:
		w.WriteHeader(http.StatusBadRequest)
	case errForbidden:
		w.WriteHeader(http.StatusForbidden)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	js, _ := json.Marshal(info)
	_, _ = w.Write(js)
}
