package models

import "time"

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

type RegisterRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	Username   string `json:"username"`
	Name       string `json:"name"`
	SecondName string `json:"second_name"`
	Role       string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token  string   `json:"token"`
	UserID int64    `json:"user_id"`
	Roles  []string `json:"roles"`
}
