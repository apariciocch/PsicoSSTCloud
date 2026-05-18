package config

// Roles del sistema
const (
	RoleAdmin      = "role_admin"
	RoleSupervisor = "role_supervisor"
	RoleObserver   = "role_observer"
	RoleWorker     = "role_worker"
	RoleAuditor    = "role_auditor"
)

// Permisos del sistema
const (
	// Users
	PermUsersView   = "users.view"
	PermUsersCreate = "users.create"
	PermUsersEdit   = "users.edit"
	PermUsersDelete = "users.delete"

	// Observations
	PermObservationsView   = "observations.view"
	PermObservationsCreate = "observations.create"
	PermObservationsEdit   = "observations.edit"
	PermObservationsDelete = "observations.delete"

	// Employees
	PermEmployeesView   = "employees.view"
	PermEmployeesCreate = "employees.create"
	PermEmployeesEdit   = "employees.edit"

	// Incidents
	PermIncidentsView   = "incidents.view"
	PermIncidentsCreate = "incidents.create"
	PermIncidentsEdit   = "incidents.edit"

	// Corrective Actions
	PermCorrectiveActionsView   = "corrective_actions.view"
	PermCorrectiveActionsCreate = "corrective_actions.create"
	PermCorrectiveActionsEdit   = "corrective_actions.edit"

	// Reports
	PermReportsGenerate = "reports.generate"
	PermReportsExport   = "reports.export"

	// Analytics
	PermAnalyticsView = "analytics.view"

	// Audit
	PermAuditLogsView = "audit_logs.view"

	// Settings
	PermSettingsEdit = "settings.edit"
)

// Estados de observación
const (
	ObservationStatusSafe   = "safe"
	ObservationStatusUnsafe = "unsafe"
	ObservationStatusMixed  = "mixed"
)

// Tipos de conducta
const (
	BehaviorTypeSafe   = "safe"
	BehaviorTypeUnsafe = "unsafe"
)

// Niveles de riesgo
const (
	RiskLevelLow    = "bajo"
	RiskLevelMedium = "medio"
	RiskLevelHigh   = "alto"
)

// Severidad
const (
	SeverityLow      = "low"
	SeverityMedium   = "medium"
	SeverityHigh     = "high"
	SeverityCritical = "critical"
)

// Tipos de incidente
const (
	IncidentTypeAccident = "accident"
	IncidentTypeNearMiss = "near_miss"
	IncidentTypeHazard   = "hazard"
)

// Estados de incidente
const (
	IncidentStatusOpen       = "open"
	IncidentStatusInProgress = "in_progress"
	IncidentStatusClosed     = "closed"
)

// Estados de acciones correctivas
const (
	ActionStatusOpen       = "open"
	ActionStatusInProgress = "in_progress"
	ActionStatusCompleted  = "completed"
	ActionStatusOverdue    = "overdue"
)

// Tipos de acciones correctivas
const (
	ActionTypeCorrectiveAction = "corrective"
	ActionTypePreventive       = "preventive"
	ActionTypeImprovement      = "improvement"
)

// Tipos de archivos
const (
	FileTypeImage    = "image"
	FileTypeVideo    = "video"
	FileTypeDocument = "document"
	FileTypeOther    = "other"
)

// Tipos de notificación
const (
	NotificationTypeInfo     = "info"
	NotificationTypeWarning  = "warning"
	NotificationTypeCritical = "critical"
)

// Estados de reporte
const (
	ReportStatusPending    = "pending"
	ReportStatusGenerating = "generating"
	ReportStatusCompleted  = "completed"
	ReportStatusFailed     = "failed"
)

// Formatos de reporte
const (
	ReportFormatPDF   = "pdf"
	ReportFormatExcel = "excel"
	ReportFormatJSON  = "json"
)

// Tipos de reporte
const (
	ReportTypeObservations = "observations"
	ReportTypeIncidents    = "incidents"
	ReportTypeKPIs         = "kpis"
	ReportTypeMonthly      = "monthly"
	ReportTypeCustom       = "custom"
)

// Acciones de auditoría
const (
	AuditActionCreate  = "create"
	AuditActionUpdate  = "update"
	AuditActionDelete  = "delete"
	AuditActionLogin   = "login"
	AuditActionLogout  = "logout"
	AuditActionExport  = "export"
)

// Estados de auditoría
const (
	AuditStatusSuccess = "success"
	AuditStatusFailure = "failure"
)

// Períodos de KPI
const (
	KPIPeriodDaily   = "daily"
	KPIPeriodWeekly  = "weekly"
	KPIPeriodMonthly = "monthly"
	KPIPeriodYearly  = "yearly"
)

// Límites del sistema
const (
	MaxLoginAttempts      = 5
	TokenExpiration       = 900           // 15 minutos
	TokenRefreshExpiration = 604800       // 7 días
	PasswordMinLength     = 8
	PasswordMaxLength     = 128
	UsernameMinLength     = 3
	UsernameMaxLength     = 50
	DefaultPageSize       = 20
	MaxPageSize           = 100
)

// Documentos
const (
	DocumentTypeCC       = "CC"
	DocumentTypeCE       = "CE"
	DocumentTypePassport = "PASSPORT"
)

// Colecciones PocketBase
const (
	CollectionUsers              = "users"
	CollectionRoles             = "roles"
	CollectionWorkAreas         = "work_areas"
	CollectionEmployees         = "employees"
	CollectionObservations      = "observations"
	CollectionBehaviors         = "behaviors"
	CollectionBehaviorCategories = "behavior_categories"
	CollectionIncidents         = "incidents"
	CollectionCorrectiveActions = "corrective_actions"
	CollectionAttachments       = "attachments"
	CollectionNotifications     = "notifications"
	CollectionAuditLogs         = "audit_logs"
	CollectionKPIs              = "kpis"
	CollectionReports           = "reports"
)

// Mensajes de error comunes
const (
	ErrorUnauthorized       = "No autorizado"
	ErrorForbidden          = "Acceso denegado"
	ErrorNotFound           = "No encontrado"
	ErrorAlreadyExists      = "Ya existe"
	ErrorInvalidRequest     = "Solicitud inválida"
	ErrorInvalidCredentials = "Credenciales inválidas"
	ErrorUserLocked         = "Usuario bloqueado temporalmente"
	ErrorWeakPassword       = "Contraseña débil"
)
