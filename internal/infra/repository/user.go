package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leonardonicola/golerplate/internal/domain/entity"
	"github.com/leonardonicola/golerplate/pkg/constants"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByCPF(ctx context.Context, cpf string) (*entity.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Check email
	exists, err := r.emailExists(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New(constants.ErrMsgEmailInUse)
	}

	// Check CPF
	exists, err = r.cpfExists(ctx, user.CPF)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New(constants.ErrMsgCPFInUse)
	}

	query := `
    INSERT INTO users (id, full_name, email, cpf, age, password)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING id, full_name, email, cpf, age
  `

	err = tx.QueryRow(ctx, query, uuid.NewString(), user.FullName, user.Email, user.CPF, user.Age, user.Password).Scan(&user.ID, &user.FullName, &user.Email, &user.CPF, &user.Age)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	user := &entity.User{}

	query := `
    SELECT id, full_name, email, cpf, age, created_at, updated_at
    FROM users
    WHERE id = $1 AND deleted_at IS NULL
  `

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.CPF,
		&user.Age,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.New(constants.ErrMsgUserNotFound)
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := &entity.User{}

	query := `
    SELECT id, full_name, email, cpf, age, password, created_at, updated_at
    FROM users
    WHERE email = $1 AND deleted_at IS NULL
  `

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.CPF,
		&user.Age,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.New(constants.ErrMsgUserNotFound)
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByCPF(ctx context.Context, cpf string) (*entity.User, error) {
	user := &entity.User{}

	query := `
    SELECT id, full_name, email, cpf, age, password, created_at, updated_at
    FROM users
    WHERE cpf = $1 AND deleted_at IS NULL
  `

	err := r.db.QueryRow(ctx, query, cpf).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.CPF,
		&user.Age,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.New(constants.ErrMsgUserNotFound)
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) emailExists(ctx context.Context, email string) (bool, error) {
	var count int
	query := `
    SELECT COUNT(*) FROM users 
    WHERE email = $1 AND deleted_at IS NULL
  `

	err := r.db.QueryRow(ctx, query, email).Scan(&count)
	return count > 0, err
}

func (r *userRepository) cpfExists(ctx context.Context, cpf string) (bool, error) {
	var count int
	query := `
    SELECT COUNT(*) FROM users 
    WHERE cpf = $1 AND deleted_at IS NULL
  `

	err := r.db.QueryRow(ctx, query, cpf).Scan(&count)
	return count > 0, err
}
