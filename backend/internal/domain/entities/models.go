package entities
package entities

import "time"

// User representa un usuario del sistema
type User struct {
	ID              string     `json:"id"`
	Email           string     `json:"email"`
	Username        string     `json:"username"`
	PasswordHash    string     `json:"-"` // Nunca serializar
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	RoleID          string     `json:"role_id"`
	IsActive        bool       `json:"is_active"`
	IsSuperAdmin    bool       `json:"is_super_admin"`
	Phone           *string    `json:"phone"`
	Department      *string    `json:"department"`
	LastLogin       *time.Time `json:"last_login"`
	LoginAttempts   int        `json:"login_attempts"`
	LockedUntil     *time.Time `json:"-"` // No serializar
	MFAEnabled      bool       `json:"mfa_enabled"`
	CreatedByID     *string    `json:"created_by_id"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// Role representa un rol del sistema
type Role struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Description  string    `json:"description"`
	Permissions  []string  `json:"permissions"`
	Level        int       `json:"level"`
	IsSystemRole bool      `json:"is_system_role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// WorkArea representa un área de trabajo
type WorkArea struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	Description  *string   `json:"description"`
	Location     string    `json:"location"`
	RiskLevel    string    `json:"risk_level"` // bajo, medio, alto
	SupervisorID *string   `json:"supervisor_id"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Employee representa un trabajador
type Employee struct {
	ID             string    `json:"id"`
	DocumentType   string    `json:"document_type"` // CC, CE, PASSPORT
	DocumentNumber string    `json:"document_number"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Email          *string   `json:"email"`
	Phone          *string   `json:"phone"`
	Position       string    `json:"position"`
	Department     string    `json:"department"`
	WorkAreaID     string    `json:"work_area_id"`
	UserID         *string   `json:"user_id"`
	HireDate       time.Time `json:"hire_date"`
	IsActive       bool      `json:"is_active"`
	Notes          *string   `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Observation representa una observación SBC
type Observation struct {
	ID                      string    `json:"id"`
	ObservationNumber       string    `json:"observation_number"`
	EmployeeID              string    `json:"employee_id"`
	ObserverID              string    `json:"observer_id"`
	WorkAreaID              string    `json:"work_area_id"`
	ObservationDate         time.Time `json:"observation_date"`
	TaskDescription         string    `json:"task_description"`
	DurationMinutes         int       `json:"duration_minutes"`
	HasSafeBehaviors        bool      `json:"has_safe_behaviors"`
	HasUnsafeBehaviors      bool      `json:"has_unsafe_behaviors"`
	OverallStatus           string    `json:"overall_status"` // safe, unsafe, mixed
	Comments                *string   `json:"comments"`
	Recommendations         *string   `json:"recommendations"`
	RequiresAction          bool      `json:"requires_action"`
	RequiresIncidentReport  bool      `json:"requires_incident_report"`
	CreatedByID             string    `json:"created_by_id"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

// Behavior representa una conducta en una observación
type Behavior struct {
	ID                 string    `json:"id"`
	ObservationID      string    `json:"observation_id"`
	BehaviorCategoryID string    `json:"behavior_category_id"`
	Type               string    `json:"type"` // safe, unsafe
	Sequence           int       `json:"sequence"`
	Description        string    `json:"description"`
	Timestamp          time.Time `json:"timestamp"`
	EvidencePhotoID    *string   `json:"evidence_photo_id"`
	CorrectionApplied  bool      `json:"correction_applied"`
	CorrectionDetails  *string   `json:"correction_details"`
	CreatedAt          time.Time `json:"created_at"`
}

// BehaviorCategory representa una categoría de conducta
type BehaviorCategory struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Code            string    `json:"code"`
	Type            string    `json:"type"` // safe, unsafe
	Description     string    `json:"description"`
	Severity        string    `json:"severity"` // low, medium, high, critical
	PreventionTips  *string   `json:"prevention_tips"`
	IsActive        bool      `json:"is_active"`
	Order           int       `json:"order"`
	CreatedAt       time.Time `json:"created_at"`
}

