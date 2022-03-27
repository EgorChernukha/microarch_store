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
	"store/pkg/stock/app"
)

const PathPrefix = "/api/v1/"
const PathPrefixInternal = "/internal/api/v1/"

const (
	listPositionsEndpoint   = PathPrefix + "position/list"
	positionEndpoint        = PathPrefix + "position"
	specPositionEndpoint    = positionEndpoint + "/{id}"
	addPositionEndpoint     = positionEndpoint + "/add"
	reservePositionEndpoint = PathPrefixInternal + "position/reserve"
	topUpPositionEndpoint   = specPositionEndpoint + "/topup"
)

const (
	errorCodeUnknown          = 0
	errorCodePositionNotFound = 1
)

const authTokenHeader = "X-Auth-Token"

var errUnauthorized = errors.New("not authorized")
var errForbidden = errors.New("access denied")
var errBadRequest = errors.New("bad request")

type addPositionRequest struct {
	Title string `json:"title"`
	Count int    `json:"count"`
}

type addPositionResponse struct {
	ID string `json:"id"`
}

type topUpPositionRequest struct {
	Count int `json:"count"`
}

type reservePositionRequestPosition struct {
	PositionID string `json:"position_id"`
	OrderID    string `json:"order_id"`
	Count      int    `json:"count"`
}

type reservePositionRequest struct {
	Positions []reservePositionRequestPosition `json:"positions"`
}

type Server interface {
	Start()
}

func NewServer(router *mux.Router, tokenParser jwt.TokenParser, trUnitFactory app.TransactionalUnitFactory, positionQueryService app.PositionQueryService) Server {
	return &server{
		router:               router,
		tokenParser:          tokenParser,
		trUnitFactory:        trUnitFactory,
		positionQueryService: positionQueryService,
	}
}

type server struct {
	router               *mux.Router
	tokenParser          jwt.TokenParser
	trUnitFactory        app.TransactionalUnitFactory
	positionQueryService app.PositionQueryService
}

type errorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (s *server) Start() {
	s.router.Methods(http.MethodGet).Path(specPositionEndpoint).Handler(s.makeHandlerFunc(s.getStockPositionEndpoint))
	s.router.Methods(http.MethodGet).Path(listPositionsEndpoint).Handler(s.makeHandlerFunc(s.listStockPositionsEndpoint))
	s.router.Methods(http.MethodPost).Path(addPositionEndpoint).Handler(s.makeHandlerFunc(s.addStockPositionEndpoint))
	s.router.Methods(http.MethodPost).Path(topUpPositionEndpoint).Handler(s.makeHandlerFunc(s.topUpStockPositionEndpoint))
	s.router.Methods(http.MethodPost).Path(reservePositionEndpoint).Handler(s.makeHandlerFunc(s.reserveStockPositionEndpoint))
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

func (s *server) makeHandlerFunc(handler func(http.ResponseWriter, *http.Request, app.PositionService) error) http.HandlerFunc {
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
			service := app.NewPositionService(provider.PositionRepository(), provider.OrderPositionRepository())
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

func (s *server) getStockPositionEndpoint(w http.ResponseWriter, r *http.Request, _ app.PositionService) error {
	vars := mux.Vars(r)
	id := vars["id"]

	positionID, err := uuid.FromString(id)
	if err != nil {
		return errBadRequest
	}

	positionData, err := s.positionQueryService.FindPositionByID(positionID)
	if err != nil {
		return err
	}

	writeResponse(w, positionData)
	return nil
}

func (s *server) listStockPositionsEndpoint(w http.ResponseWriter, r *http.Request, _ app.PositionService) error {
	positionData, err := s.positionQueryService.ListPositions()
	if err != nil {
		return err
	}

	writeResponse(w, positionData)
	return nil
}

func (s *server) addStockPositionEndpoint(w http.ResponseWriter, r *http.Request, positionService app.PositionService) error {
	var requestData addPositionRequest
	bytesBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(bytesBody, &requestData); err != nil {
		return err
	}

	positionID, err := positionService.AddPosition(requestData.Title, requestData.Count)
	if err != nil {
		return err
	}

	response := addPositionResponse{ID: positionID.String()}

	writeResponse(w, response)
	return nil
}

func (s *server) topUpStockPositionEndpoint(w http.ResponseWriter, r *http.Request, positionService app.PositionService) error {
	vars := mux.Vars(r)
	id := vars["id"]

	positionID, err := uuid.FromString(id)
	if err != nil {
		return errBadRequest
	}

	var requestData topUpPositionRequest
	bytesBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(bytesBody, &requestData); err != nil {
		return err
	}

	if err = positionService.TopUpPosition(positionID, requestData.Count); err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func (s *server) reserveStockPositionEndpoint(w http.ResponseWriter, r *http.Request, positionService app.PositionService) error {
	var requestData reservePositionRequest
	bytesBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(bytesBody, &requestData); err != nil {
		return err
	}

	var positionsInput []app.ReservePositionInputItem
	for _, requestItem := range requestData.Positions {
		orderID, err := uuid.FromString(requestItem.OrderID)
		if err != nil {
			return err
		}
		positionID, err := uuid.FromString(requestItem.PositionID)
		if err != nil {
			return err
		}

		positionsInput = append(positionsInput, app.ReservePositionInputItem{
			OrderID:    orderID,
			PositionID: positionID,
			Count:      requestItem.Count,
		})
	}

	if err = positionService.ReservePosition(app.ReservePositionInput{Items: positionsInput}); err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
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
	case app.ErrOrderPositionNotExists:
		info.Code = errorCodePositionNotFound
		w.WriteHeader(http.StatusNotFound)
	case app.ErrPositionNotExists:
		info.Code = errorCodePositionNotFound
		w.WriteHeader(http.StatusNotFound)
	case app.ErrOrderPositionNotFound:
		info.Code = errorCodePositionNotFound
		w.WriteHeader(http.StatusNotFound)
	case app.ErrPositionNotFound:
		info.Code = errorCodePositionNotFound
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
