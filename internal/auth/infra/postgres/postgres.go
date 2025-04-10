package postgres

import (
	"context"

	"github.com/0x0FACED/pvz-avito/internal/auth/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewAuthPostgresRepository(pgx *pgxpool.Pool) *AuthPostgresRepository {
	return &AuthPostgresRepository{
		pool: pgx,
	}
}

func (r *AuthPostgresRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `
		INSERT INTO avito.users (email, password_hash, role)
		VALUES (@email, @password_hash, @role)
		RETURNING id, email, role
	`

	args := pgx.NamedArgs{
		"email":         user.Email,
		"password_hash": user.Password,
		"role":          user.Role,
	}

	var created domain.User
	created.Password = user.Password

	err := r.pool.QueryRow(ctx, query, args).Scan(
		&created.ID, &created.Email, &created.Role,
	)
	if err != nil {
		// TODO: handle err
		return nil, err
	}

	return &created, nil
}

func (r *AuthPostgresRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, role
		FROM avito.users
		WHERE email = @email
	`

	args := pgx.NamedArgs{
		"email": email,
	}

	user := domain.User{}

	err := r.pool.QueryRow(ctx, query, args).Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		// TODO: handle err
		return nil, err
	}

	return &user, nil
}
