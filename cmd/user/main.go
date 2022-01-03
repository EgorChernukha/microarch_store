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

	"store/pkg/common/infrastructure/prometheus"
	"store/pkg/user/app"
	"store/pkg/user/domain"
	"store/pkg/user/infrastructure/jwt"
	"store/pkg/user/infrastructure/mysql"
	"store/pkg/user/infrastructure/transport"
)

const (
	appID = "user"
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
	metricsHandler.AddMetricsHandler(router, "/monitoring")
	metricsHandler.AddCommonMetricsMiddleware(router)

	userRepository := mysql.NewUserRepository(client)
	userDomainService := domain.NewUserService(userRepository)
	userService := app.NewUserService(userDomainService)
	userQueryService := mysql.NewUserQueryService(client)
	tokenParser := jwt.NewTokenParser(cnf.JWTSecret)

	server := transport.NewServer(router, tokenParser, userService, userQueryService)
	server.Start()

	return &http.Server{
		Addr:    ":8000",
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
