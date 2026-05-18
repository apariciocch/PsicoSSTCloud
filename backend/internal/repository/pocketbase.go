package repository

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/apariciocch/psicosstcloud/internal/domain/entities"
	"github.com/pocketbase/pocketbase"
)

// UserRepository interfaz para operaciones con usuarios
type UserRepository interface {
	GetByID(ctx context.Context, id string) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetAll(ctx context.Context, page, pageSize int) ([]*entities.User, int, error)
	Create(ctx context.Context, user *entities.User) (*entities.User, error)
	Update(ctx context.Context, id string, user *entities.User) (*entities.User, error)
	Delete(ctx context.Context, id string) error
}

// RoleRepository interfaz para operaciones con roles
type RoleRepository interface {
	GetByID(ctx context.Context, id string) (*entities.Role, error)
	GetAll(ctx context.Context, page, pageSize int) ([]*entities.Role, int, error)
	Create(ctx context.Context, role *entities.Role) (*entities.Role, error)
	Update(ctx context.Context, id string, role *entities.Role) (*entities.Role, error)
	Delete(ctx context.Context, id string) error
}

// PocketBaseClient cliente para PocketBase
type PocketBaseClient struct {
	url       string
	authToken string
	http      *http.Client
}

// NewPocketBaseClient crea un nuevo cliente de PocketBase
func NewPocketBaseClient(pbURL, adminEmail, adminPassword string) (*PocketBaseClient, error) {
	client := &PocketBaseClient{
		url:  pbURL,
		http: &http.Client{},
	}

	// Por ahora, solo retornamos el cliente sin autenticación
	// La autenticación se implementará cuando se necesite
	return client, nil
}

// GetRecord obtiene un registro
func (c *PocketBaseClient) GetRecord(ctx context.Context, collection, id string) (map[string]interface{}, error) {
	// Implementación con PocketBase SDK
	// Será llenada cuando se use el SDK real
	return nil, nil
}

// GetList obtiene una lista de registros
func (c *PocketBaseClient) GetList(ctx context.Context, collection string, page, pageSize int, filter string) ([]map[string]interface{}, int, error) {
	// Implementación con PocketBase SDK
	return nil, 0, nil
}

// Create crea un nuevo registro
func (c *PocketBaseClient) Create(ctx context.Context, collection string, data map[string]interface{}) (map[string]interface{}, error) {
	// Implementación con PocketBase SDK
	return nil, nil
}

// Update actualiza un registro
func (c *PocketBaseClient) Update(ctx context.Context, collection, id string, data map[string]interface{}) (map[string]interface{}, error) {
	// Implementación con PocketBase SDK
	return nil, nil
}

// Delete elimina un registro
func (c *PocketBaseClient) Delete(ctx context.Context, collection, id string) error {
	// Implementación con PocketBase SDK
	return nil
}

// GetApp retorna la app de PocketBase para acceso directo
// NOTA: En la versión actual, retorna nil ya que se usa cliente HTTP
func (c *PocketBaseClient) GetApp() *pocketbase.PocketBase {
	return nil
}

// Close cierra la conexión
func (c *PocketBaseClient) Close() error {
	// No hay conexión persistente en el cliente HTTP
	return nil
}

// ==================== IN-MEMORY IMPLEMENTATIONS ====================

// InMemoryUserRepository implementación en memoria de UserRepository
type InMemoryUserRepository struct {
	users map[string]*entities.User
}

