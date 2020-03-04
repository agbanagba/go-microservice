package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/agbanagba/go-microservice/data"

	"github.com/gorilla/mux"
)

// Products ...
type Products struct {
	l *log.Logger
}

// KeyProduct is a key for the product in the request context
type KeyProduct struct{}

// GenericError is a generic error message returned by the server
type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a list of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}

// NewProduct ...
func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

// AddProduct ...
func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")
	product := r.Context().Value(KeyProduct{}).(*data.Product)

	p.l.Printf("Product %#v", product)
	data.AddProduct(product)

}

// GetProducts ...
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")
	listProducts := data.GetProducts()

	// serialize data to JSON
	err := listProducts.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

// UpdateProducts ...
func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	p.l.Println("Handle PUT Products")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
	}

	product := r.Context().Value(KeyProduct{}).(*data.Product)

	err = data.UpdateProduct(id, product)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	// write a no content success header
	rw.WriteHeader(http.StatusNoContent)
}
