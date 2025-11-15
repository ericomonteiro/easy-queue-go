package repositories

import (
	"context"
	"easy-queue-go/src/internal/models"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// BusinessRepository defines the interface for business operations
type BusinessRepository interface {
	Create(ctx context.Context, business *models.Business) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Business, error)
	FindByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*models.Business, error)
	FindAll(ctx context.Context) ([]*models.Business, error)
	Update(ctx context.Context, business *models.Business) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// businessRepository implements BusinessRepository
type businessRepository struct {
	pool *pgxpool.Pool
}

// NewBusinessRepository creates a new instance of BusinessRepository
func NewBusinessRepository(pool *pgxpool.Pool) BusinessRepository {
	return &businessRepository{
		pool: pool,
	}
}

// Create inserts a new business into the database
func (r *businessRepository) Create(ctx context.Context, business *models.Business) error {
	query := `
		INSERT INTO businesses (id, owner_id, name, description, address, phone, email, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.pool.Exec(ctx, query,
		business.ID,
		business.OwnerID,
		business.Name,
		business.Description,
		business.Address,
		business.Phone,
		business.Email,
		business.IsActive,
		business.CreatedAt,
		business.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create business: %w", err)
	}

	return nil
}

// FindByID retrieves a business by ID
func (r *businessRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Business, error) {
	query := `
		SELECT id, owner_id, name, description, address, phone, email, is_active, created_at, updated_at
		FROM businesses
		WHERE id = $1
	`

	business := &models.Business{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&business.ID,
		&business.OwnerID,
		&business.Name,
		&business.Description,
		&business.Address,
		&business.Phone,
		&business.Email,
		&business.IsActive,
		&business.CreatedAt,
		&business.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("business not found")
		}
		return nil, fmt.Errorf("failed to find business: %w", err)
	}

	return business, nil
}

// FindByOwnerID retrieves all businesses owned by a specific user
func (r *businessRepository) FindByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*models.Business, error) {
	query := `
		SELECT id, owner_id, name, description, address, phone, email, is_active, created_at, updated_at
		FROM businesses
		WHERE owner_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query businesses: %w", err)
	}
	defer rows.Close()

	var businesses []*models.Business
	for rows.Next() {
		business := &models.Business{}
		err := rows.Scan(
			&business.ID,
			&business.OwnerID,
			&business.Name,
			&business.Description,
			&business.Address,
			&business.Phone,
			&business.Email,
			&business.IsActive,
			&business.CreatedAt,
			&business.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan business: %w", err)
		}
		businesses = append(businesses, business)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating businesses: %w", err)
	}

	return businesses, nil
}

// FindAll returns all businesses
func (r *businessRepository) FindAll(ctx context.Context) ([]*models.Business, error) {
	query := `
		SELECT id, owner_id, name, description, address, phone, email, is_active, created_at, updated_at
		FROM businesses
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query businesses: %w", err)
	}
	defer rows.Close()

	var businesses []*models.Business
	for rows.Next() {
		business := &models.Business{}
		err := rows.Scan(
			&business.ID,
			&business.OwnerID,
			&business.Name,
			&business.Description,
			&business.Address,
			&business.Phone,
			&business.Email,
			&business.IsActive,
			&business.CreatedAt,
			&business.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan business: %w", err)
		}
		businesses = append(businesses, business)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating businesses: %w", err)
	}

	return businesses, nil
}

// Update updates an existing business
func (r *businessRepository) Update(ctx context.Context, business *models.Business) error {
	query := `
		UPDATE businesses
		SET name = $2, description = $3, address = $4, phone = $5, email = $6, is_active = $7, updated_at = $8
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query,
		business.ID,
		business.Name,
		business.Description,
		business.Address,
		business.Phone,
		business.Email,
		business.IsActive,
		business.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update business: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("business not found")
	}

	return nil
}

// Delete removes a business from the database
func (r *businessRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM businesses WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete business: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("business not found")
	}

	return nil
}
