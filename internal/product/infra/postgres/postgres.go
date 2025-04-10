package postgres

import (
	"context"

	"github.com/0x0FACED/pvz-avito/internal/product/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewProductPostgresRepository(pgx *pgxpool.Pool) *ProductPostgresRepository {
	return &ProductPostgresRepository{pool: pgx}
}

func (r *ProductPostgresRepository) Create(ctx context.Context, product *domain.Product) (*domain.Product, error) {
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

	var created domain.Product
	created.Type = product.Type
	created.ReceptionID = product.ReceptionID

	err := r.pool.QueryRow(ctx, query, args).Scan(&created.ID, &created.DateTime)
	if err != nil {
		// TODO: handle err
		return nil, err
	}

	return &created, nil
}

func (r *ProductPostgresRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	query := `
		SELECT id, date_time, type, reception_id
		FROM avito.products
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	product := domain.Product{}
	err := r.pool.QueryRow(ctx, query, args).Scan(
		&product.ID,
		&product.DateTime,
		&product.Type,
		&product.ReceptionID,
	)
	if err != nil {
		// TODO: handle err
		return nil, err
	}

	return &product, nil
}

func (r *ProductPostgresRepository) GetLastByReception(ctx context.Context, receptionID string) (*domain.Product, error) {
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

	product := domain.Product{}
	err := r.pool.QueryRow(ctx, query, args).Scan(
		&product.ID,
		&product.DateTime,
		&product.Type,
		&product.ReceptionID,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			// TODO: handle err
			return nil, nil
		}
		// TODO: handle err
		return nil, err
	}

	return &product, nil
}

func (r *ProductPostgresRepository) DeleteLastFromReception(ctx context.Context, receptionID string) error {
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

	args := pgx.NamedArgs{
		"reception_id": receptionID,
	}

	_, err := r.pool.Exec(ctx, query, args)
	return err
}

func (r *ProductPostgresRepository) ListByReception(ctx context.Context, receptionID string) ([]*domain.Product, error) {
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
		// TODO: handle err
		return nil, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var product domain.Product
		err := rows.Scan(
			&product.ID,
			&product.DateTime,
			&product.Type,
			&product.ReceptionID,
		)
		if err != nil {
			// TODO: handle err
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}
