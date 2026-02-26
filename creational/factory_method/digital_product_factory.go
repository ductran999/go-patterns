package factorymethod

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidProductLink = errors.New("invalid download link for digital product")
)

// DigitalProductFactory creates digital products
type DigitalProductFactory struct{}

func (f *DigitalProductFactory) CreateProduct(name string, price float64, extra string) (Product, error) {
	baseProduct := baseProduct{
		name:  name,
		price: price,
	}

	if err := baseProduct.validateName(); err != nil {
		return nil, err
	}

	if err := baseProduct.validatePrice(); err != nil {
		return nil, err
	}

	if strings.Trim(extra, " ") == "" {
		return nil, ErrInvalidProductLink
	}

	return &DigitalProduct{
		baseProduct:  baseProduct,
		downloadLink: extra,
	}, nil
}

// DigitalProduct represents a downloadable item like an e-book or software
type DigitalProduct struct {
	baseProduct

	downloadLink string
}

func (d *DigitalProduct) GetDetails() string {
	return fmt.Sprintf("Digital Product: %s, Price: $%.2f, Download: %s", d.name, d.price, d.downloadLink)
}

func (d *DigitalProduct) CalculateShippingCost() float64 {
	// No shipping cost for digital products
	return 0.0
}
