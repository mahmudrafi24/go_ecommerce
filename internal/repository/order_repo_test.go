package repository

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"testing"
	"testing/quick"

	"github.com/jackc/pgx/v5/pgxpool"

	"go-ecommerce-demo/internal/models"
)

// **Feature: go-ecommerce-demo, Property 14: Data persistence round-trip (Order)**
// **Validates: Requirements 6.1, 6.2**
//
// For any entity (Order), storing to PostgreSQL and retrieving should return
// an equivalent entity with matching field values.

// testDB holds the database pool for tests
var testDB *pgxpool.Pool

// TestMain sets up and tears down the test database connection
func TestMain(m *testing.M) {
	ctx := context.Background()

	// Get database URL from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/ecommerce_test?sslmode=disable"
	}

	// Connect to database
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		fmt.Println("Skipping repository tests - no database connection")
		os.Exit(0)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		fmt.Printf("Failed to ping database: %v\n", err)
		fmt.Println("Skipping repository tests - database not responding")
		pool.Close()
		os.Exit(0)
	}

	testDB = pool

	// Run tests
	code := m.Run()

	// Cleanup
	pool.Close()
	os.Exit(code)
}

// cleanupTestData removes test data from the database
func cleanupTestData(ctx context.Context, t *testing.T) {
	// Delete in order to respect foreign key constraints
	_, err := testDB.Exec(ctx, "DELETE FROM order_items")
	if err != nil {
		t.Logf("Warning: failed to cleanup order_items: %v", err)
	}
	_, err = testDB.Exec(ctx, "DELETE FROM orders")
	if err != nil {
		t.Logf("Warning: failed to cleanup orders: %v", err)
	}
	_, err = testDB.Exec(ctx, "DELETE FROM cart_items")
	if err != nil {
		t.Logf("Warning: failed to cleanup cart_items: %v", err)
	}
	_, err = testDB.Exec(ctx, "DELETE FROM products")
	if err != nil {
		t.Logf("Warning: failed to cleanup products: %v", err)
	}
	_, err = testDB.Exec(ctx, "DELETE FROM users")
	if err != nil {
		t.Logf("Warning: failed to cleanup users: %v", err)
	}
}

// createTestUser creates a user for testing and returns its ID
func createTestUser(ctx context.Context, r *rand.Rand) (*models.User, error) {
	userRepo := NewUserRepository(testDB)
	user := &models.User{
		Email:        fmt.Sprintf("test_%d@example.com", r.Int63()),
		PasswordHash: "hashedpassword123",
		Name:         fmt.Sprintf("Test User %d", r.Intn(10000)),
		Role:         "user",
	}
	err := userRepo.Create(ctx, user)
	return user, err
}

// createTestProduct creates a product for testing and returns it
func createTestProduct(ctx context.Context, r *rand.Rand) (*models.Product, error) {
	productRepo := NewProductRepository(testDB)
	product := &models.Product{
		Name:        fmt.Sprintf("Test Product %d", r.Intn(10000)),
		Description: fmt.Sprintf("Description for product %d", r.Intn(10000)),
		Price:       float64(r.Intn(10000)+1) / 100.0, // 0.01 to 100.00
		Stock:       r.Intn(1000) + 1,                 // 1 to 1000
	}
	err := productRepo.Create(ctx, product)
	return product, err
}

// generateRandomOrderItem creates a random order item with a valid product
func generateRandomOrderItem(r *rand.Rand, productID int64, productPrice float64) models.OrderItem {
	quantity := r.Intn(10) + 1 // 1 to 10
	return models.OrderItem{
		ProductID: productID,
		Quantity:  quantity,
		Price:     productPrice,
	}
}

// floatsEqual compares two floats with a small epsilon for floating point precision
func floatsEqual(a, b float64) bool {
	epsilon := 0.01
	return math.Abs(a-b) < epsilon
}

// ordersEqual compares two orders for equality (excluding auto-generated fields)
func ordersEqual(original, retrieved *models.Order) bool {
	if original.UserID != retrieved.UserID {
		return false
	}
	if !floatsEqual(original.Total, retrieved.Total) {
		return false
	}
	if original.Status != retrieved.Status {
		return false
	}
	if len(original.Items) != len(retrieved.Items) {
		return false
	}

	// Compare items (order should be preserved by ORDER BY id)
	for i := range original.Items {
		origItem := original.Items[i]
		retItem := retrieved.Items[i]

		if origItem.ProductID != retItem.ProductID {
			return false
		}
		if origItem.Quantity != retItem.Quantity {
			return false
		}
		if !floatsEqual(origItem.Price, retItem.Price) {
			return false
		}
	}

	return true
}

