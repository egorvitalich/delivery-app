package v1

import (
	"net/http"

	"DostavkaProject/internal/api/v1/handlers"
	"DostavkaProject/internal/service"
)

// Deps — всё, что нужно роутеру.
type Deps struct {
	Auth    *handlers.AuthHandler
	Version *handlers.VersionHandler
	JWT     *service.JWTService
}

func NewRouter(d Deps) http.Handler {
	mux := http.NewServeMux()

	// --- публичные ---
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})
	mux.HandleFunc("GET /api/version", d.Version.Get)
	mux.HandleFunc("POST /api/register", d.Auth.Register)
	mux.HandleFunc("POST /api/login", d.Auth.Login)

	// --- защищённые ---
	auth := handlers.AuthMiddleware(d.JWT)
	mux.Handle("GET /api/me", auth(http.HandlerFunc(d.Auth.Me)))

	// Пример для будущих роутов менеджера:
	// managerOnly := handlers.RequireRole("manager")
	// mux.Handle("GET /api/manager/orders",
	//     auth(managerOnly(http.HandlerFunc(someHandler))))

	return mux
}
