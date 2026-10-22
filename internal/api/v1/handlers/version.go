package handlers

import (
	"database/sql"
	"net/http"
)

type VersionHandler struct {
	db *sql.DB
}

func NewVersionHandler(db *sql.DB) *VersionHandler {
	return &VersionHandler{db: db}
}

func (h *VersionHandler) Get(w http.ResponseWriter, r *http.Request) {
	var version string
	err := h.db.QueryRow("SELECT version FROM application_info LIMIT 1").Scan(&version)
	if err != nil {
		writeErr(w, http.StatusNotFound, "Версия не найдена")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"version": version})
}
