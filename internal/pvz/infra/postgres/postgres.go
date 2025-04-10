package postgres

import (
	"context"
	"time"

	product_domain "github.com/0x0FACED/pvz-avito/internal/product/domain"
	pvz_domain "github.com/0x0FACED/pvz-avito/internal/pvz/domain"
	reception_domain "github.com/0x0FACED/pvz-avito/internal/reception/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PVZPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPVZPostgresRepository(pgx *pgxpool.Pool) *PVZPostgresRepository {
	return &PVZPostgresRepository{pool: pgx}
}

func (r *PVZPostgresRepository) Create(ctx context.Context, pvz *pvz_domain.PVZ) (*pvz_domain.PVZ, error) {
	query := `
		INSERT INTO avito.pvz (id, registration_date, city)
		VALUES (@id, @registration_date, @city)
		RETURNING id, registration_date, city
	`

	args := pgx.NamedArgs{
		"id":                pvz.ID,
		"registration_date": pvz.RegistrationDate,
		"city":              pvz.City,
	}

	err := r.pool.QueryRow(ctx, query, args).Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City)
	if err != nil {
		// TODO: handle err
		return nil, err
	}

	return pvz, nil
}

func (r *PVZPostgresRepository) FindByID(ctx context.Context, id string) (*pvz_domain.PVZ, error) {
	query := `
		SELECT id, registration_date, city
		FROM avito.pvz
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	pvz := pvz_domain.PVZ{}

	err := r.pool.QueryRow(ctx, query, args).Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City)
	if err != nil {
		// TODO: handle err
		return nil, err
	}

	return &pvz, nil
}

func (r *PVZPostgresRepository) ListWithReceptions(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]*pvz_domain.PVZWithReceptions, error) {
	query := `
		SELECT p.id, p.registration_date, p.city,
		       r.id, r.date_time, r.pvz_id, r.status,
		       pr.id, pr.date_time, pr.type, pr.reception_id
		FROM avito.pvz p
		LEFT JOIN avito.receptions r ON r.pvz_id = p.id
		LEFT JOIN avito.products pr ON pr.reception_id = r.id
		WHERE (r.date_time BETWEEN @start_date AND @end_date OR @start_date IS NULL)
		ORDER BY p.registration_date DESC, r.date_time DESC
		LIMIT @limit OFFSET @offset
	`

	offset := (page - 1) * limit
	args := pgx.NamedArgs{
		"start_date": startDate,
		"end_date":   endDate,
		"limit":      limit,
		"offset":     offset,
	}

	rows, err := r.pool.Query(ctx, query, args)
	if err != nil {
		// TODO: handle err
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]*pvz_domain.PVZWithReceptions)
	receptionsMap := make(map[string]*pvz_domain.ReceptionWithProducts)

	for rows.Next() {
		var (
			pvzID            string
			pvzRegDate       time.Time
			pvzCity          pvz_domain.City
			receptionID      *string
			receptionDT      *time.Time
			receptionPVZID   *string
			receptionStatus  *reception_domain.Status
			productID        *string
			productDT        *time.Time
			productType      *product_domain.ProductType
			productReception *string
		)

		err := rows.Scan(
			&pvzID, &pvzRegDate, &pvzCity,
			&receptionID, &receptionDT, &receptionPVZID, &receptionStatus,
			&productID, &productDT, &productType, &productReception,
		)
		if err != nil {
			// TODO: handle err
			return nil, err
		}

		if _, ok := result[pvzID]; !ok {
			result[pvzID] = &pvz_domain.PVZWithReceptions{
				PVZ: &pvz_domain.PVZ{
					ID:               pvzID,
					RegistrationDate: pvzRegDate,
					City:             pvz_domain.City(pvzCity),
				},
				Receptions: []*pvz_domain.ReceptionWithProducts{},
			}
		}

		if receptionID != nil {
			receptionKey := *receptionID
			if _, ok := receptionsMap[receptionKey]; !ok {
				receptionsMap[receptionKey] = &pvz_domain.ReceptionWithProducts{
					Reception: &reception_domain.Reception{
						ID:       *receptionID,
						DateTime: *receptionDT,
						PVZID:    *receptionPVZID,
						Status:   *receptionStatus,
					},
					Products: []*product_domain.Product{},
				}
				result[pvzID].Receptions = append(result[pvzID].Receptions, receptionsMap[receptionKey])
			}

			if productID != nil {
				receptionsMap[receptionKey].Products = append(receptionsMap[receptionKey].Products, &product_domain.Product{
					ID:          *productID,
					DateTime:    *productDT,
					Type:        *productType,
					ReceptionID: *productReception,
				})
			}
		}
	}

	var finalResult []*pvz_domain.PVZWithReceptions
	for _, v := range result {
		finalResult = append(finalResult, v)
	}

	return finalResult, nil
}