// TestProperty14_OrderPersistenceRoundTrip tests that orders can be stored and retrieved
func TestProperty14_OrderPersistenceRoundTrip(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	ctx := context.Background()
	cleanupTestData(ctx, t)
	defer cleanupTestData(ctx, t)

	config := &quick.Config{
		MaxCount: 100,
	}

	// Property: For any valid order, storing and retrieving should return equivalent data
	f := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		orderRepo := NewOrderRepository(testDB)

		// Create a test user (required for foreign key)
		user, err := createTestUser(ctx, r)
		if err != nil {
			t.Logf("Failed to create test user: %v", err)
			return false
		}

		// Create 1-5 test products for order items
		numProducts := r.Intn(5) + 1
		products := make([]*models.Product, numProducts)
		for i := 0; i < numProducts; i++ {
			product, err := createTestProduct(ctx, r)
			if err != nil {
				t.Logf("Failed to create test product: %v", err)
				return false
			}
			products[i] = product
		}

		// Create order items from products
		items := make([]models.OrderItem, numProducts)
		var total float64
		for i, product := range products {
			items[i] = generateRandomOrderItem(r, product.ID, product.Price)
			total += float64(items[i].Quantity) * items[i].Price
		}

		// Create the order
		statuses := []string{models.OrderStatusPending, models.OrderStatusCompleted, models.OrderStatusCancelled}
		order := &models.Order{
			UserID: user.ID,
			Items:  items,
			Total:  total,
			Status: statuses[r.Intn(len(statuses))],
		}

		// Store the order
		err = orderRepo.Create(ctx, order)
		if err != nil {
			t.Logf("Failed to create order: %v", err)
			return false
		}

		// Verify ID was assigned
		if order.ID == 0 {
			t.Log("Order ID was not assigned after creation")
			return false
		}

		// Retrieve the order
		retrieved, err := orderRepo.GetByID(ctx, order.ID)
		if err != nil {
			t.Logf("Failed to retrieve order: %v", err)
			return false
		}

		// Compare original and retrieved
		if !ordersEqual(order, retrieved) {
			t.Logf("Order mismatch: original=%+v, retrieved=%+v", order, retrieved)
			return false
		}

		return true
	}

	if err := quick.Check(f, config); err != nil {
		t.Errorf("Property 14 failed (Order round-trip): %v", err)
	}
}

// TestProperty14_OrderWithMultipleItemsRoundTrip tests orders with varying item counts
func TestProperty14_OrderWithMultipleItemsRoundTrip(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	ctx := context.Background()
	cleanupTestData(ctx, t)
	defer cleanupTestData(ctx, t)

	config := &quick.Config{
		MaxCount: 100,
	}

	// Property: Orders with any number of items (1-10) should round-trip correctly
	f := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		orderRepo := NewOrderRepository(testDB)

		// Create a test user
		user, err := createTestUser(ctx, r)
		if err != nil {
			t.Logf("Failed to create test user: %v", err)
			return false
		}

		// Create varying number of products (1-10)
		numProducts := r.Intn(10) + 1
		products := make([]*models.Product, numProducts)
		for i := 0; i < numProducts; i++ {
			product, err := createTestProduct(ctx, r)
			if err != nil {
				t.Logf("Failed to create test product: %v", err)
				return false
			}
			products[i] = product
		}

		// Create order items
		items := make([]models.OrderItem, numProducts)
		var total float64
		for i, product := range products {
			items[i] = generateRandomOrderItem(r, product.ID, product.Price)
			total += float64(items[i].Quantity) * items[i].Price
		}

		order := &models.Order{
			UserID: user.ID,
			Items:  items,
			Total:  total,
			Status: models.OrderStatusPending,
		}

		// Store and retrieve
		err = orderRepo.Create(ctx, order)
		if err != nil {
			t.Logf("Failed to create order: %v", err)
			return false
		}

		retrieved, err := orderRepo.GetByID(ctx, order.ID)
		if err != nil {
			t.Logf("Failed to retrieve order: %v", err)
			return false
		}

		// Verify item count matches
		if len(retrieved.Items) != numProducts {
			t.Logf("Item count mismatch: expected %d, got %d", numProducts, len(retrieved.Items))
			return false
		}

		return ordersEqual(order, retrieved)
	}

	if err := quick.Check(f, config); err != nil {
		t.Errorf("Property 14 failed (multiple items round-trip): %v", err)
	}
}

// TestProperty14_OrderListByUserRoundTrip tests that orders can be listed by user
func TestProperty14_OrderListByUserRoundTrip(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	ctx := context.Background()
	cleanupTestData(ctx, t)
	defer cleanupTestData(ctx, t)

	config := &quick.Config{
		MaxCount: 50, // Fewer iterations since we create multiple orders per iteration
	}

	// Property: All orders created for a user should be retrievable via ListByUserID
	f := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		orderRepo := NewOrderRepository(testDB)

		// Create a test user
		user, err := createTestUser(ctx, r)
		if err != nil {
			t.Logf("Failed to create test user: %v", err)
			return false
		}

		// Create 1-5 orders for this user
		numOrders := r.Intn(5) + 1
		createdOrders := make([]*models.Order, numOrders)

		for i := 0; i < numOrders; i++ {
			// Create a product for this order
			product, err := createTestProduct(ctx, r)
			if err != nil {
				t.Logf("Failed to create test product: %v", err)
				return false
			}

			item := generateRandomOrderItem(r, product.ID, product.Price)
			total := float64(item.Quantity) * item.Price

			order := &models.Order{
				UserID: user.ID,
				Items:  []models.OrderItem{item},
				Total:  total,
				Status: models.OrderStatusPending,
			}

			err = orderRepo.Create(ctx, order)
			if err != nil {
				t.Logf("Failed to create order: %v", err)
				return false
			}
			createdOrders[i] = order
		}

		// List orders for user
		orders, total, err := orderRepo.ListByUserID(ctx, user.ID, 1, 100)
		if err != nil {
			t.Logf("Failed to list orders: %v", err)
			return false
		}

		// Verify count matches
		if total != numOrders {
			t.Logf("Order count mismatch: expected %d, got %d", numOrders, total)
			return false
		}

		if len(orders) != numOrders {
			t.Logf("Returned orders mismatch: expected %d, got %d", numOrders, len(orders))
			return false
		}

		// Verify all orders belong to the user
		for _, order := range orders {
			if order.UserID != user.ID {
				t.Logf("Order user ID mismatch: expected %d, got %d", user.ID, order.UserID)
				return false
			}
		}

		return true
	}

	if err := quick.Check(f, config); err != nil {
		t.Errorf("Property 14 failed (list by user round-trip): %v", err)
	}
}
