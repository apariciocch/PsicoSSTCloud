package middleware
package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/apariciocch/psicosstcloud/internal/service/jwt"
	"github.com/apariciocch/psicosstcloud/internal/pkg/response"
)

const (
	UserIDKey = "user_id"
	RoleIDKey = "role_id"
	EmailKey  = "email"
	PermissionsKey = "permissions"
	IsSuperAdminKey = "is_super_admin"
)

// JWTAuth middleware para validar JWT
func JWTAuth(jwtManager *jwt.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Obtener token del header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Error(w, http.StatusUnauthorized, "AUTH_001", "Token no proporcionado")
				return
			}

			// Extraer token del formato "Bearer <token>"
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.Error(w, http.StatusUnauthorized, "AUTH_003", "Formato de token inválido")
				return
			}

			tokenString := parts[1]

			// Validar token
			claims, err := jwtManager.ValidateToken(tokenString)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "AUTH_002", "Token expirado o inválido")
				return
			}

			// Añadir datos del usuario al contexto
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, RoleIDKey, claims.RoleID)
			ctx = context.WithValue(ctx, EmailKey, claims.Email)
			ctx = context.WithValue(ctx, PermissionsKey, claims.Permissions)
			ctx = context.WithValue(ctx, IsSuperAdminKey, claims.IsSuperAdmin)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// CheckPermission middleware para verificar permisos
func CheckPermission(requiredPermissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Super admin tiene todos los permisos
			isSuperAdmin, ok := r.Context().Value(IsSuperAdminKey).(bool)
			if ok && isSuperAdmin {
				next.ServeHTTP(w, r)
				return
			}

			// Verificar permisos
			permissionsInterface := r.Context().Value(PermissionsKey)
			permissions, ok := permissionsInterface.([]string)
			if !ok {
				response.Error(w, http.StatusForbidden, "AUTH_004", "Acceso denegado")
				return
			}

			hasPermission := false
			for _, required := range requiredPermissions {
				for _, userPerm := range permissions {
					if userPerm == required {
						hasPermission = true
						break
					}
				}
				if hasPermission {
					break
				}
			}

			if !hasPermission {
				response.Error(w, http.StatusForbidden, "AUTH_004", "No tiene permisos suficientes")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetUserIDFromContext obtiene el ID de usuario del contexto
func GetUserIDFromContext(r *http.Request) (string, error) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

// GetUserFromContext obtiene datos de usuario del contexto
func GetUserFromContext(r *http.Request) map[string]interface{} {
	return map[string]interface{}{
		"user_id": r.Context().Value(UserIDKey),
		"role_id": r.Context().Value(RoleIDKey),
		"email": r.Context().Value(EmailKey),
		"is_super_admin": r.Context().Value(IsSuperAdminKey),
	}
}
