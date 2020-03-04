package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

// Product represents a product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// Products ...
type Products []*Product

// ErrProductNotFound is a structured error
var ErrProductNotFound = fmt.Errorf("Product not found")

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
		DeletedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd343",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
		DeletedOn:   time.Now().UTC().String(),
	},
}

// Validate validates a product
func (p *Product) Validate() error {
	validate := validator.New()

	validate.RegisterValidation("sku", func(fl validator.FieldLevel) bool {
		// sku is of the format abc-abcd-abcdef
		re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]`)
		matches := re.FindAllString(fl.Field().String(), -1)

		if len(matches) != 1 {
			return false
		}

		return true
	})
	return validate.Struct(p)
}

// The sku validation can be defined as a separate function and passed on to the register validation function
// func validateSKU(f1 validator.FieldLevel) {
// }

// FromJSON ...
func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

// ToJSON returns a json representation of the product list
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

// GetProducts returns list of products in the application
func GetProducts() Products {
	return productList
}

// AddProduct ...
func AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

// UpdateProduct ...
// This is likely to overwrite other properties of the product
func UpdateProduct(id int, p *Product) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}

	// Update product list
	p.ID = id
	productList[pos] = p
	return nil
}

func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, -1, ErrProductNotFound
}

func getNextID() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}
