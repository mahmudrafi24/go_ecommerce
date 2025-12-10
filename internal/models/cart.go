package models

import (
	"time"
)

// CartItem represents an item in a user's shopping cart
type CartItem struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	ProductID int64     `json:"product_id"`
	Product   *Product  `json:"product,omitempty"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}

// Cart represents a user's shopping cart
type Cart struct {
	UserID int64      `json:"user_id"`
	Items  []CartItem `json:"items"`
	Total  float64    `json:"total"`
}

// CalculateTotal calculates the total price of all items in the cart
func (c *Cart) CalculateTotal() float64 {
	var total float64
	for _, item := range c.Items {
		if item.Product != nil {
			total += float64(item.Quantity) * item.Product.Price
		}
	}
	c.Total = total
	return total
}

// ItemCount returns the number of distinct items in the cart
func (c *Cart) ItemCount() int {
	return len(c.Items)
}

// TotalQuantity returns the total quantity of all items in the cart
func (c *Cart) TotalQuantity() int {
	var total int
	for _, item := range c.Items {
		total += item.Quantity
	}
	return total
}

// IsEmpty returns true if the cart has no items
func (c *Cart) IsEmpty() bool {
	return len(c.Items) == 0
}

// FindItem finds a cart item by product ID
func (c *Cart) FindItem(productID int64) *CartItem {
	for i := range c.Items {
		if c.Items[i].ProductID == productID {
			return &c.Items[i]
		}
	}
	return nil
}

// AddItemInput represents the input for adding an item to cart
type AddItemInput struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

// Validate validates the add item input
func (a *AddItemInput) Validate() ValidationErrors {
	errs := make(ValidationErrors)

	if a.ProductID <= 0 {
		errs["product_id"] = "product_id must be a positive integer"
	}
	if a.Quantity <= 0 {
		errs["quantity"] = "quantity must be a positive integer"
	}

	return errs
}

// UpdateQuantityInput represents the input for updating cart item quantity
type UpdateQuantityInput struct {
	Quantity int `json:"quantity"`
}

// Validate validates the update quantity input
func (u *UpdateQuantityInput) Validate() ValidationErrors {
	errs := make(ValidationErrors)

	if u.Quantity <= 0 {
		errs["quantity"] = "quantity must be a positive integer"
	}

	return errs
}
