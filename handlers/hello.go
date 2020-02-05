package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Hello ...
type Hello struct {
	l *log.Logger
}

// NewHello returns a new Hello handler with a logger
func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Hello World")

	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// all of this can be replaced with http Error
		http.Error(rw, "Oops", http.StatusBadRequest)
		// rw.WriteHeader(http.StatusBadRequest)
		// rw.Write([]byte("Oops"))
		// return
	}

	fmt.Fprintf(rw, "Welcome %s\n", d)
}
