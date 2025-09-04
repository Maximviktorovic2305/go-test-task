package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"effective-mobile-subscription-service/internal/models"
)

// SubscriptionRepository defines the interface for subscription data operations
type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *models.Subscription) error
	GetByID(ctx context.Context, id string) (*models.Subscription, error)
	Update(ctx context.Context, subscription *models.Subscription) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filters ListFilters) ([]*models.Subscription, error)
	CalculateCost(ctx context.Context, filters CostFilters) (int, error)
}

// ListFilters represents filters for listing subscriptions
type ListFilters struct {
	Page        int
	Limit       int
	UserID      string
	ServiceName string
}

// CostFilters represents filters for cost calculation
type CostFilters struct {
	UserID      string
	ServiceName string
	FromDate    time.Time
	ToDate      time.Time
}

// subscriptionRepository implements SubscriptionRepository
type subscriptionRepository struct {
	db *pgxpool.Pool
}

// NewSubscriptionRepository creates a new subscription repository
func NewSubscriptionRepository(db *pgxpool.Pool) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

// Create inserts a new subscription into the database
func (r *subscriptionRepository) Create(ctx context.Context, subscription *models.Subscription) error {
	query := `
		INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(ctx, query,
		subscription.ID,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate,
		subscription.EndDate,
		subscription.CreatedAt,
		subscription.UpdatedAt,
	)
	return err
}

// GetByID retrieves a subscription by its ID
func (r *subscriptionRepository) GetByID(ctx context.Context, id string) (*models.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions
		WHERE id = $1
	`
	subscription := &models.Subscription{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&subscription.ID,
		&subscription.ServiceName,
		&subscription.Price,
		&subscription.UserID,
		&subscription.StartDate,
		&subscription.EndDate,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("subscription not found")
		}
		return nil, err
	}
	return subscription, nil
}

// Update updates an existing subscription
func (r *subscriptionRepository) Update(ctx context.Context, subscription *models.Subscription) error {
	query := `
		UPDATE subscriptions
		SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5, updated_at = $6
		WHERE id = $7
	`
	_, err := r.db.Exec(ctx, query,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate,
		subscription.EndDate,
		subscription.UpdatedAt,
		subscription.ID,
	)
	return err
}

// Delete removes a subscription by its ID
func (r *subscriptionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("subscription not found")
	}
	
	return nil
}

// List retrieves subscriptions with optional filters and pagination
func (r *subscriptionRepository) List(ctx context.Context, filters ListFilters) ([]*models.Subscription, error) {
	// Base query
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions
		WHERE 1=1
	`
	
	// Parameters for the query
	params := []interface{}{}
	paramCount := 1
	
	// Apply filters
	if filters.UserID != "" {
		query += fmt.Sprintf(" AND user_id = $%d", paramCount)
		params = append(params, filters.UserID)
		paramCount++
	}
	
	if filters.ServiceName != "" {
		query += fmt.Sprintf(" AND service_name = $%d", paramCount)
		params = append(params, filters.ServiceName)
		paramCount++
	}
	
	// Apply pagination
	if filters.Limit <= 0 {
		filters.Limit = 10
	}
	if filters.Page <= 0 {
		filters.Page = 1
	}
	
	offset := (filters.Page - 1) * filters.Limit
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", paramCount, paramCount+1)
	params = append(params, filters.Limit, offset)
	
	// Execute query
	rows, err := r.db.Query(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	// Scan results
	subscriptions := []*models.Subscription{}
	for rows.Next() {
		subscription := &models.Subscription{}
		err := rows.Scan(
			&subscription.ID,
			&subscription.ServiceName,
			&subscription.Price,
			&subscription.UserID,
			&subscription.StartDate,
			&subscription.EndDate,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}
	
	return subscriptions, nil
}

// CalculateCost calculates the total cost of subscriptions with optional filters
func (r *subscriptionRepository) CalculateCost(ctx context.Context, filters CostFilters) (int, error) {
	query := `
		SELECT COALESCE(SUM(price), 0)
		FROM subscriptions
		WHERE 1=1
	`
	
	// Parameters for the query
	params := []interface{}{}
	paramCount := 1
	
	// Apply filters
	if filters.UserID != "" {
		query += fmt.Sprintf(" AND user_id = $%d", paramCount)
		params = append(params, filters.UserID)
		paramCount++
	}
	
	if filters.ServiceName != "" {
		query += fmt.Sprintf(" AND service_name = $%d", paramCount)
		params = append(params, filters.ServiceName)
		paramCount++
	}
	
	// Apply date range filter
	if !filters.FromDate.IsZero() {
		query += fmt.Sprintf(" AND (end_date IS NULL OR end_date >= $%d)", paramCount)
		params = append(params, filters.FromDate)
		paramCount++
	}
	
	if !filters.ToDate.IsZero() {
		query += fmt.Sprintf(" AND start_date <= $%d", paramCount)
		params = append(params, filters.ToDate)
	}
	
	// Execute query
	var totalCost int
	err := r.db.QueryRow(ctx, query, params...).Scan(&totalCost)
	if err != nil {
		return 0, err
	}
	
	return totalCost, nil
}