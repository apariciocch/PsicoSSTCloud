package handler

import (
	"encoding/json"
	"net/http"

	"github.com/apariciocch/psicosstcloud/internal/dto"
	"github.com/apariciocch/psicosstcloud/internal/middleware"
	"github.com/apariciocch/psicosstcloud/internal/pkg/response"
	"github.com/apariciocch/psicosstcloud/internal/service"
)

// AuthHandler maneja operaciones de autenticación
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler crea un nuevo handler de autenticación
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login maneja el login del usuario
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request format")
		return
	}

	// Validar entrada
	if req.Email == "" {
		response.ErrorWithField(w, http.StatusBadRequest, "VALIDATION_ERROR", "Email is required", "email")
		return
	}

	if req.Password == "" {
		response.ErrorWithField(w, http.StatusBadRequest, "VALIDATION_ERROR", "Password is required", "password")
		return
	}

	// Realizar login
	loginResp, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Invalid credentials")
		return
	}

	response.Success(w, http.StatusOK, loginResp)
}

// Logout maneja el logout del usuario
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized")
		return
	}

	err = h.authService.Logout(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "SYSTEM_001", "Error processing logout")
		return
	}

	response.Success(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

// RefreshToken maneja la renovación del token
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request format")
		return
	}

	if req.RefreshToken == "" {
		response.ErrorWithField(w, http.StatusBadRequest, "VALIDATION_ERROR", "Refresh token is required", "refresh_token")
		return
	}

	loginResp, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "AUTH_002", "Invalid refresh token")
		return
	}

	response.Success(w, http.StatusOK, loginResp)
}

// ChangePassword maneja el cambio de contraseña
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized")
		return
	}

	var req dto.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request format")
		return
	}

	// Validar entrada
	if req.CurrentPassword == "" || req.NewPassword == "" {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Missing required fields")
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		response.ErrorWithField(w, http.StatusBadRequest, "VALIDATION_ERROR", "New passwords don't match", "confirm_password")
		return
	}

	err = h.authService.ChangePassword(r.Context(), userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid current password")
		return
	}

	response.Success(w, http.StatusOK, map[string]string{"message": "Password changed successfully"})
}
