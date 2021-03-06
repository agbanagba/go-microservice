// Code generated by go-swagger; DO NOT EDIT.

package products

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/agbanagba/go-microservice/product-api/client/models"
)

// ListProductsReader is a Reader for the ListProducts structure.
type ListProductsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListProductsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListProductsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewListProductsBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewListProductsOK creates a ListProductsOK with default headers values
func NewListProductsOK() *ListProductsOK {
	return &ListProductsOK{}
}

/*ListProductsOK handles this case with default header values.

A list of products returns in the response
*/
type ListProductsOK struct {
	Payload []*models.Product
}

func (o *ListProductsOK) Error() string {
	return fmt.Sprintf("[GET /products][%d] listProductsOK  %+v", 200, o.Payload)
}

func (o *ListProductsOK) GetPayload() []*models.Product {
	return o.Payload
}

func (o *ListProductsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListProductsBadRequest creates a ListProductsBadRequest with default headers values
func NewListProductsBadRequest() *ListProductsBadRequest {
	return &ListProductsBadRequest{}
}

/*ListProductsBadRequest handles this case with default header values.

Generic message returned as a string
*/
type ListProductsBadRequest struct {
	Payload *models.GenericError
}

func (o *ListProductsBadRequest) Error() string {
	return fmt.Sprintf("[GET /products][%d] listProductsBadRequest  %+v", 400, o.Payload)
}

func (o *ListProductsBadRequest) GetPayload() *models.GenericError {
	return o.Payload
}

func (o *ListProductsBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GenericError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
