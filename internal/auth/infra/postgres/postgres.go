package postgres

import (
	"context"
	"errors"
	"fmt"

	auth_domain "github.com/0x0FACED/pvz-avito/internal/auth/domain"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (r *AuthPostgresRepository) Create(ctx context.Context, user *auth_domain.User) (*auth_domain.User, error) {
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

	var created auth_domain.User
	created.Password = user.Password

	err := r.pool.QueryRow(ctx, query, args).Scan(
		&created.ID, &created.Email, &created.Role,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr); pgErr.Code == "23505" {
			return nil, fmt.Errorf("%w: %w", auth_domain.ErrUserAlreadyExists, err)
		}
		// in openapi wrote that 201 and 400 codes only.
		// so using domain.ErrInternalDatabase as invalid request (400)
		return nil, fmt.Errorf("%w: %w", auth_domain.ErrInternalDatabase, err)
	}

	return &created, nil
}

func (r *AuthPostgresRepository) FindByEmail(ctx context.Context, email string) (*auth_domain.User, error) {
	query := `
		SELECT id, email, password_hash, role
		FROM avito.users
		WHERE email = @email
	`

	args := pgx.NamedArgs{
		"email": email,
	}

	user := auth_domain.User{}

	err := r.pool.QueryRow(ctx, query, args).Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %w", auth_domain.ErrUserNotFound, err)
		}
		return nil, fmt.Errorf("%w: %w", auth_domain.ErrInternalDatabase, err)
	}

	return &user, nil
}
