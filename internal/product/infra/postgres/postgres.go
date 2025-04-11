package postgres

import (
	"context"
	"errors"
	"fmt"

	product_domain "github.com/0x0FACED/pvz-avito/internal/product/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewProductPostgresRepository(pgx *pgxpool.Pool) *ProductPostgresRepository {
	return &ProductPostgresRepository{pool: pgx}
}

func (r *ProductPostgresRepository) Create(ctx context.Context, product *product_domain.Product) (*product_domain.Product, error) {
	query := `
		INSERT INTO avito.products (id, date_time, type, reception_id)
		VALUES (@id, @date_time, @type, @reception_id)
		RETURNING id, date_time
	`

	args := pgx.NamedArgs{
		"id":           product.ID,
		"date_time":    product.DateTime,
		"type":         product.Type,
		"reception_id": product.ReceptionID,
	}

	var created product_domain.Product
	created.Type = product.Type
	created.ReceptionID = product.ReceptionID

	err := r.pool.QueryRow(ctx, query, args).Scan(&created.ID, &created.DateTime)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return nil, fmt.Errorf("%w: %w", product_domain.ErrReceptionNotFound, err)
			}
		}
		return nil, fmt.Errorf("%w: %w", product_domain.ErrInternalDatabase, err)
	}

	return &created, nil
}

func (r *ProductPostgresRepository) GetByID(ctx context.Context, id string) (*product_domain.Product, error) {
	query := `
		SELECT id, date_time, type, reception_id
		FROM avito.products
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	product := product_domain.Product{}
	err := r.pool.QueryRow(ctx, query, args).Scan(
		&product.ID,
		&product.DateTime,
		&product.Type,
		&product.ReceptionID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %w", product_domain.ErrProductNotFound, err)
		}
		return nil, fmt.Errorf("%w: %w", product_domain.ErrInternalDatabase, err)
	}

	return &product, nil
}

func (r *ProductPostgresRepository) GetLastByReception(ctx context.Context, receptionID string) (*product_domain.Product, error) {
	query := `
		SELECT id, date_time, type, reception_id
		FROM avito.products
		WHERE reception_id = @reception_id
		ORDER BY date_time DESC
		LIMIT 1
	`

	args := pgx.NamedArgs{
		"reception_id": receptionID,
	}

	product := product_domain.Product{}
	err := r.pool.QueryRow(ctx, query, args).Scan(
		&product.ID,
		&product.DateTime,
		&product.Type,
		&product.ReceptionID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %w", product_domain.ErrProductNotFound, err)
		}
		return nil, fmt.Errorf("%w: %w", product_domain.ErrInternalDatabase, err)
	}

	return &product, nil
}

func (r *ProductPostgresRepository) DeleteLastFromReception(ctx context.Context, receptionID string) error {
	var count int

	checkQuery := `
		SELECT COUNT(*) 
		FROM avito.products 
		WHERE reception_id = @reception_id
	`

	args := pgx.NamedArgs{
		"reception_id": receptionID,
	}

	err := r.pool.QueryRow(ctx, checkQuery, args).Scan(&count)
	if err != nil {
		return fmt.Errorf("%w: %w", product_domain.ErrInternalDatabase, err)
	}

	if count == 0 {
		return fmt.Errorf("%w: no products found for reception_id: %s", product_domain.ErrNoProductsToDelete, receptionID)
	}

	query := `
		DELETE FROM avito.products
		WHERE id = (
			SELECT id
			FROM avito.products
			WHERE reception_id = @reception_id
			ORDER BY date_time DESC
			LIMIT 1
		)
	`

	_, err = r.pool.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("%w: %w", product_domain.ErrInternalDatabase, err)
	}

	return nil
}

func (r *ProductPostgresRepository) ListByReception(ctx context.Context, receptionID string) ([]*product_domain.Product, error) {
	query := `
		SELECT id, date_time, type, reception_id
		FROM avito.products
		WHERE reception_id = @reception_id
		ORDER BY date_time DESC
	`

	args := pgx.NamedArgs{
		"reception_id": receptionID,
	}

	rows, err := r.pool.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", product_domain.ErrInternalDatabase, err)
	}
	defer rows.Close()

	var products []*product_domain.Product
	for rows.Next() {
		var product product_domain.Product
		err := rows.Scan(
			&product.ID,
			&product.DateTime,
			&product.Type,
			&product.ReceptionID,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", product_domain.ErrInternalDatabase, err)
		}
		products = append(products, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %w", product_domain.ErrInternalDatabase, err)
	}

	return products, nil
}
