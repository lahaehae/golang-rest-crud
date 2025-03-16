package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lahaehae/crud_project/internal/models"
)

type SubUserRepository struct {
    pool *pgxpool.Pool
}

func NewSubUserRepository(pool *pgxpool.Pool) *SubUserRepository {
    return &SubUserRepository{
        pool: pool,
    }
}

func (r *SubUserRepository) Create(ctx context.Context, subUser *models.SubUser) error {
    query := `
        INSERT INTO sub_users (id, owner_id, name, email, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
    _, err := r.pool.Exec(ctx, query,
        subUser.ID,
        subUser.OwnerID,
        subUser.Name,
        subUser.Email,
        subUser.CreatedAt,
        subUser.UpdatedAt,
    )
    return err
}

func (r *SubUserRepository) GetByOwner(ctx context.Context, ownerID string) ([]models.SubUser, error) {
    query := `
        SELECT id, owner_id, name, email, created_at, updated_at
        FROM sub_users
        WHERE owner_id = $1
    `
    rows, err := r.pool.Query(ctx, query, ownerID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var subUsers []models.SubUser
    for rows.Next() {
        var su models.SubUser
        err := rows.Scan(
            &su.ID,
            &su.OwnerID,
            &su.Name,
            &su.Email,
            &su.CreatedAt,
            &su.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        subUsers = append(subUsers, su)
    }
    return subUsers, nil
}

// Другие методы CRUD...