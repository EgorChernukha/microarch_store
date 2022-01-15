package transport

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"store/pkg/billing/domain"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"

	"store/pkg/billing/app"
	"store/pkg/common/infrastructure/jwt"
)

const PathPrefix = "/api/v1/"
const PathPrefixInternal = "/internal/api/v1/"
const authTokenHeader = "X-Auth-Token"

const (
	accountEndpoint = PathPrefix + "account"
	paymentEndpoint = PathPrefixInternal + "payment"
)

const (
	errorCodeUnknown             = 0
	errorCodeUserAccountNotFound = 1
	errorNotEnoughBalance        = 2
	errorInvalidAmount           = 3
)

var errUnauthorized = errors.New("not authorized")
var errForbidden = errors.New("access denied")

type Server interface {
	Start()
}

func NewServer(router *mux.Router, tokenParser jwt.TokenParser, userAccountService app.UserAccountService, userAccountQueryService app.UserAccountQueryService) Server {
	return &server{
		router:                  router,
		tokenParser:             tokenParser,
		userAccountService:      userAccountService,
		userAccountQueryService: userAccountQueryService,
	}
}

type server struct {
	router                  *mux.Router
	tokenParser             jwt.TokenParser
	userAccountService      app.UserAccountService
	userAccountQueryService app.UserAccountQueryService
}

type errorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type topUpAccountInfo struct {
	Amount float64 `json:"amount"`
}

type paymentInfo struct {
	UserID string  `json:"userId"`
	Amount float64 `json:"amount"`
}

func (s *server) Start() {
	s.router.Methods(http.MethodGet).Path(accountEndpoint).Handler(s.makeHandlerFunc(s.getUserAccountEndpoint))
	s.router.Methods(http.MethodPost).Path(accountEndpoint).Handler(s.makeHandlerFunc(s.topUpAccountEndpoint))
	s.router.Methods(http.MethodPost).Path(paymentEndpoint).Handler(s.makeHandlerFunc(s.processPaymentEndpoint))
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

func (s *server) getUserAccountEndpoint(w http.ResponseWriter, r *http.Request) error {
	tokenData, err := s.extractAuthorizationData(r)
	if err != nil {
		return err
	}

	userID, err := uuid.FromString(tokenData.UserID())
	if err != nil {
		return err
	}

	userAccountData, err := s.userAccountQueryService.FindUserAccountByUserID(userID)
	if err != nil {
		return err
	}

	writeResponse(w, userAccountData)
	return nil
}

func (s *server) topUpAccountEndpoint(w http.ResponseWriter, r *http.Request) error {
	tokenData, err := s.extractAuthorizationData(r)
	if err != nil {
		return err
	}

	userID, err := uuid.FromString(tokenData.UserID())
	if err != nil {
		return err
	}

	var info topUpAccountInfo
	bytesBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(bytesBody, &info); err != nil {
		return err
	}

	if err = s.userAccountService.TopUpAccount(userID, info.Amount); err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func (s *server) processPaymentEndpoint(w http.ResponseWriter, r *http.Request) error {
	var info paymentInfo
	bytesBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(bytesBody, &info); err != nil {
		return err
	}
	userID, err := uuid.FromString(info.UserID)
	if err != nil {
		return err
	}

	if err = s.userAccountService.ProcessPayment(userID, info.Amount); err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
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
	case app.ErrUserAccountNotExists:
		info.Code = errorCodeUserAccountNotFound
		w.WriteHeader(http.StatusNotFound)
	case errForbidden:
		w.WriteHeader(http.StatusForbidden)
	case domain.ErrNotEnoughBalance:
		info.Code = errorNotEnoughBalance
		w.WriteHeader(http.StatusBadRequest)
	case domain.ErrInvalidAmount:
		info.Code = errorInvalidAmount
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
