package service

import (
	"errors"
	"strconv"
	"time"

	"DostavkaProject/pkg/models"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int64    `json:"uid"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTService(secret []byte, ttl time.Duration) *JWTService {
	return &JWTService{secret: secret, ttl: ttl}
}

func (s *JWTService) Generate(u *models.User) (string, error) {
	claims := Claims{
		UserID: u.ID,
		Roles:  u.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "dostavka",
			Subject:   strconv.FormatInt(u.ID, 10),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
}

func (s *JWTService) Parse(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	t, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil || !t.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
