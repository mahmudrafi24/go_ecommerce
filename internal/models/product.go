package models

import (
	"errors"
	"time"
)

// Product represents an item available for purchase
type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Product validation errors
var (
	ErrEmptyProductName = errors.New("product name is required")
	ErrInvalidPrice     = errors.New("price must be greater than or equal to 0")
	ErrInvalidStock     = errors.New("stock must be greater than or equal to 0")
	ErrEmptyDescription = errors.New("description is required")
)

// ProductInput represents the input for creating/updating a product
type ProductInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

// Validate validates the product input
func (p *ProductInput) Validate() ValidationErrors {
	errs := make(ValidationErrors)

	if p.Name == "" {
		errs["name"] = ErrEmptyProductName.Error()
	}
	if p.Price < 0 {
		errs["price"] = ErrInvalidPrice.Error()
	}
	if p.Stock < 0 {
		errs["stock"] = ErrInvalidStock.Error()
	}

	return errs
}

// Validate validates the product
func (p *Product) Validate() ValidationErrors {
	errs := make(ValidationErrors)

	if p.Name == "" {
		errs["name"] = ErrEmptyProductName.Error()
	}
	if p.Price < 0 {
		errs["price"] = ErrInvalidPrice.Error()
	}
	if p.Stock < 0 {
		errs["stock"] = ErrInvalidStock.Error()
	}

	return errs
}

// HasSufficientStock checks if the product has enough stock
func (p *Product) HasSufficientStock(quantity int) bool {
	return p.Stock >= quantity
}

// ReduceStock reduces the product stock by the given quantity
func (p *Product) ReduceStock(quantity int) error {
	if !p.HasSufficientStock(quantity) {
		return errors.New("insufficient stock")
	}
	p.Stock -= quantity
	return nil
}
