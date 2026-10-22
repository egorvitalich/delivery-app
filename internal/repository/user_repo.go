package repository

import (
	"context"
	"database/sql"
	"errors"

	"DostavkaProject/pkg/models"

	"github.com/lib/pq"
)

// UserRepo — репозиторий пользователей.
// Знает только про SQL, ничего не знает про HTTP, JWT, bcrypt.
type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, u *models.User) (int64, error) {
	// 1. Открываем транзакцию: два INSERT должны быть атомарны.
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var id int64
	err = tx.QueryRowContext(ctx,
		`INSERT INTO users (username, password_hash, email, name, second_name, roles)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING user_id`,
		u.Username,
		u.PasswordHash,
		u.Email,
		u.Name,
		u.SecondName,
		pq.Array(u.Roles), // []string -> text[]
	).Scan(&id)

	// 4. Обработка unique violation (SQLSTATE 23505).
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "users_email_key":
				return 0, ErrEmailTaken
			case "users_username_key":
				return 0, ErrUsernameTaken
			}
		}
		return 0, err
	}

	var insertSub string
	switch {
	case contains(u.Roles, "customer"):
		insertSub = `INSERT INTO customers (user_id) VALUES ($1)`
	case contains(u.Roles, "courier"):
		insertSub = `INSERT INTO couriers (user_id) VALUES ($1)`
	case contains(u.Roles, "manager"):
		insertSub = `INSERT INTO managers (user_id) VALUES ($1)`
	}
	if insertSub != "" {
		if _, err := tx.ExecContext(ctx, insertSub, id); err != nil {
			return 0, err
		}
	}

	// 6. Фиксируем транзакцию. После этого изменения в БД видны всем.
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}

// GetByEmail — для логина. Возвращает пользователя вместе с password_hash.
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	u := &models.User{}
	var roles []string

	err := r.db.QueryRowContext(ctx,
		`SELECT user_id, username, email, name, second_name, password_hash, roles
		 FROM users
		 WHERE email = $1`,
		email,
	).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Name,
		&u.SecondName,
		&u.PasswordHash, // <- нужен для bcrypt.CompareHashAndPassword
		pq.Array(&roles),
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	u.Roles = roles
	return u, nil
}

// GetByID — для роута /api/me. Без password_hash — он не нужен.
func (r *UserRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
	u := &models.User{}
	var roles []string

	err := r.db.QueryRowContext(ctx,
		`SELECT user_id, username, email, name, second_name, roles
		 FROM users
		 WHERE user_id = $1`,
		id,
	).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Name,
		&u.SecondName,
		pq.Array(&roles),
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	u.Roles = roles
	return u, nil
}

// contains — хелпер: есть ли строка в слайсе.
// Можно вынести в отдельный util-пакет, но пока пусть живёт тут.
func contains(s []string, v string) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}
