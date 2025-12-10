package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"go-ecommerce-demo/internal/models"
)

// CartRepository handles cart data access operations
type CartRepository struct {
	pool *pgxpool.Pool
}

// NewCartRepository creates a new CartRepository
func NewCartRepository(pool *pgxpool.Pool) *CartRepository {
	return &CartRepository{pool: pool}
}

// GetByUserID retrieves a user's cart with all items and product details
func (r *CartRepository) GetByUserID(ctx context.Context, userID int64) (*models.Cart, error) {
	query := `
		SELECT 
			ci.id, ci.user_id, ci.product_id, ci.quantity, ci.created_at,
			p.id, p.name, p.description, p.price, p.stock, p.created_at, p.updated_at
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		WHERE ci.user_id = $1
		ORDER BY ci.created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cart := &models.Cart{
		UserID: userID,
		Items:  make([]models.CartItem, 0),
	}

	for rows.Next() {
		var item models.CartItem
		var product models.Product

		err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.ProductID,
			&item.Quantity,
			&item.CreatedAt,
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		item.Product = &product
		cart.Items = append(cart.Items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	cart.CalculateTotal()
	return cart, nil
}

// AddItem adds a product to the user's cart or updates quantity if it already exists
func (r *CartRepository) AddItem(ctx context.Context, userID, productID int64, quantity int) error {
	// Use upsert to handle both insert and update cases
	query := `
		INSERT INTO cart_items (user_id, product_id, quantity, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id, product_id)
		DO UPDATE SET quantity = cart_items.quantity + $3
	`

	_, err := r.pool.Exec(ctx, query, userID, productID, quantity)
	return err
}

// UpdateItemQuantity updates the quantity of a cart item
func (r *CartRepository) UpdateItemQuantity(ctx context.Context, userID, productID int64, quantity int) error {
	query := `
		UPDATE cart_items
		SET quantity = $1
		WHERE user_id = $2 AND product_id = $3
	`

	result, err := r.pool.Exec(ctx, query, quantity, userID, productID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// RemoveItem removes a product from the user's cart
func (r *CartRepository) RemoveItem(ctx context.Context, userID, productID int64) error {
	query := `DELETE FROM cart_items WHERE user_id = $1 AND product_id = $2`

	result, err := r.pool.Exec(ctx, query, userID, productID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// Clear removes all items from the user's cart
func (r *CartRepository) Clear(ctx context.Context, userID int64) error {
	query := `DELETE FROM cart_items WHERE user_id = $1`

	_, err := r.pool.Exec(ctx, query, userID)
	return err
}

// GetItem retrieves a specific cart item
func (r *CartRepository) GetItem(ctx context.Context, userID, productID int64) (*models.CartItem, error) {
	query := `
		SELECT 
			ci.id, ci.user_id, ci.product_id, ci.quantity, ci.created_at,
			p.id, p.name, p.description, p.price, p.stock, p.created_at, p.updated_at
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		WHERE ci.user_id = $1 AND ci.product_id = $2
	`

	var item models.CartItem
	var product models.Product

	err := r.pool.QueryRow(ctx, query, userID, productID).Scan(
		&item.ID,
		&item.UserID,
		&item.ProductID,
		&item.Quantity,
		&item.CreatedAt,
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	item.Product = &product
	return &item, nil
}
