package dto

import "time"

// ============ AUTH DTOs ============

// LoginRequest solicitud de login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=128"`
}

// LoginResponse respuesta de login
type LoginResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    int           `json:"expires_in"`
	User         *UserResponse `json:"user"`
}

// RefreshTokenRequest solicitud de refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ============ USER DTOs ============

// CreateUserRequest solicitud de creación de usuario
type CreateUserRequest struct {
	Email      string  `json:"email" validate:"required,email,max=255"`
	Username   string  `json:"username" validate:"required,min=3,max=50,alphanum"`
	Password   string  `json:"password" validate:"required,min=8,max=128"`
	FirstName  string  `json:"first_name" validate:"required,max=100"`
	LastName   string  `json:"last_name" validate:"required,max=100"`
	RoleID     string  `json:"role_id" validate:"required"`
	Phone      *string `json:"phone" validate:"omitempty,max=20"`
	Department *string `json:"department" validate:"omitempty,max=100"`
}

// UpdateUserRequest solicitud de actualización de usuario
type UpdateUserRequest struct {
	FirstName  *string `json:"first_name" validate:"omitempty,max=100"`
	LastName   *string `json:"last_name" validate:"omitempty,max=100"`
	Phone      *string `json:"phone" validate:"omitempty,max=20"`
	Department *string `json:"department" validate:"omitempty,max=100"`
}

