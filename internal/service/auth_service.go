package service

import (
	"context"
	"errors"
	"strings"

	"DostavkaProject/internal/repository"
	"DostavkaProject/pkg/models"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrWeakPassword       = errors.New("password too weak")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidUsername    = errors.New("invalid username")
)

type AuthService struct {
	repo *repository.UserRepo
	jwt  *JWTService
}

func NewAuthService(repo *repository.UserRepo, jwt *JWTService) *AuthService {
	return &AuthService{repo: repo, jwt: jwt}
}

func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if !validEmail(req.Email) {
		return nil, ErrInvalidEmail
	}
	if len(req.Password) < 8 {
		return nil, ErrWeakPassword
	}
	if len(req.Username) < 3 {
		return nil, ErrInvalidUsername
	}

	switch req.Role {
	case "customer", "courier", "manager":
	default:
		req.Role = "customer"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, err
	}

	u := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		Name:         req.Name,
		SecondName:   req.SecondName,
		PasswordHash: string(hash),
		Roles:        []string{req.Role},
	}

	id, err := s.repo.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}
	u.ID = id

	token, err := s.jwt.Generate(u)
	if err != nil {
		return nil, err
	}
	return &models.AuthResponse{Token: token, UserID: id, Roles: u.Roles}, nil
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials // не палим существование email-а
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)) != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.jwt.Generate(u)
	if err != nil {
		return nil, err
	}
	return &models.AuthResponse{Token: token, UserID: u.ID, Roles: u.Roles}, nil
}

func (s *AuthService) Me(ctx context.Context, id int64) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

func validEmail(s string) bool {
	at := strings.IndexByte(s, '@')
	return at > 0 && at < len(s)-1 && strings.Contains(s[at:], ".")
}
