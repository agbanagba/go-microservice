package main

import (
	"context"
	"golearn/microservices/product/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	l := log.New(os.Stdout, "product-api:", log.LstdFlags)
	productHandler := handlers.NewProduct(l)

	// A new servemux can be created and used instead of the default servemux
	servemux := http.NewServeMux()
	servemux.Handle("/", productHandler)

	// Proper server to handle timeouts and other things like TLS configs, keep-alives etc
	s := &http.Server{
		Addr:         ":9090",
		Handler:      servemux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// handling listen and serve so it won't block
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)

	// Try to gracefully shut down any connection and after 30s kill everything
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
