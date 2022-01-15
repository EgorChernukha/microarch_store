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

	"store/pkg/auth/app"
	"store/pkg/auth/infrastructure/encoding"
	"store/pkg/auth/infrastructure/jwt"
	"store/pkg/auth/infrastructure/mysql"
	"store/pkg/auth/infrastructure/transport"
	"store/pkg/common/app/streams"
	commonmysql "store/pkg/common/infrastructure/mysql"
	"store/pkg/common/infrastructure/prometheus"
	"store/pkg/common/infrastructure/storedevent"
	commonstreams "store/pkg/common/infrastructure/streams"
	commontransport "store/pkg/common/infrastructure/transport"
)

const (
	appID = "auth"
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

	metricsHandler, err := prometheus.NewMetricsHandler(commontransport.NewEndpointLabelCollector())
	if err != nil {
		log.Fatal(err)
	}

	streamsEnvironment, err := initStreamsEnvironment(cnf)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	srv := createServer(ctx, connector, streamsEnvironment, metricsHandler, cnf)

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

func createServer(ctx context.Context, connector commonmysql.Connector, streamsEnvironment streams.Environment, metricsHandler prometheus.MetricsHandler, cnf *config) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/health", healthEndpoint).Methods(http.MethodGet)
	metricsHandler.AddMetricsHandler(router, "/monitoring")
	metricsHandler.AddCommonMetricsMiddleware(router)

	eventStore, err := storedevent.NewEventSender(ctx, mysql.NewEventStore(connector.Client()), streamsEnvironment)
	if err != nil {
		log.Fatal(err)
	}
	trUnitFactory := mysql.NewTransactionalUnitFactory(connector.Client())
	userRepository := mysql.NewUserRepository(connector.Client())
	userService := app.NewUserService(trUnitFactory, userRepository, eventStore, encoding.NewPasswordEncoder())
	sessionRepository := mysql.NewSessionRepository(connector.Client())
	tokenGenerator := jwt.NewTokenGenerator(cnf.JWTSecret)

	server := transport.NewServer(router, userService, sessionRepository, tokenGenerator)
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
