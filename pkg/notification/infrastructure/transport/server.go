package transport

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"

	"store/pkg/common/infrastructure/jwt"

	"store/pkg/notification/app"
)

const PathPrefix = "/api/v1/"
const authTokenHeader = "X-Auth-Token"

const (
	listUserNotificationsEndpoint = PathPrefix + "notification/list"
)

const (
	errorCodeUnknown = 0
)

var errUnauthorized = errors.New("not authorized")
var errForbidden = errors.New("access denied")

type Server interface {
	Start()
}

func NewServer(router *mux.Router, tokenParser jwt.TokenParser, userNotificationQueryService app.UserNotificationQueryService) Server {
	return &server{
		router:                       router,
		tokenParser:                  tokenParser,
		userNotificationQueryService: userNotificationQueryService,
	}
}

type server struct {
	router                       *mux.Router
	tokenParser                  jwt.TokenParser
	userNotificationQueryService app.UserNotificationQueryService
}

func (s *server) Start() {
	s.router.Methods(http.MethodGet).Path(listUserNotificationsEndpoint).Handler(s.makeHandlerFunc(s.listUserNotificationsEndpoint))
}

type errorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
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

func (s *server) listUserNotificationsEndpoint(w http.ResponseWriter, r *http.Request) error {
	tokenData, err := s.extractAuthorizationData(r)
	if err != nil {
		return err
	}

	userID, err := uuid.FromString(tokenData.UserID())
	if err != nil {
		return err
	}

	userOrdersData, err := s.userNotificationQueryService.ListUserNotifications(userID)
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
