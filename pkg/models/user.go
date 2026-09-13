package models

import "time"

// User — доменная модель пользователя. Соответствует таблице users.
// Тег json:"-" на PasswordHash означает: никогда не отдавать хеш в JSON.
type User struct {
	ID           int64     `json:"id"           db:"user_id"`
	Username     string    `json:"username"     db:"username"`
	Email        string    `json:"email"        db:"email"`
	Name         string    `json:"name"         db:"name"`
	SecondName   string    `json:"second_name"  db:"second_name"`
	PasswordHash string    `json:"-"            db:"password_hash"`
	Roles        []string  `json:"roles"        db:"roles"`
	CreatedAt    time.Time `json:"created_at,omitempty" db:"created_at"`
}

// RegisterRequest — тело POST /api/register.
type RegisterRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	Username   string `json:"username"`
	Name       string `json:"name"`
	SecondName string `json:"second_name"`
	Role       string `json:"role"` // customer | courier | manager
}

// LoginRequest — тело POST /api/login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse — ответ на успешную регистрацию/логин.
type AuthResponse struct {
	Token  string   `json:"token"`
	UserID int64    `json:"user_id"`
	Roles  []string `json:"roles"`
}
