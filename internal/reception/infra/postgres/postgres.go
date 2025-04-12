package postgres

import (
	"context"
	"errors"
	"fmt"

	reception_domain "github.com/0x0FACED/pvz-avito/internal/reception/domain"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReceptionPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewReceptionPostgresRepository(pgx *pgxpool.Pool) *ReceptionPostgresRepository {
	return &ReceptionPostgresRepository{pool: pgx}
}

func (r *ReceptionPostgresRepository) Create(ctx context.Context, reception *reception_domain.Reception) (*reception_domain.Reception, error) {
	query := `
		INSERT INTO avito.receptions (id, date_time, pvz_id, status)
		VALUES (@id, @date_time, @pvz_id, @status)
		RETURNING id, date_time, pvz_id, status
	`

	args := pgx.NamedArgs{
		"id":        reception.ID,
		"date_time": reception.DateTime,
		"pvz_id":    reception.PVZID,
		"status":    reception.Status,
	}

	var created reception_domain.Reception
	err := r.pool.QueryRow(ctx, query, args).Scan(
		&created.ID, &created.DateTime, &created.PVZID, &created.Status,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return nil, fmt.Errorf("%w: %w", reception_domain.ErrPVZNotFound, err)
			}
		}
		return nil, fmt.Errorf("%w: %w", reception_domain.ErrInternalDatabase, err)
	}

	return &created, nil
}

func (r *ReceptionPostgresRepository) FindByID(ctx context.Context, id string) (*reception_domain.Reception, error) {
	query := `
		SELECT id, date_time, pvz_id, status
		FROM avito.receptions
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	reception := reception_domain.Reception{}
	err := r.pool.QueryRow(ctx, query, args).Scan(
		&reception.ID,
		&reception.DateTime,
		&reception.PVZID,
		&reception.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %w", reception_domain.ErrReceptionNotFound, err)
		}
		return nil, fmt.Errorf("%w: %w", reception_domain.ErrInternalDatabase, err)
	}

	return &reception, nil
}

func (r *ReceptionPostgresRepository) FindLastOpenByPVZ(ctx context.Context, pvzID string) (*reception_domain.Reception, error) {
	query := `
		SELECT id, date_time, pvz_id, status
		FROM avito.receptions
		WHERE pvz_id = @pvz_id AND status = 'in_progress'
		ORDER BY date_time DESC
		LIMIT 1
	`

	args := pgx.NamedArgs{
		"pvz_id": pvzID,
	}

	reception := reception_domain.Reception{}

	err := r.pool.QueryRow(ctx, query, args).Scan(
		&reception.ID,
		&reception.DateTime,
		&reception.PVZID,
		&reception.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %w", reception_domain.ErrNoOpenReception, err)
		}
		return nil, fmt.Errorf("%w: %w", reception_domain.ErrInternalDatabase, err)
	}

	return &reception, nil
}

func (r *ReceptionPostgresRepository) CloseLastReception(ctx context.Context, pvzID string) (*reception_domain.Reception, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", reception_domain.ErrInternalDatabase, err)
	}
	defer tx.Rollback(ctx)

	query := `
		UPDATE avito.receptions
		SET status = 'close'
		WHERE id = (
			SELECT id
			FROM avito.receptions
			WHERE pvz_id = @pvz_id AND status = 'in_progress'
			ORDER BY date_time DESC
			LIMIT 1
		)
		RETURNING id, date_time, pvz_id, status
	`

	args := pgx.NamedArgs{
		"pvz_id": pvzID,
	}

	reception := reception_domain.Reception{}
	err = tx.QueryRow(ctx, query, args).Scan(
		&reception.ID,
		&reception.DateTime,
		&reception.PVZID,
		&reception.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %w", reception_domain.ErrReceptionNotFound, err)
		}
		return nil, fmt.Errorf("%w: %w", reception_domain.ErrInternalDatabase, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%w: %w", reception_domain.ErrInternalDatabase, err)
	}

	return &reception, nil
}

func (r *ReceptionPostgresRepository) ListByPVZ(ctx context.Context, pvzID string) ([]*reception_domain.Reception, error) {
	query := `
		SELECT id, date_time, pvz_id, status
		FROM avito.receptions
		WHERE pvz_id = @pvz_id
		ORDER BY date_time DESC
	`

	args := pgx.NamedArgs{
		"pvz_id": pvzID,
	}

	rows, err := r.pool.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", reception_domain.ErrInternalDatabase, err)
	}
	defer rows.Close()

	var receptions []*reception_domain.Reception
	for rows.Next() {
		var reception reception_domain.Reception
		err := rows.Scan(
			&reception.ID,
			&reception.DateTime,
			&reception.PVZID,
			&reception.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", reception_domain.ErrInternalDatabase, err)
		}
		receptions = append(receptions, &reception)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %w", reception_domain.ErrInternalDatabase, err)
	}

	return receptions, nil
}
