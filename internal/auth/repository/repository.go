package repository

import (
	"context"

	"github.com/Prizze/TaskScheduler/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepoAuth struct {
	pool *pgxpool.Pool
}

func NewRepoAuth(pool *pgxpool.Pool) *RepoAuth {
	return &RepoAuth{
		pool: pool,
	}
}

func (r *RepoAuth) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	var user models.User
	err := r.pool.QueryRow(ctx, getUserByID, userID).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *RepoAuth) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.pool.QueryRow(ctx, getUserByEmail, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *RepoAuth) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	var createdUser models.User
	err := r.pool.QueryRow(ctx, createUser, user.Email, user.PasswordHash, user.Name).Scan(
		&createdUser.ID,
		&createdUser.Email,
		&createdUser.PasswordHash,
		&createdUser.Name,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func (r *RepoAuth) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	var id int64
	err := r.pool.QueryRow(ctx, checkEmail, email).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
