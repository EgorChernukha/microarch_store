package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"

	"store/pkg/common/infrastructure/amqp"
	"store/pkg/common/infrastructure/jwt"
	commonmysql "store/pkg/common/infrastructure/mysql"
	"store/pkg/common/infrastructure/prometheus"
	transportcommon "store/pkg/common/infrastructure/transport"

	appintegrationevent "store/pkg/billing/app/integrationevent"
	"store/pkg/billing/infrastructure/integrationevent"
	"store/pkg/billing/infrastructure/mysql"
	"store/pkg/billing/infrastructure/transport"
)

const (
	appID = "billing"
)

func main() {
	cnf, err := parseEnv()
	if err != nil {
		log.Fatal(err)
	}

	connector := commonmysql.NewConnector()
	err = connector.MigrateUp(cnf.dsn(), cnf.MigrationsDir)
	if err != nil {
		log.Fatal(err)
	}
	err = connector.Open(cnf.dsn(), commonmysql.Config{
		MaxConnections:     cnf.DBMaxConn,
		ConnectionLifetime: time.Duration(cnf.DBConnectionLifetime) * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	metricsHandler, err := prometheus.NewMetricsHandler(transportcommon.NewEndpointLabelCollector())
	if err != nil {
		log.Fatal(err)
	}

	amqpConnection := amqp.NewAMQPConnection(&amqp.Config{Host: cnf.AMQPHost, User: cnf.AMQPUser, Password: cnf.AMQPPassword})
	integrationEventTransport := integrationevent.NewIntegrationEventsTransport(false)
	amqpConnection.AddChannel(integrationEventTransport)
	eventHandler := integrationevent.NewIntegrationEventHandler([]appintegrationevent.Handler{appintegrationevent.NewHandler()})
	integrationEventTransport.SetHandler(eventHandler)

	err = amqpConnection.Start()
	if err != nil {
		log.Fatal(err)
	}
	// noinspection GoUnhandledErrorResult
	defer amqpConnection.Stop()

	srv := createServer(connector.Client(), metricsHandler, cnf)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-done
	if err := srv.Shutdown(context.TODO()); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}

func createServer(client commonmysql.Client, metricsHandler prometheus.MetricsHandler, cnf *config) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/health", healthEndpoint).Methods(http.MethodGet)
	metricsHandler.AddMetricsHandler(router, "/monitoring")
	metricsHandler.AddCommonMetricsMiddleware(router)
	tokenParser := jwt.NewTokenParser(cnf.JWTSecret)

	userAccountQueryService := mysql.NewUserAccountQueryService(client)

	server := transport.NewServer(router, tokenParser, userAccountQueryService)
	server.Start()

	return &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
}

func healthEndpoint(w http.ResponseWriter, r *http.Request) {
	json := simplejson.New()
	json.Set("status", "OK")

	payload, err := json.MarshalJSON()
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(payload)
}
