package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	ghandlers "github.com/gorilla/handlers"

	"github.com/agbanagba/go-microservice/images-api/files"
	"github.com/agbanagba/go-microservice/images-api/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/env"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")
var logLevel = env.String("LOG_LEVEL", false, "debug", "Log level output [debug, info, trace]")
var basePath = env.String("BASE_PATH", false, "./imagestore", "Path to save images")

func main() {
	env.Parse()

	l := hclog.New(&hclog.LoggerOptions{
		Name:  "images-api",
		Level: hclog.LevelFromString(*logLevel),
	})

	sl := l.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true})

	// local storage class. max filesize 5MB
	stor, err := files.NewLocal(*basePath, 1024*1000*5)
	if err != nil {
		l.Error("Unable to create storage", "error", err, stor)
		os.Exit(1)
	}

	filehandler := handlers.NewFiles(stor, l)

	// serve mux for registering handlers
	sm := mux.NewRouter()

	ph := sm.Methods(http.MethodPost).Subrouter()
	ph.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", filehandler.UploadREST)
	// ph.HandleFunc()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.Handle(
		"/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}",
		http.StripPrefix("/images/", http.FileServer(http.Dir(*basePath))),
	)

	// CORS handler allowing all origins to access product api
	corsHandler := ghandlers.CORS(ghandlers.AllowedOrigins([]string{"*"}))

	s := http.Server{
		Addr:         *bindAddress,
		Handler:      corsHandler(sm),
		ErrorLog:     sl,
		IdleTimeout:  120 * time.Second,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	go func() {
		l.Info("Starting server", "bind_address", *bindAddress)

		err := s.ListenAndServe()
		if err != nil {
			l.Error("Unable to start server", "error", err)
			os.Exit(1)
		}
	}()

	// trap a sigterm and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	l.Info("Shutting down server with", "signal", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
