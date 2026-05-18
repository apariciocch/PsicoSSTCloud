package service
package service

import (
	"context"

	"github.com/apariciocch/psicosstcloud/internal/domain/entities"
	"github.com/apariciocch/psicosstcloud/internal/dto"
)

// AuthService interfaz para operaciones de autenticación
type AuthService interface {
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	Logout(ctx context.Context, userID string) error
	RefreshToken(ctx context.Context, refreshToken string) (*dto.LoginResponse, error)
	ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error
	ValidateToken(tokenString string) (map[string]interface{}, error)
}

// UserService interfaz para operaciones con usuarios
type UserService interface {
	GetByID(ctx context.Context, id string) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetAll(ctx context.Context, page, pageSize int) ([]*entities.User, int, error)
	Create(ctx context.Context, req *dto.CreateUserRequest) (*entities.User, error)
	Update(ctx context.Context, id string, req *dto.UpdateUserRequest) (*entities.User, error)
	Delete(ctx context.Context, id string) error
}

// EmployeeService interfaz para operaciones con empleados
type EmployeeService interface {
	GetByID(ctx context.Context, id string) (*entities.Employee, error)
	GetAll(ctx context.Context, page, pageSize int, filter *dto.ObservationFilter) ([]*entities.Employee, int, error)
	Create(ctx context.Context, req *dto.CreateEmployeeRequest) (*entities.Employee, error)
	Update(ctx context.Context, id string, req *dto.UpdateEmployeeRequest) (*entities.Employee, error)
	Delete(ctx context.Context, id string) error
}

// ObservationService interfaz para operaciones con observaciones
type ObservationService interface {
	GetByID(ctx context.Context, id string) (*entities.Observation, error)
	GetAll(ctx context.Context, page, pageSize int, filter *dto.ObservationFilter) ([]*entities.Observation, int, error)
	Create(ctx context.Context, req *dto.CreateObservationRequest, userID string) (*entities.Observation, error)
	Update(ctx context.Context, id string, req *dto.UpdateObservationRequest) (*entities.Observation, error)
	Delete(ctx context.Context, id string) error
}

// UpdateObservationRequest para actualizar observaciones
type UpdateObservationRequest struct {
	Comments               *string
	Recommendations        *string
	RequiresAction         *bool
	RequiresIncidentReport *bool
	OverallStatus          *string
}

// AnalyticsService interfaz para operaciones de analytics
type AnalyticsService interface {
	GetDashboard(ctx context.Context, workAreaID *string, period string) (*dto.DashboardResponse, error)
	GetObservationStats(ctx context.Context, workAreaID *string) (map[string]interface{}, error)
	GetIncidentStats(ctx context.Context, workAreaID *string) (map[string]interface{}, error)
	CalculateKPIs(ctx context.Context, workAreaID *string, period string) ([]*entities.KPI, error)
}

// IncidentService interfaz para operaciones con incidentes
type IncidentService interface {
	GetByID(ctx context.Context, id string) (*entities.Incident, error)
	GetAll(ctx context.Context, page, pageSize int, filter *dto.IncidentFilter) ([]*entities.Incident, int, error)
	Create(ctx context.Context, req *dto.CreateIncidentRequest, userID string) (*entities.Incident, error)
	Update(ctx context.Context, id string, data map[string]interface{}) (*entities.Incident, error)
	Close(ctx context.Context, id, reason string) (*entities.Incident, error)
	Delete(ctx context.Context, id string) error
}

// ReportService interfaz para operaciones con reportes
type ReportService interface {
	GenerateObservationReport(ctx context.Context, workAreaID *string, startDate, endDate string, format string) (*entities.Report, error)
	GenerateIncidentReport(ctx context.Context, workAreaID *string, startDate, endDate string, format string) (*entities.Report, error)
	GenerateKPIReport(ctx context.Context, workAreaID *string, period string) (*entities.Report, error)
	GetByID(ctx context.Context, id string) (*entities.Report, error)
	GetAll(ctx context.Context, page, pageSize int) ([]*entities.Report, int, error)
}

// NotificationService interfaz para operaciones con notificaciones
type NotificationService interface {
	GetByUserID(ctx context.Context, userID string, page, pageSize int) ([]*entities.Notification, int, error)
	GetUnread(ctx context.Context, userID string) ([]*entities.Notification, error)
	MarkAsRead(ctx context.Context, notificationID string) error
	MarkAllAsRead(ctx context.Context, userID string) error
	Send(ctx context.Context, notification *entities.Notification) error
}

// AuditService interfaz para operaciones de auditoría
type AuditService interface {
	Log(ctx context.Context, audit *entities.AuditLog) error
	GetByUserID(ctx context.Context, userID string, page, pageSize int) ([]*entities.AuditLog, int, error)
	GetByEntity(ctx context.Context, entityType, entityID string) ([]*entities.AuditLog, error)
}
