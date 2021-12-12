package transport

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"

	"store/pkg/auth/app"
	"store/pkg/auth/infrastructure/jwt"
)

const PathPrefix = "/api/v1/"

const (
	registerUserEndpoint = PathPrefix + "register"
	authEndpoint         = PathPrefix + "auth"
	loginEndpoint        = PathPrefix + "login"
	logoutEndpoint       = PathPrefix + "logout"
)

const (
	errorCodeUnknown         = 0
	errorCodeUserNotFound    = 1
	errUserAlreadyExists     = 2
	errInvalidLogin          = 3
	errorCodeInvalidPassword = 4
)

const sessionCookieName = "session_id"
const sessionLifetime = time.Minute * 30
const authTokenHeader = "X-Auth-Token"

var errUnauthorized = errors.New("not authorized")

type Server interface {
	Start()
}

func NewServer(router *mux.Router, authService app.UserService, sessionRepository app.SessionRepository, tokenGenerator jwt.TokenGenerator) Server {
	return &server{
		router:            router,
		userService:       authService,
		sessionRepository: sessionRepository,
		tokenGenerator:    tokenGenerator,
	}
}

type userAuthData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type createdUserInfo struct {
	UserID string `json:"id"`
}

type errorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type server struct {
	router            *mux.Router
	userService       app.UserService
	sessionRepository app.SessionRepository
	tokenGenerator    jwt.TokenGenerator
}

func (s *server) Start() {
	s.router.Methods(http.MethodPost).Path(registerUserEndpoint).Handler(s.makeHandlerFunc(s.registerUserEndpoint))
	s.router.Methods(http.MethodPost).Path(loginEndpoint).Handler(s.makeHandlerFunc(s.loginEndpoint))
	s.router.Methods(http.MethodPost).Path(logoutEndpoint).Handler(s.makeHandlerFunc(s.logoutEndpoint))
	s.router.Path(authEndpoint).Handler(s.makeHandlerFunc(s.authEndpoint))
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

func (s *server) registerUserEndpoint(w http.ResponseWriter, r *http.Request) error {
	var info userAuthData
	bytesBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(bytesBody, &info); err != nil {
		return err
	}

	userID, err := s.userService.AddUser(info.Login, info.Password)
	if err != nil {
		return err
	}
	writeResponse(w, createdUserInfo{UserID: uuid.UUID(userID).String()})
	return nil
}

func (s *server) loginEndpoint(w http.ResponseWriter, r *http.Request) error {
	var info userAuthData
	bytesBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(bytesBody, &info); err != nil {
		return err
	}

	user, err := s.userService.FindUserByLoginAndPassword(info.Login, info.Password)
	if err != nil {
		return err
	}
	session := app.Session{
		ID:        app.SessionID(uuid.NewV1()),
		UserID:    user.ID,
		ValidTill: time.Now().Add(sessionLifetime),
	}
	err = s.sessionRepository.Store(&session)
	if err != nil {
		return err
	}

	setSessionCookie(w, &session.ID)
	w.WriteHeader(http.StatusOK)
	return nil
}

func (s *server) logoutEndpoint(w http.ResponseWriter, r *http.Request) error {
	if sessionID, err := getSessionIDFromRequest(r); err == nil {
		err = s.sessionRepository.Remove(sessionID)
		if err != nil {
			return err
		}
	}

	setSessionCookie(w, nil)
	w.WriteHeader(http.StatusOK)
	return nil
}

func (s *server) authEndpoint(w http.ResponseWriter, r *http.Request) error {
	sessionID, err := getSessionIDFromRequest(r)
	if err != nil {
		return errUnauthorized
	}
	session, err := s.sessionRepository.FindOneByID(sessionID)
	if err != nil {
		if errors.Cause(err) == app.ErrSessionNotFound {
			return errUnauthorized
		}
		return err
	}
	user, err := s.userService.FindUserByID(session.UserID)
	if err != nil {
		return err
	}
	session.ValidTill = time.Now().Add(sessionLifetime)
	_ = s.sessionRepository.Store(session)

	token, err := s.tokenGenerator.GenerateToken(uuid.UUID(user.ID).String(), string(user.Login))
	if err != nil {
		return err
	}

	w.Header().Set(authTokenHeader, token)
	w.WriteHeader(http.StatusOK)
	return nil
}

func getSessionIDFromRequest(r *http.Request) (app.SessionID, error) {
	sessionIDCookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return app.SessionID{}, err
	}

	sessionID, err := uuid.FromString(sessionIDCookie.Value)

	return app.SessionID(sessionID), nil
}

func setSessionCookie(w http.ResponseWriter, sessionID *app.SessionID) {
	c := &http.Cookie{
		Name:     sessionCookieName,
		Path:     "/",
		HttpOnly: true,
	}
	if sessionID != nil {
		c.Value = uuid.UUID(*sessionID).String()
	} else {
		// delete cookie
		c.MaxAge = -1
	}

	http.SetCookie(w, c)
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
	case app.ErrUserNotFound:
		info.Code = errorCodeUserNotFound
		w.WriteHeader(http.StatusNotFound)
	case app.ErrUserAlreadyExists:
		info.Code = errUserAlreadyExists
		w.WriteHeader(http.StatusBadRequest)
	case app.ErrInvalidLogin:
		info.Code = errInvalidLogin
		w.WriteHeader(http.StatusBadRequest)
	case app.ErrInvalidPassword:
		info.Code = errorCodeInvalidPassword
		w.WriteHeader(http.StatusBadRequest)
	case errUnauthorized:
		w.WriteHeader(http.StatusUnauthorized)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	js, _ := json.Marshal(info)
	_, _ = w.Write(js)
}
