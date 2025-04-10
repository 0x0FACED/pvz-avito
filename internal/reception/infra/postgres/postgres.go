package postgres

import (
	"context"

	"github.com/0x0FACED/pvz-avito/internal/reception/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReceptionPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewReceptionPostgresRepository(pgx *pgxpool.Pool) *ReceptionPostgresRepository {
	return &ReceptionPostgresRepository{pool: pgx}
}

func (r *ReceptionPostgresRepository) Create(ctx context.Context, reception *domain.Reception) error {
	query := `
		INSERT INTO avito.receptions (pvz_id, status)
		VALUES (@pvz_id, 'in_progress')
		RETURNING id, date_time
	`

	args := pgx.NamedArgs{
		"pvz_id": reception.PVZID,
	}

	err := r.pool.QueryRow(ctx, query, args).Scan(&reception.ID, &reception.DateTime)
	if err != nil {
		return err
	}

	reception.Status = "in_progress"
	return nil
}

func (r *ReceptionPostgresRepository) FindByID(ctx context.Context, id string) (*domain.Reception, error) {
	query := `
		SELECT id, date_time, pvz_id, status
		FROM avito.receptions
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	reception := domain.Reception{}
	err := r.pool.QueryRow(ctx, query, args).Scan(
		&reception.ID,
		&reception.DateTime,
		&reception.PVZID,
		&reception.Status,
	)
	if err != nil {
		return nil, err
	}

	return &reception, nil
}

func (r *ReceptionPostgresRepository) FindLastOpenByPVZ(ctx context.Context, pvzID string) (*domain.Reception, error) {
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

	reception := domain.Reception{}
	err := r.pool.QueryRow(ctx, query, args).Scan(
		&reception.ID,
		&reception.DateTime,
		&reception.PVZID,
		&reception.Status,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &reception, nil
}

func (r *ReceptionPostgresRepository) CloseLastReception(ctx context.Context, pvzID string) (*domain.Reception, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
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

	reception := domain.Reception{}
	err = tx.QueryRow(ctx, query, args).Scan(
		&reception.ID,
		&reception.DateTime,
		&reception.PVZID,
		&reception.Status,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &reception, nil
}

func (r *ReceptionPostgresRepository) ListByPVZ(ctx context.Context, pvzID string) ([]*domain.Reception, error) {
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
		return nil, err
	}
	defer rows.Close()

	var receptions []*domain.Reception
	for rows.Next() {
		var reception domain.Reception
		err := rows.Scan(
			&reception.ID,
			&reception.DateTime,
			&reception.PVZID,
			&reception.Status,
		)
		if err != nil {
			return nil, err
		}
		receptions = append(receptions, &reception)
	}

	return receptions, nil
}
