package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/agbanagba/go-microservice/data"
)

// MiddlewareProductValidation ...
func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")

		product := &data.Product{}
		err := data.FromJSON(product, r.Body)
		if err != nil {
			p.l.Println("ERROR deserializing product", err)

			rw.WriteHeader(http.StatusBadRequest)

			// TODO Use this here instead of http.Error
			// data.ToJSON(&GenericError{Message: err.Error()}, rw)

			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		// validate the product right here
		err = product.Validate()
		if err != nil {
			p.l.Println("ERROR validating product", err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating product: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
