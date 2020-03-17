package main

import (
	"net/http"
	"os"

	"github.com/agbanagba/go-microservice/images-api/files"
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
		l.Error("Unable to create storage", "error", err)
		os.Exit(1)
	}

	// serve mux for registering handlers
	sm := mux.NewRouter()

	s := http.Server{
		Addr:     *bindAddress,
		Handler:  sm,
		ErrorLog: sl,
	}
}