// Incident representa un incidente de seguridad
type Incident struct {
	ID                      string    `json:"id"`
	IncidentNumber          string    `json:"incident_number"`
	ObservationID           *string   `json:"observation_id"`
	EmployeeID              string    `json:"employee_id"`
	WorkAreaID              string    `json:"work_area_id"`
	IncidentDate            time.Time `json:"incident_date"`
	IncidentType            string    `json:"incident_type"` // accident, near_miss, hazard
	Severity                string    `json:"severity"`      // low, medium, high, critical
	Description             string    `json:"description"`
	Injuries                *string   `json:"injuries"`
	PropertyDamage          *string   `json:"property_damage"`
	RootCause               *string   `json:"root_cause"`
	Status                  string    `json:"status"` // open, in_progress, closed
	ReportedByID            string    `json:"reported_by_id"`
	AssignedToID            *string   `json:"assigned_to_id"`
	InvestigationNotes      *string   `json:"investigation_notes"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
	ClosedAt                *time.Time `json:"closed_at"`
}

// CorrectiveAction representa una acción correctiva
type CorrectiveAction struct {
	ID                   string    `json:"id"`
	ActionNumber         string    `json:"action_number"`
	ObservationID        *string   `json:"observation_id"`
	IncidentID           *string   `json:"incident_id"`
	ActionType           string    `json:"action_type"` // corrective, preventive, improvement
	Description          string    `json:"description"`
	ResponsibleID        string    `json:"responsible_id"`
	Priority             string    `json:"priority"` // low, medium, high, critical
	DueDate              time.Time `json:"due_date"`
	Status               string    `json:"status"` // open, in_progress, completed, overdue
	CompletionDate       *time.Time `json:"completion_date"`
	Evidence             *string   `json:"evidence"`
	EffectivenessCheck   *string   `json:"effectiveness_check"`
	CreatedByID          string    `json:"created_by_id"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// Attachment representa un archivo adjunto
type Attachment struct {
	ID                 string    `json:"id"`
	ObservationID      *string   `json:"observation_id"`
	IncidentID         *string   `json:"incident_id"`
	CorrectiveActionID *string   `json:"corrective_action_id"`
	FileName           string    `json:"file_name"`
	FileType           string    `json:"file_type"` // image, video, document, other
	MimeType           string    `json:"mime_type"`
	FileSize           int64     `json:"file_size"`
	FilePath           string    `json:"file_path"`
	Description        *string   `json:"description"`
	UploadedByID       string    `json:"uploaded_by_id"`
	CreatedAt          time.Time `json:"created_at"`
}

// Notification representa una notificación
type Notification struct {
	ID                string     `json:"id"`
	RecipientID       string     `json:"recipient_id"`
	NotificationType  string     `json:"notification_type"` // info, warning, critical
	Title             string     `json:"title"`
	Message           string     `json:"message"`
	RelatedEntityType *string    `json:"related_entity_type"`
	RelatedEntityID   *string    `json:"related_entity_id"`
	IsRead            bool       `json:"is_read"`
	ReadAt            *time.Time `json:"read_at"`
	ActionURL         *string    `json:"action_url"`
	CreatedAt         time.Time  `json:"created_at"`
}

// AuditLog representa un registro de auditoría
type AuditLog struct {
	ID           string                 `json:"id"`
	UserID       string                 `json:"user_id"`
	Action       string                 `json:"action"`
	EntityType   string                 `json:"entity_type"`
	EntityID     string                 `json:"entity_id"`
	OldValues    map[string]interface{} `json:"old_values"`
	NewValues    map[string]interface{} `json:"new_values"`
	ChangeSummary *string               `json:"change_summary"`
	IPAddress    string                 `json:"ip_address"`
	UserAgent    string                 `json:"user_agent"`
	Status       string                 `json:"status"` // success, failure
	ErrorMessage *string                `json:"error_message"`
	CreatedAt    time.Time              `json:"created_at"`
}

// KPI representa un indicador clave de desempeño
type KPI struct {
	ID                       string     `json:"id"`
	KPIName                  string     `json:"kpi_name"`
	KPICode                  string     `json:"kpi_code"`
	Period                   string     `json:"period"` // daily, weekly, monthly, yearly
	PeriodStart              time.Time  `json:"period_start"`
	PeriodEnd                time.Time  `json:"period_end"`
	WorkAreaID               *string    `json:"work_area_id"`
	TotalObservations        int        `json:"total_observations"`
	SafeBehaviorsCount       int        `json:"safe_behaviors_count"`
	UnsafeBehaviorsCount     int        `json:"unsafe_behaviors_count"`
	SafePercentage           float64    `json:"safe_percentage"`
	UnsafePercentage         float64    `json:"unsafe_percentage"`
	IncidentsCount           int        `json:"incidents_count"`
	HighSeverityIncidents    int        `json:"high_severity_incidents"`
	CorrectiveActionsOpen    int        `json:"corrective_actions_open"`
	CorrectiveActionsOverdue int        `json:"corrective_actions_overdue"`
	CorrectiveActionsCompleted int      `json:"corrective_actions_completed"`
	ObservationFrequency     float64    `json:"observation_frequency"`
	EmployeesObserved        int        `json:"employees_observed"`
	TrendVsPrevious          *string    `json:"trend_vs_previous"`
	CalculatedAt             time.Time  `json:"calculated_at"`
}

// Report representa un reporte generado
type Report struct {
	ID            string    `json:"id"`
	ReportName    string    `json:"report_name"`
	ReportType    string    `json:"report_type"` // observations, incidents, kpis, monthly, custom
	PeriodStart   time.Time `json:"period_start"`
	PeriodEnd     time.Time `json:"period_end"`
	WorkAreaID    *string   `json:"work_area_id"`
	Format        string    `json:"format"` // pdf, excel, json
	FilePath      string    `json:"file_path"`
	FileSize      int64     `json:"file_size"`
	GeneratedByID string    `json:"generated_by_id"`
	Status        string    `json:"status"` // pending, generating, completed, failed
	ErrorMessage  *string   `json:"error_message"`
	CreatedAt     time.Time `json:"created_at"`
	GeneratedAt   *time.Time `json:"generated_at"`
}
