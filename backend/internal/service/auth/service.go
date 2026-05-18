package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/apariciocch/psicosstcloud/internal/config"
	"github.com/apariciocch/psicosstcloud/internal/domain/entities"
	"github.com/apariciocch/psicosstcloud/internal/dto"
	"github.com/apariciocch/psicosstcloud/internal/repository"
	"github.com/apariciocch/psicosstcloud/internal/service/jwt"
	"github.com/apariciocch/psicosstcloud/internal/service/password"
)

// AuthService implementación de autenticación
type AuthService struct {
	userRepo     repository.UserRepository
	roleRepo     repository.RoleRepository
	jwtManager   *jwt.Manager
	passwordMgr  *password.Manager
	auditService AuditService
}

// AuditService para logs de auditoría
type AuditService interface {
	Log(ctx context.Context, audit *entities.AuditLog) error
}

// UserRepository para acceso de usuarios
type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetByID(ctx context.Context, id string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
}

// RoleRepository para acceso de roles
type RoleRepository interface {
	GetByID(ctx context.Context, id string) (*entities.Role, error)
}

// NewAuthService crea un nuevo servicio de autenticación
func NewAuthService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	jwtManager *jwt.Manager,
	passwordMgr *password.Manager,
	auditService AuditService,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		roleRepo:     roleRepo,
		jwtManager:   jwtManager,
		passwordMgr:  passwordMgr,
		auditService: auditService,
	}
}

// Login autentica un usuario
func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Buscar usuario
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Verificar si el usuario está activo
	if !user.IsActive {
		return nil, fmt.Errorf("user account is inactive")
	}

	// Verificar si el usuario está bloqueado
	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		return nil, fmt.Errorf("user account is temporarily locked")
	}

	// Verificar contraseña
	if !s.passwordMgr.Compare(user.PasswordHash, req.Password) {
		// Incrementar intentos fallidos
		user.LoginAttempts++
		if user.LoginAttempts >= config.MaxLoginAttempts {
			lockTime := time.Now().Add(time.Duration(30) * time.Minute)
			user.LockedUntil = &lockTime
		}
		_, _ = s.userRepo.Update(ctx, user.ID, user)

		// Auditoría
		s.logFailedLogin(ctx, user.Email, "Invalid password")

		return nil, fmt.Errorf("invalid password")
	}

	// Resetear intentos fallidos
	user.LoginAttempts = 0
	user.LockedUntil = nil
	now := time.Now()
	user.LastLogin = &now
	_, _ = s.userRepo.Update(ctx, user.ID, user)

	// Obtener rol
	role, err := s.roleRepo.GetByID(ctx, user.RoleID)
	if err != nil {
		return nil, fmt.Errorf("role not found")
	}

	// Generar tokens
	accessToken, expiresIn, err := s.jwtManager.GenerateToken(
		user.ID,
		user.Email,
		user.RoleID,
		role.Permissions,
		user.IsSuperAdmin,
	)
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %w", err)
	}

	refreshToken, _, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("error generating refresh token: %w", err)
	}

	// Auditoría
	s.logSuccessLogin(ctx, user.ID, user.Email)

	// Preparar respuesta
	userResp := &dto.UserResponse{
		ID:           user.ID,
		Email:        user.Email,
		Username:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		RoleID:       user.RoleID,
		IsActive:     user.IsActive,
		IsSuperAdmin: user.IsSuperAdmin,
		Phone:        user.Phone,
		Department:   user.Department,
		LastLogin:    user.LastLogin,
		MFAEnabled:   user.MFAEnabled,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		User:         userResp,
	}, nil
}

// Logout cierra la sesión del usuario
func (s *AuthService) Logout(ctx context.Context, userID string) error {
	// Auditoría
	audit := &entities.AuditLog{
		UserID:     userID,
		Action:     config.AuditActionLogout,
		EntityType: "auth",
		EntityID:   userID,
		Status:     config.AuditStatusSuccess,
		CreatedAt:  time.Now(),
	}
	return s.auditService.Log(ctx, audit)
}

// RefreshToken renueva el access token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.LoginResponse, error) {
	// Validar refresh token
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Obtener usuario
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("user account is inactive")
	}

	// Obtener rol
	role, err := s.roleRepo.GetByID(ctx, user.RoleID)
	if err != nil {
		return nil, fmt.Errorf("role not found")
	}

	// Generar nuevo access token
	accessToken, expiresIn, err := s.jwtManager.GenerateToken(
		user.ID,
		user.Email,
		user.RoleID,
		role.Permissions,
		user.IsSuperAdmin,
	)
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %w", err)
	}

	// Generar nuevo refresh token
	newRefreshToken, _, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("error generating refresh token: %w", err)
	}

	userResp := &dto.UserResponse{
		ID:           user.ID,
		Email:        user.Email,
		Username:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		RoleID:       user.RoleID,
		IsActive:     user.IsActive,
		IsSuperAdmin: user.IsSuperAdmin,
		Phone:        user.Phone,
		Department:   user.Department,
		LastLogin:    user.LastLogin,
		MFAEnabled:   user.MFAEnabled,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    expiresIn,
		User:         userResp,
	}, nil
}

// ChangePassword cambia la contraseña del usuario
func (s *AuthService) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	// Obtener usuario
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Verificar contraseña actual
	if !s.passwordMgr.Compare(user.PasswordHash, currentPassword) {
		return fmt.Errorf("current password is incorrect")
	}

	// Validar nueva contraseña
	valid, errors := s.passwordMgr.ValidateStrength(newPassword)
	if !valid {
		return fmt.Errorf("weak password: %v", errors)
	}

	// Hash de nueva contraseña
	newHash, err := s.passwordMgr.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	// Actualizar usuario
	user.PasswordHash = newHash
	if _, err := s.userRepo.Update(ctx, user.ID, user); err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	// Auditoría
	audit := &entities.AuditLog{
		UserID:        userID,
		Action:        "change_password",
		EntityType:    "user",
		EntityID:      userID,
		ChangeSummary: stringPtr("Password changed"),
		Status:        config.AuditStatusSuccess,
		CreatedAt:     time.Now(),
	}
	_ = s.auditService.Log(ctx, audit)

	return nil
}

// ValidateToken valida un token y retorna sus claims
func (s *AuthService) ValidateToken(tokenString string) (map[string]interface{}, error) {
	claims, err := s.jwtManager.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"user_id":        claims.UserID,
		"email":          claims.Email,
		"role_id":        claims.RoleID,
		"permissions":    claims.Permissions,
		"is_super_admin": claims.IsSuperAdmin,
	}, nil
}

// Helper methods
func (s *AuthService) logFailedLogin(ctx context.Context, email, reason string) {
	audit := &entities.AuditLog{
		Action:       config.AuditActionLogin,
		EntityType:   "auth",
		EntityID:     email,
		Status:       config.AuditStatusFailure,
		ErrorMessage: &reason,
		CreatedAt:    time.Now(),
	}
	_ = s.auditService.Log(ctx, audit)
}

func (s *AuthService) logSuccessLogin(ctx context.Context, userID, email string) {
	audit := &entities.AuditLog{
		UserID:     userID,
		Action:     config.AuditActionLogin,
		EntityType: "auth",
		EntityID:   userID,
		Status:     config.AuditStatusSuccess,
		CreatedAt:  time.Now(),
	}
	_ = s.auditService.Log(ctx, audit)
}

func stringPtr(s string) *string {
	return &s
}
