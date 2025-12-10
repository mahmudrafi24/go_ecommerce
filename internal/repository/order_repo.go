package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"go-ecommerce-demo/internal/models"
)

// OrderRepository handles order data access operations
type OrderRepository struct {
	pool *pgxpool.Pool
}

// NewOrderRepository creates a new OrderRepository
func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

// Create inserts a new order with its items using a transaction
func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Insert order
	orderQuery := `
		INSERT INTO orders (user_id, total, status, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, created_at
	`

	err = tx.QueryRow(ctx, orderQuery,
		order.UserID,
		order.Total,
		order.Status,
	).Scan(&order.ID, &order.CreatedAt)

	if err != nil {
		return err
	}

	// Insert order items
	itemQuery := `
		INSERT INTO order_items (order_id, product_id, quantity, price)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	for i := range order.Items {
		err = tx.QueryRow(ctx, itemQuery,
			order.ID,
			order.Items[i].ProductID,
			order.Items[i].Quantity,
			order.Items[i].Price,
		).Scan(&order.Items[i].ID)

		if err != nil {
			return err
		}
		order.Items[i].OrderID = order.ID
	}

	return tx.Commit(ctx)
}

// GetByID retrieves an order by its ID with all items
func (r *OrderRepository) GetByID(ctx context.Context, id int64) (*models.Order, error) {
	// Get order
	orderQuery := `
		SELECT id, user_id, total, status, created_at
		FROM orders
		WHERE id = $1
	`

	order := &models.Order{}
	err := r.pool.QueryRow(ctx, orderQuery, id).Scan(
		&order.ID,
		&order.UserID,
		&order.Total,
		&order.Status,
		&order.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// Get order items
	itemsQuery := `
		SELECT id, order_id, product_id, quantity, price
		FROM order_items
		WHERE order_id = $1
		ORDER BY id
	`

	rows, err := r.pool.Query(ctx, itemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	order.Items = make([]models.OrderItem, 0)
	for rows.Next() {
		var item models.OrderItem
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
		)
		if err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return order, nil
}

// ListByUserID retrieves a paginated list of orders for a user
func (r *OrderRepository) ListByUserID(ctx context.Context, userID int64, page, pageSize int) ([]models.Order, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM orders WHERE user_id = $1`
	err := r.pool.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get orders
	ordersQuery := `
		SELECT id, user_id, total, status, created_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, ordersQuery, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	orders := make([]models.Order, 0)
	for rows.Next() {
		var o models.Order
		err := rows.Scan(
			&o.ID,
			&o.UserID,
			&o.Total,
			&o.Status,
			&o.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	// Load items for each order
	for i := range orders {
		items, err := r.getOrderItems(ctx, orders[i].ID)
		if err != nil {
			return nil, 0, err
		}
		orders[i].Items = items
	}

	return orders, total, nil
}

// getOrderItems retrieves items for a specific order
func (r *OrderRepository) getOrderItems(ctx context.Context, orderID int64) ([]models.OrderItem, error) {
	query := `
		SELECT id, order_id, product_id, quantity, price
		FROM order_items
		WHERE order_id = $1
		ORDER BY id
	`

	rows, err := r.pool.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.OrderItem, 0)
	for rows.Next() {
		var item models.OrderItem
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

// BeginTx starts a new transaction and returns it
// This allows services to manage complex transactions spanning multiple repositories
func (r *OrderRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.pool.Begin(ctx)
}

// CreateWithTx creates an order within an existing transaction
func (r *OrderRepository) CreateWithTx(ctx context.Context, tx pgx.Tx, order *models.Order) error {
	// Insert order
	orderQuery := `
		INSERT INTO orders (user_id, total, status, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, created_at
	`

	err := tx.QueryRow(ctx, orderQuery,
		order.UserID,
		order.Total,
		order.Status,
	).Scan(&order.ID, &order.CreatedAt)

	if err != nil {
		return err
	}

	// Insert order items
	itemQuery := `
		INSERT INTO order_items (order_id, product_id, quantity, price)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	for i := range order.Items {
		err = tx.QueryRow(ctx, itemQuery,
			order.ID,
			order.Items[i].ProductID,
			order.Items[i].Quantity,
			order.Items[i].Price,
		).Scan(&order.Items[i].ID)

		if err != nil {
			return err
		}
		order.Items[i].OrderID = order.ID
	}

	return nil
}
