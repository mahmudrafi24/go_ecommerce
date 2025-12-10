package models

import (
	"time"
)

// Order status constants
const (
	OrderStatusPending   = "pending"
	OrderStatusCompleted = "completed"
	OrderStatusCancelled = "cancelled"
)

// OrderItem represents an item in an order
type OrderItem struct {
	ID        int64   `json:"id"`
	OrderID   int64   `json:"order_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"` // Price at time of order
}

// Order represents a confirmed purchase
type Order struct {
	ID        int64       `json:"id"`
	UserID    int64       `json:"user_id"`
	Items     []OrderItem `json:"items"`
	Total     float64     `json:"total"`
	Status    string      `json:"status"` // "pending", "completed", "cancelled"
	CreatedAt time.Time   `json:"created_at"`
}

// CalculateTotal calculates the total price of all items in the order
func (o *Order) CalculateTotal() float64 {
	var total float64
	for _, item := range o.Items {
		total += float64(item.Quantity) * item.Price
	}
	o.Total = total
	return total
}

// ItemCount returns the number of distinct items in the order
func (o *Order) ItemCount() int {
	return len(o.Items)
}

// TotalQuantity returns the total quantity of all items in the order
func (o *Order) TotalQuantity() int {
	var total int
	for _, item := range o.Items {
		total += item.Quantity
	}
	return total
}

// IsPending returns true if the order is pending
func (o *Order) IsPending() bool {
	return o.Status == OrderStatusPending
}

// IsCompleted returns true if the order is completed
func (o *Order) IsCompleted() bool {
	return o.Status == OrderStatusCompleted
}

// IsCancelled returns true if the order is cancelled
func (o *Order) IsCancelled() bool {
	return o.Status == OrderStatusCancelled
}

// NewOrderFromCart creates a new order from a cart
func NewOrderFromCart(cart *Cart) *Order {
	order := &Order{
		UserID: cart.UserID,
		Status: OrderStatusPending,
		Items:  make([]OrderItem, 0, len(cart.Items)),
	}

	for _, cartItem := range cart.Items {
		if cartItem.Product != nil {
			orderItem := OrderItem{
				ProductID: cartItem.ProductID,
				Quantity:  cartItem.Quantity,
				Price:     cartItem.Product.Price,
			}
			order.Items = append(order.Items, orderItem)
		}
	}

	order.CalculateTotal()
	return order
}
