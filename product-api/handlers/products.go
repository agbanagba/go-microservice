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

// getProductID returns the product id
func getProductID(r *http.Request) int {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// this should not happen if correct id is given
		panic(err)
	}
	return id
}

// AddProduct ...
func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")
	product := r.Context().Value(KeyProduct{}).(*data.Product)

	p.l.Printf("Product %#v", product)
	data.AddProduct(product)

}

// GetProducts returns a list of products
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {

	// swagger:route GET /products products listProducts
	// Returns a list of products
	//
	// 	Responses:
	//		200: productsResponse
	//		400: errorResponse

	p.l.Println("Handle GET Products")
	rw.Header().Add("Content-Type", "application/json")
	listProducts := data.GetProducts()

	// serialize data to JSON
	err := data.ToJSON(listProducts, rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

// GetProduct retuns a single product with the an id
func (p *Products) GetProduct(rw http.ResponseWriter, r *http.Request) {

	// swagger:route GET /products/{id} products listSingleProduct
	// Retuns a single product from the database
	//
	// 	Responses:
	//		200: productsResponse
	//		400: errorResponse

	rw.Header().Add("Content-Type", "application/json")
	id := getProductID(r)

	p.l.Println("[DEBUG] get record id", id)
	prod, err := data.GetProductByID(id)

	switch err {
	case nil:
	case data.ErrProductNotFound:
		p.l.Println("[ERROR] fetching product", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		p.l.Println("[ERROR] fetching product", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.ToJSON(prod, rw)
	if err != nil {
		p.l.Println("[ERROR] serializing product", err)
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

// DeleteProduct removes a product from product list
func (p *Products) DeleteProduct(rw http.ResponseWriter, r *http.Request) {

	// swagger:route DELETE /products/{id} products deleteProduct
	//
	// Deletes a product
	//
	// Responses:
	// 		201: noContentResponse
	// 404: errorResponse
	// 501: errorResponse

	rw.Header().Add("Content-Type", "application/json")
	id := getProductID(r)

	p.l.Println("[DEBUG] deleting product with id", id)
	err := data.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		p.l.Println("[ERROR] deleting record id does not exist.")

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	if err != nil {
		p.l.Println("[ERROR] deleting record", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
