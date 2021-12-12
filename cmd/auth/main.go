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
	"store/pkg/common/infrastructure/prometheus"
)

const (
	appID = "auth"
)

func main() {
	cnf, err := parseEnv()
	if err != nil {
		log.Fatal(err)
	}

	connector := mysql.NewConnector()
	err = connector.MigrateUp(cnf.dsn(), cnf.MigrationsDir)
	if err != nil {
		log.Fatal(err)
	}
	err = connector.Open(cnf.dsn(), mysql.Config{
		MaxConnections:     cnf.DBMaxConn,
		ConnectionLifetime: time.Duration(cnf.DBConnectionLifetime) * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	metricsHandler, err := prometheus.NewMetricsHandler(transport.NewEndpointLabelCollector())
	if err != nil {
		log.Fatal(err)
	}

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

func createServer(client mysql.Client, metricsHandler prometheus.MetricsHandler, cnf *config) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/health", healthEndpoint).Methods(http.MethodGet)
	metricsHandler.AddMetricsHandler(router, "/metrics")
	metricsHandler.AddCommonMetricsMiddleware(router)

	userRepository := mysql.NewUserRepository(client)
	userService := app.NewUserService(userRepository, encoding.NewPasswordEncoder())
	sessionRepository := mysql.NewSessionRepository(client)
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
