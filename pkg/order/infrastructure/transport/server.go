package transport

import (
	"encoding/json"
	"io/ioutil"
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
	createOrderEndpoint  = PathPrefix + "order"
	ordersEndpoint       = PathPrefix + "orders"
	specificOderEndpoint = PathPrefix + "order/{id}"
)

const (
	errorCodeUnknown                     = 0
	errorCodeOrderNotFound               = 1
	errorCodePaymentFailed               = 2
	errorCodeReserveOrderDeliveryFailed  = 3
	errorCodeReserveOrderPositionsFailed = 4
)

var errUnauthorized = errors.New("not authorized")
var errForbidden = errors.New("access denied")
var errBadRequest = errors.New("bad request")

type Server interface {
	Start()
}

func NewServer(router *mux.Router, tokenParser jwt.TokenParser, userOrderService app.UserOrderService, userOrderQueryService app.UserOrderQueryService) Server {
	return &server{
		router:                router,
		tokenParser:           tokenParser,
		userOrderService:      userOrderService,
		userOrderQueryService: userOrderQueryService,
	}
}

type server struct {
	router                *mux.Router
	tokenParser           jwt.TokenParser
	userOrderService      app.UserOrderService
	userOrderQueryService app.UserOrderQueryService
}

type errorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type createOrderRequest struct {
	Price      float64 `json:"price"`
	PositionID string  `json:"positionID"`
	Count      int     `json:"count"`
}

type createOrderResponse struct {
	ID string `json:"id"`
}

func (s *server) Start() {
	s.router.Methods(http.MethodPost).Path(createOrderEndpoint).Handler(s.makeHandlerFunc(s.createOrderEndpoint))
	s.router.Methods(http.MethodGet).Path(specificOderEndpoint).Handler(s.makeHandlerFunc(s.getOrderEndpoint))
	s.router.Methods(http.MethodGet).Path(ordersEndpoint).Handler(s.makeHandlerFunc(s.listOrdersEndpoint))
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

	var requestData createOrderRequest
	bytesBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(bytesBody, &requestData); err != nil {
		return err
	}

	userID, err := uuid.FromString(tokenData.UserID())
	if err != nil {
		return err
	}

	positionID, err := uuid.FromString(requestData.PositionID)
	if err != nil {
		return err
	}

	orderID, err := s.userOrderService.Create(app.UserID(userID), requestData.Price, app.PositionID(positionID), requestData.Count)
	if err != nil {
		return err
	}
	response := createOrderResponse{ID: uuid.UUID(orderID).String()}
	writeResponse(w, response)
	return nil
}

func (s *server) getOrderEndpoint(w http.ResponseWriter, r *http.Request) error {
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

	writeResponse(w, userOrderData)
	return nil
}

func (s *server) listOrdersEndpoint(w http.ResponseWriter, r *http.Request) error {
	tokenData, err := s.extractAuthorizationData(r)
	if err != nil {
		return err
	}

	userID, err := uuid.FromString(tokenData.UserID())
	if err != nil {
		return err
	}

	userOrdersData, err := s.userOrderQueryService.ListUserOrdersByUserIDs(userID)
	if err != nil {
		return err
	}

	writeResponse(w, userOrdersData)
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
	case app.ErrPaymentFailed:
		info.Code = errorCodePaymentFailed
		w.WriteHeader(http.StatusBadRequest)
	case app.ErrReserveOrderDeliveryFailed:
		info.Code = errorCodeReserveOrderDeliveryFailed
		w.WriteHeader(http.StatusBadRequest)
	case app.ErrReserveOrderPositionsFailed:
		info.Code = errorCodeReserveOrderPositionsFailed
		w.WriteHeader(http.StatusBadRequest)
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
