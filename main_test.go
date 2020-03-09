package main

import (
	"fmt"
	"testing"

	"github.com/agbanagba/go-microservice/client/client"
	"github.com/agbanagba/go-microservice/client/client/products"
)

func TestClient(t *testing.T) {
	config := client.DefaultTransportConfig().WithHost("localhost:9090")
	c := client.NewHTTPClientWithConfig(nil, config)
	params := products.NewListProductsParams()
	prod, err := c.Products.ListProducts(params)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(prod)
}