// NewInMemoryUserRepository crea un nuevo repositorio en memoria
func NewInMemoryUserRepository() *InMemoryUserRepository {
	// Crear usuarios de prueba
	users := make(map[string]*entities.User)

	// Usuario admin de prueba
	adminUser := &entities.User{
		ID:            "1",
		Email:         "admin@example.com",
		Username:      "admin",
		PasswordHash:  "$2a$10$IJKjCgI2DknPahDwkWzPkO/MC5E2Sqc6/H6wiBEbl5hVPuGtH7bye", // password123
		FirstName:     "Admin",
		LastName:      "User",
		RoleID:        "admin",
		IsActive:      true,
		IsSuperAdmin:  true,
		Phone:         nil,
		Department:    nil,
		LastLogin:     nil,
		LoginAttempts: 0,
		LockedUntil:   nil,
		MFAEnabled:    false,
		CreatedAt:     parseTime("2024-01-01T00:00:00Z"),
		UpdatedAt:     parseTime("2024-01-01T00:00:00Z"),
	}
	users[adminUser.ID] = adminUser

	return &InMemoryUserRepository{
		users: users,
	}
}

func (r *InMemoryUserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	if user, exists := r.users[id]; exists {
		return user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (r *InMemoryUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (r *InMemoryUserRepository) GetAll(ctx context.Context, page, pageSize int) ([]*entities.User, int, error) {
	users := make([]*entities.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, len(users), nil
}

func (r *InMemoryUserRepository) Create(ctx context.Context, user *entities.User) (*entities.User, error) {
	if _, exists := r.users[user.ID]; exists {
		return nil, fmt.Errorf("user already exists")
	}
	r.users[user.ID] = user
	return user, nil
}

func (r *InMemoryUserRepository) Update(ctx context.Context, id string, user *entities.User) (*entities.User, error) {
	if _, exists := r.users[id]; !exists {
		return nil, fmt.Errorf("user not found")
	}
	user.ID = id
	r.users[id] = user
	return user, nil
}

func (r *InMemoryUserRepository) Delete(ctx context.Context, id string) error {
	delete(r.users, id)
	return nil
}

// InMemoryRoleRepository implementación en memoria de RoleRepository
type InMemoryRoleRepository struct {
	roles map[string]*entities.Role
}

// NewInMemoryRoleRepository crea un nuevo repositorio de roles en memoria
func NewInMemoryRoleRepository() *InMemoryRoleRepository {
	roles := make(map[string]*entities.Role)

	// Rol admin
	adminRole := &entities.Role{
		ID:           "admin",
		Name:         "Administrator",
		Slug:         "admin",
		Description:  "Administrator role",
		Permissions:  []string{"users.view", "users.create", "users.edit", "users.delete"},
		Level:        1,
		IsSystemRole: true,
		CreatedAt:    parseTime("2024-01-01T00:00:00Z"),
		UpdatedAt:    parseTime("2024-01-01T00:00:00Z"),
	}
	roles[adminRole.ID] = adminRole

	return &InMemoryRoleRepository{
		roles: roles,
	}
}

func (r *InMemoryRoleRepository) GetByID(ctx context.Context, id string) (*entities.Role, error) {
	if role, exists := r.roles[id]; exists {
		return role, nil
	}
	return nil, fmt.Errorf("role not found")
}

func (r *InMemoryRoleRepository) GetAll(ctx context.Context, page, pageSize int) ([]*entities.Role, int, error) {
	roles := make([]*entities.Role, 0, len(r.roles))
	for _, role := range r.roles {
		roles = append(roles, role)
	}
	return roles, len(roles), nil
}

func (r *InMemoryRoleRepository) Create(ctx context.Context, role *entities.Role) (*entities.Role, error) {
	if _, exists := r.roles[role.ID]; exists {
		return nil, fmt.Errorf("role already exists")
	}
	r.roles[role.ID] = role
	return role, nil
}

func (r *InMemoryRoleRepository) Update(ctx context.Context, id string, role *entities.Role) (*entities.Role, error) {
	if _, exists := r.roles[id]; !exists {
		return nil, fmt.Errorf("role not found")
	}
	role.ID = id
	r.roles[id] = role
	return role, nil
}

func (r *InMemoryRoleRepository) Delete(ctx context.Context, id string) error {
	delete(r.roles, id)
	return nil
}

// Helper function para parsear fechas
func parseTime(timeStr string) time.Time {
	t, _ := time.Parse(time.RFC3339, timeStr)
	return t
}
