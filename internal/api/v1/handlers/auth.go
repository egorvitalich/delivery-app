package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"DostavkaProject/internal/repository"
	"DostavkaProject/internal/service"
	"DostavkaProject/pkg/models"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid json")
		return
	}

	resp, err := h.svc.Register(r.Context(), req)
	switch {
	case errors.Is(err, repository.ErrEmailTaken):
		writeErr(w, http.StatusConflict, "email already used")
	case errors.Is(err, repository.ErrUsernameTaken):
		writeErr(w, http.StatusConflict, "username already used")
	case errors.Is(err, service.ErrWeakPassword):
		writeErr(w, http.StatusBadRequest, "password min 8 chars")
	case errors.Is(err, service.ErrInvalidEmail):
		writeErr(w, http.StatusBadRequest, "invalid email")
	case errors.Is(err, service.ErrInvalidUsername):
		writeErr(w, http.StatusBadRequest, "username min 3 chars")
	case err != nil:
		writeErr(w, http.StatusInternalServerError, "internal error")
	default:
		writeJSON(w, http.StatusCreated, resp)
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid json")
		return
	}

	resp, err := h.svc.Login(r.Context(), req)
	if errors.Is(err, service.ErrInvalidCredentials) {
		writeErr(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	id, ok := UserIDFromCtx(r.Context())
	if !ok {
		writeErr(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	u, err := h.svc.Me(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		writeErr(w, http.StatusNotFound, "user not found")
		return
	}
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, u)
}
