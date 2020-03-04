package data

import "testing"

func TestProductValidation(t *testing.T) {
	p := &Product{
		Name:  "Tega",
		Price: 12,
		SKU:   "abc-abc-abc",
	}

	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}

}
