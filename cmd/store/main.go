package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/health", healthEndpoint).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

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