// ChangePasswordRequest cambio de contraseña
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=128"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// UserResponse respuesta con datos de usuario
type UserResponse struct {
	ID           string     `json:"id"`
	Email        string     `json:"email"`
	Username     string     `json:"username"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	RoleID       string     `json:"role_id"`
	IsActive     bool       `json:"is_active"`
	IsSuperAdmin bool       `json:"is_super_admin"`
	Phone        *string    `json:"phone"`
	Department   *string    `json:"department"`
	LastLogin    *time.Time `json:"last_login"`
	MFAEnabled   bool       `json:"mfa_enabled"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// ============ EMPLOYEE DTOs ============

// CreateEmployeeRequest solicitud de creación de empleado
type CreateEmployeeRequest struct {
	DocumentType   string    `json:"document_type" validate:"required,oneof=CC CE PASSPORT"`
	DocumentNumber string    `json:"document_number" validate:"required,max=20"`
	FirstName      string    `json:"first_name" validate:"required,max=100"`
	LastName       string    `json:"last_name" validate:"required,max=100"`
	Email          *string   `json:"email" validate:"omitempty,email"`
	Phone          *string   `json:"phone" validate:"omitempty,max=20"`
	Position       string    `json:"position" validate:"required,max=100"`
	Department     string    `json:"department" validate:"required,max=100"`
	WorkAreaID     string    `json:"work_area_id" validate:"required"`
	HireDate       time.Time `json:"hire_date" validate:"required"`
	Notes          *string   `json:"notes" validate:"omitempty,max=500"`
}

// UpdateEmployeeRequest solicitud de actualización de empleado
type UpdateEmployeeRequest struct {
	Email      *string `json:"email" validate:"omitempty,email"`
	Phone      *string `json:"phone" validate:"omitempty,max=20"`
	Position   *string `json:"position" validate:"omitempty,max=100"`
	Department *string `json:"department" validate:"omitempty,max=100"`
	WorkAreaID *string `json:"work_area_id" validate:"omitempty"`
	Notes      *string `json:"notes" validate:"omitempty,max=500"`
	IsActive   *bool   `json:"is_active"`
}

// EmployeeResponse respuesta de empleado
type EmployeeResponse struct {
	ID             string    `json:"id"`
	DocumentType   string    `json:"document_type"`
	DocumentNumber string    `json:"document_number"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Email          *string   `json:"email"`
	Phone          *string   `json:"phone"`
	Position       string    `json:"position"`
	Department     string    `json:"department"`
	WorkAreaID     string    `json:"work_area_id"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ============ OBSERVATION DTOs ============

// CreateObservationRequest solicitud de creación de observación
type CreateObservationRequest struct {
	EmployeeID             string                  `json:"employee_id" validate:"required"`
	WorkAreaID             string                  `json:"work_area_id" validate:"required"`
	TaskDescription        string                  `json:"task_description" validate:"required,max=500"`
	DurationMinutes        int                     `json:"duration_minutes" validate:"required,min=5,max=480"`
	HasSafeBehaviors       bool                    `json:"has_safe_behaviors"`
	HasUnsafeBehaviors     bool                    `json:"has_unsafe_behaviors"`
	Comments               *string                 `json:"comments" validate:"omitempty,max=2000"`
	Recommendations        *string                 `json:"recommendations" validate:"omitempty,max=2000"`
	RequiresAction         bool                    `json:"requires_action"`
	RequiresIncidentReport bool                    `json:"requires_incident_report"`
	Behaviors              []CreateBehaviorRequest `json:"behaviors" validate:"required,min=1"`
}

// CreateBehaviorRequest solicitud de creación de conducta
type CreateBehaviorRequest struct {
	BehaviorCategoryID string    `json:"behavior_category_id" validate:"required"`
	Type               string    `json:"type" validate:"required,oneof=safe unsafe"`
	Sequence           int       `json:"sequence" validate:"required,min=1"`
	Description        string    `json:"description" validate:"required,max=500"`
	Timestamp          time.Time `json:"timestamp" validate:"required"`
	CorrectionApplied  bool      `json:"correction_applied"`
	CorrectionDetails  *string   `json:"correction_details" validate:"omitempty,max=500"`
}

// ObservationResponse respuesta de observación
type ObservationResponse struct {
	ID                     string    `json:"id"`
	ObservationNumber      string    `json:"observation_number"`
	EmployeeID             string    `json:"employee_id"`
	ObserverID             string    `json:"observer_id"`
	WorkAreaID             string    `json:"work_area_id"`
	ObservationDate        time.Time `json:"observation_date"`
	TaskDescription        string    `json:"task_description"`
	DurationMinutes        int       `json:"duration_minutes"`
	HasSafeBehaviors       bool      `json:"has_safe_behaviors"`
	HasUnsafeBehaviors     bool      `json:"has_unsafe_behaviors"`
	OverallStatus          string    `json:"overall_status"`
	Comments               *string   `json:"comments"`
	Recommendations        *string   `json:"recommendations"`
	RequiresAction         bool      `json:"requires_action"`
	RequiresIncidentReport bool      `json:"requires_incident_report"`
	SafeCount              int       `json:"safe_count"`
	UnsafeCount            int       `json:"unsafe_count"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

// ============ UPDATE OBSERVATION DTOs ============

// UpdateObservationRequest solicitud de actualización de observación
type UpdateObservationRequest struct {
	Comments               *string `json:"comments" validate:"omitempty,max=2000"`
	Recommendations        *string `json:"recommendations" validate:"omitempty,max=2000"`
	RequiresAction         *bool   `json:"requires_action"`
	RequiresIncidentReport *bool   `json:"requires_incident_report"`
	OverallStatus          *string `json:"overall_status" validate:"omitempty,oneof=completed pending archived"`
}

// ============ INCIDENT DTOs ============

// CreateIncidentRequest solicitud de creación de incidente
type CreateIncidentRequest struct {
	ObservationID          *string   `json:"observation_id"`
	EmployeeID             string    `json:"employee_id" validate:"required"`
	WorkAreaID             string    `json:"work_area_id" validate:"required"`
	IncidentDate           time.Time `json:"incident_date" validate:"required"`
	IncidentType           string    `json:"incident_type" validate:"required,oneof=accident near_miss hazard"`
	Severity               string    `json:"severity" validate:"required,oneof=low medium high critical"`
	Description            string    `json:"description" validate:"required,max=2000"`
	Injuries               *string   `json:"injuries" validate:"omitempty,max=500"`
	PropertyDamage         *string   `json:"property_damage" validate:"omitempty,max=500"`
	RequiresIncidentReport bool      `json:"requires_incident_report"`
}

// ============ DASHBOARD DTOs ============

// DashboardResponse respuesta del dashboard
type DashboardResponse struct {
	Period                 string          `json:"period"`
	PeriodStart            time.Time       `json:"period_start"`
	PeriodEnd              time.Time       `json:"period_end"`
	TotalObservations      int             `json:"total_observations"`
	SafeBehaviorsPercent   float64         `json:"safe_behaviors_percent"`
	UnsafeBehaviorsPercent float64         `json:"unsafe_behaviors_percent"`
	TotalIncidents         int             `json:"total_incidents"`
	CriticalIncidents      int             `json:"critical_incidents"`
	IncidentsByArea        []AreaIncidents `json:"incidents_by_area"`
	CriticalAreas          []AreaRisk      `json:"critical_areas"`
	OverdueActions         int             `json:"overdue_actions"`
	CompletedActions       int             `json:"completed_actions"`
	OpenActions            int             `json:"open_actions"`
	Trend                  string          `json:"trend"` // up, down, stable
	MostObservedAreas      []AreaStats     `json:"most_observed_areas"`
	TopBehaviors           []BehaviorStat  `json:"top_behaviors"`
}

// AreaIncidents incidentes por área
type AreaIncidents struct {
	AreaID   string `json:"area_id"`
	AreaName string `json:"area_name"`
	Count    int    `json:"count"`
	Severity string `json:"severity"`
}

// AreaRisk riesgo por área
type AreaRisk struct {
	AreaID        string  `json:"area_id"`
	AreaName      string  `json:"area_name"`
	UnsafePercent float64 `json:"unsafe_percent"`
	IncidentCount int     `json:"incident_count"`
	RiskLevel     string  `json:"risk_level"`
}

// AreaStats estadísticas de área
type AreaStats struct {
	AreaID           string  `json:"area_id"`
	AreaName         string  `json:"area_name"`
	ObservationCount int     `json:"observation_count"`
	SafePercent      float64 `json:"safe_percent"`
	IncidentCount    int     `json:"incident_count"`
}

// BehaviorStat estadística de conducta
type BehaviorStat struct {
	BehaviorID      string `json:"behavior_id"`
	BehaviorName    string `json:"behavior_name"`
	Type            string `json:"type"`
	OccurrenceCount int    `json:"occurrence_count"`
	Severity        string `json:"severity"`
}

// ============ PAGINATION DTOs ============

// PaginatedResponse respuesta paginada genérica
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	Total      int         `json:"total"`
	TotalPages int         `json:"total_pages"`
	HasNext    bool        `json:"has_next"`
	HasPrev    bool        `json:"has_prev"`
}

// ============ ERROR RESPONSE ============

// ErrorResponse respuesta de error estándar
type ErrorResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Field   string      `json:"field,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

// ============ FILTER DTOs ============

// ObservationFilter filtros para observaciones
type ObservationFilter struct {
	EmployeeID    *string
	WorkAreaID    *string
	ObserverID    *string
	OverallStatus *string
	StartDate     *time.Time
	EndDate       *time.Time
	Page          int
	PageSize      int
}

// IncidentFilter filtros para incidentes
type IncidentFilter struct {
	EmployeeID *string
	WorkAreaID *string
	Severity   *string
	Status     *string
	StartDate  *time.Time
	EndDate    *time.Time
	Page       int
	PageSize   int
}
