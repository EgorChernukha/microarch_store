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

	"store/pkg/common/app/streams"
	commonintegrationevent "store/pkg/common/infrastructure/integrationevent"
	"store/pkg/common/infrastructure/jwt"
	commonmysql "store/pkg/common/infrastructure/mysql"
	"store/pkg/common/infrastructure/prometheus"
	commonstreams "store/pkg/common/infrastructure/streams"
	transportcommon "store/pkg/common/infrastructure/transport"

	"store/pkg/stock/app"
	"store/pkg/stock/infrastructure/integrationevent"
	"store/pkg/stock/infrastructure/mysql"
	"store/pkg/stock/infrastructure/transport"
)

const (
	appID = "stock"
)

func main() {
	cnf, err := parseEnv()
	if err != nil {
		log.Fatal(err)
	}

	connector := commonmysql.NewConnector()
	defer connector.Close()
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

	streamsEnvironment, err := initStreamsEnvironment(cnf)
	if err != nil {
		log.Fatal(err)
	}

	metricsHandler, err := prometheus.NewMetricsHandler(transportcommon.NewEndpointLabelCollector())
	if err != nil {
		log.Fatal(err)
	}

	srv := createServer(connector, streamsEnvironment, metricsHandler, cnf)

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

func initStreamsEnvironment(cfg *config) (streams.Environment, error) {
	return commonstreams.NewEnvironment(appID,
		streams.Config{
			Host:     cfg.AMQPHost,
			Port:     cfg.AMQPPort,
			User:     cfg.AMQPUser,
			Password: cfg.AMQPPassword,
		})
}

func createServer(connector commonmysql.Connector, streamsEnvironment streams.Environment, metricsHandler prometheus.MetricsHandler, cnf *config) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/health", healthEndpoint).Methods(http.MethodGet)
	metricsHandler.AddMetricsHandler(router, "/monitoring")
	metricsHandler.AddCommonMetricsMiddleware(router)

	tokenParser := jwt.NewTokenParser(cnf.JWTSecret)

	trUnitFactory := mysql.NewTransactionalUnitFactory(connector.Client())
	eventHandler := app.NewEventHandler(trUnitFactory, integrationevent.NewEventParser())

	if err := commonintegrationevent.StartEventConsumer(streamsEnvironment, eventHandler); err != nil {
		log.Fatal(err)
	}

	positionQueryService := mysql.NewPositionQueryService(connector.Client())

	server := transport.NewServer(router, tokenParser, trUnitFactory, positionQueryService)
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
