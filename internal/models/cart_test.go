package models

import (
	"math"
	"math/rand"
	"testing"
	"testing/quick"
)

// **Feature: go-ecommerce-demo, Property 10: Cart total equals sum of item prices**
// **Validates: Requirements 4.4**
//
// For any cart state, the cart total should equal the sum of (item.quantity * item.product.price)
// for all items in the cart.

// generateRandomProduct creates a random product with valid price and stock
func generateRandomProduct(r *rand.Rand, id int64) *Product {
	// Generate price between 0.01 and 10000.00 (in cents, then convert)
	price := float64(r.Intn(1000000)+1) / 100.0 // 0.01 to 10000.00

	return &Product{
		ID:          id,
		Name:        "Test Product",
		Description: "Test Description",
		Price:       price,
		Stock:       r.Intn(1000) + 1, // 1 to 1000
	}
}

// generateRandomCartItem creates a random cart item with a product
func generateRandomCartItem(r *rand.Rand, id int64, userID int64) CartItem {
	product := generateRandomProduct(r, id)
	quantity := r.Intn(100) + 1 // 1 to 100

	return CartItem{
		ID:        id,
		UserID:    userID,
		ProductID: product.ID,
		Product:   product,
		Quantity:  quantity,
	}
}

// generateRandomCart creates a random cart with 0-20 items
func generateRandomCart(r *rand.Rand) Cart {
	userID := int64(r.Intn(10000) + 1)
	numItems := r.Intn(21) // 0 to 20 items

	items := make([]CartItem, numItems)
	for i := 0; i < numItems; i++ {
		items[i] = generateRandomCartItem(r, int64(i+1), userID)
	}

	return Cart{
		UserID: userID,
		Items:  items,
		Total:  0, // Will be calculated
	}
}

// manuallyCalculateTotal calculates the expected total by summing item prices
func manuallyCalculateTotal(cart *Cart) float64 {
	var total float64
	for _, item := range cart.Items {
		if item.Product != nil {
			total += float64(item.Quantity) * item.Product.Price
		}
	}
	return total
}

// floatsEqual compares two floats with a small epsilon for floating point precision
func floatsEqual(a, b float64) bool {
	epsilon := 0.0001
	return math.Abs(a-b) < epsilon
}

// TestProperty10_CartTotalEqualsSumOfItemPrices tests that cart total equals sum of item prices
func TestProperty10_CartTotalEqualsSumOfItemPrices(t *testing.T) {
	config := &quick.Config{
		MaxCount: 100,
	}

	// Property: For any cart, CalculateTotal() should return the sum of (quantity * price) for all items
	f := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		cart := generateRandomCart(r)

		// Calculate expected total manually
		expectedTotal := manuallyCalculateTotal(&cart)

		// Calculate total using the Cart method
		actualTotal := cart.CalculateTotal()

		// The calculated total should match the expected total
		if !floatsEqual(actualTotal, expectedTotal) {
			return false
		}

		// The cart's Total field should also be updated
		if !floatsEqual(cart.Total, expectedTotal) {
			return false
		}

		return true
	}

	if err := quick.Check(f, config); err != nil {
		t.Errorf("Property 10 failed: %v", err)
	}
}

// TestProperty10_EmptyCartTotalIsZero tests that an empty cart has zero total
func TestProperty10_EmptyCartTotalIsZero(t *testing.T) {
	config := &quick.Config{
		MaxCount: 100,
	}

	// Property: For any empty cart, the total should be 0
	f := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		userID := int64(r.Intn(10000) + 1)

		cart := Cart{
			UserID: userID,
			Items:  []CartItem{},
			Total:  0,
		}

		actualTotal := cart.CalculateTotal()

		return actualTotal == 0 && cart.Total == 0
	}

	if err := quick.Check(f, config); err != nil {
		t.Errorf("Property 10 failed (empty cart): %v", err)
	}
}

// TestProperty10_CartWithNilProductsHandled tests that items with nil products don't affect total
func TestProperty10_CartWithNilProductsHandled(t *testing.T) {
	config := &quick.Config{
		MaxCount: 100,
	}

	// Property: Items with nil products should not contribute to the total
	f := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		userID := int64(r.Intn(10000) + 1)
		numItems := r.Intn(10) + 1 // 1 to 10 items

		items := make([]CartItem, numItems)
		var expectedTotal float64

		for i := 0; i < numItems; i++ {
			item := generateRandomCartItem(r, int64(i+1), userID)

			// Randomly set some products to nil
			if r.Float32() < 0.3 { // 30% chance of nil product
				item.Product = nil
			} else {
				expectedTotal += float64(item.Quantity) * item.Product.Price
			}

			items[i] = item
		}

		cart := Cart{
			UserID: userID,
			Items:  items,
			Total:  0,
		}

		actualTotal := cart.CalculateTotal()

		return floatsEqual(actualTotal, expectedTotal)
	}

	if err := quick.Check(f, config); err != nil {
		t.Errorf("Property 10 failed (nil products): %v", err)
	}
}
