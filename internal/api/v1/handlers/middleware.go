package handlers

import (
	"context"
	"net/http"
	"strings"

	"DostavkaProject/internal/service"
)

type ctxKey string

const (
	ctxUserID ctxKey = "uid"
	ctxRoles  ctxKey = "roles"
)

// AuthMiddleware проверяет Bearer-токен и кладёт uid/roles в context запроса.
func AuthMiddleware(jwt *service.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hdr := r.Header.Get("Authorization")
			if !strings.HasPrefix(hdr, "Bearer ") {
				writeErr(w, http.StatusUnauthorized, "missing token")
				return
			}
			claims, err := jwt.Parse(strings.TrimPrefix(hdr, "Bearer "))
			if err != nil {
				writeErr(w, http.StatusUnauthorized, "invalid token")
				return
			}
			ctx := context.WithValue(r.Context(), ctxUserID, claims.UserID)
			ctx = context.WithValue(ctx, ctxRoles, claims.Roles)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole пускает дальше только если у пользователя есть хотя бы одна из ролей.
// Использовать ТОЛЬКО после AuthMiddleware.
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRoles, _ := r.Context().Value(ctxRoles).([]string)
			for _, need := range roles {
				for _, have := range userRoles {
					if need == have {
						next.ServeHTTP(w, r)
						return
					}
				}
			}
			writeErr(w, http.StatusForbidden, "forbidden")
		})
	}
}

// UserIDFromCtx достаёт uid, положенный AuthMiddleware.
func UserIDFromCtx(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(ctxUserID).(int64)
	return id, ok
}
