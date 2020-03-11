// Package handlers classification of Product API
//
// Documentation of Product API
//
//	Schemes: http
// 	BasePath: /
//	Version: 1.0.0
//
// 	Consumes:
//	- application/json
//
// 	Produces:
// 	- application/json
// swagger:meta
package handlers

import "github.com/agbanagba/go-microservice/data"

// Types here are used for documentation purposes with go-swagger
// and not used in any handler

// A list of products returns in the response
// swagger:response productsResponse
type productsResponseWrapper struct {
	// All products in the system
	// in: body
	Body []data.Product
}

// A single product
// swagger:response productResponse
type productResponseWrapper struct {
	// Newly created product
	// in: body
	Body data.Product
}

// Generic message returned as a string
// swagger:response errorResponse
type errorResponseWrapper struct {
	// Error description
	// in: body
	Body GenericError
}

// Validation error as an array of strings
// swagger:response errorValidation
type errorValidationWrapper struct {
	// Collection of errors
	// in: body
	Body ValidationError
}

// No content returned
// swagger:response noContentResponse
type noContentResponseWrapper struct{}

// swagger:parameters updateProduct createProduct
type productParamsWrapper struct {
	// Product data structure to Update or Create.
	// Note: the id field is ignored by update and create operations
	// in: body
	// required: true
	Body data.Product
}

// swagger:parameters listSingleProduct deleteProduct
type productIDParamsWrapper struct {
	// The id of the product for which the operation relates
	// in: path
	// required: true
	ID int `json:"id"`
}
