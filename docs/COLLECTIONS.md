# Diseño Completo de Colecciones PocketBase

## 1. COLECCIONES PRINCIPALES

---

## 1.1 Colección: `users`

**Descripción**: Tabla central de usuarios del sistema. Administra acceso, autenticación y perfiles.

### Campos:

| Campo | Tipo | Descripción | Required | Validación |
|-------|------|-------------|----------|-----------|
| `id` | String (UID) | ID único generado | ✓ | Auto PocketBase |
| `email` | Email | Email único del usuario | ✓ | unique, valid email |
| `username` | String | Nombre de usuario | ✓ | unique, 3-50 chars |
| `password_hash` | String | Hash bcrypt contraseña | ✓ | minLength: 60 |
| `first_name` | String | Primer nombre | ✓ | 2-50 chars |
| `last_name` | String | Apellido | ✓ | 2-50 chars |
| `role_id` | Relation | Referencia a `roles` | ✓ | many-to-one |
| `is_active` | Bool | Usuario activo | ✓ | default: true |
| `is_super_admin` | Bool | Super administrador | ✓ | default: false |
| `phone` | String | Teléfono contacto | ✗ | max: 20 chars |
| `department` | String | Departamento | ✗ | max: 100 chars |
| `last_login` | DateTime | Último login | ✗ | Auto timestamp |
| `login_attempts` | Int | Intentos fallidos | ✓ | default: 0 |
| `locked_until` | DateTime | Bloqueado hasta | ✗ | Para rate limiting |
| `mfa_enabled` | Bool | Autenticación de 2 factores | ✓ | default: false |
| `created_by_id` | Relation | Quién creó el usuario | ✗ | ref to users |
| `created_at` | DateTime | Fecha creación | ✓ | Auto now |
| `updated_at` | DateTime | Fecha actualización | ✓ | Auto now |

### Índices:
```
CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE UNIQUE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_role_id ON users(role_id);
CREATE INDEX idx_users_is_active ON users(is_active);
CREATE INDEX idx_users_created_at ON users(created_at);
```

### Reglas de Acceso:

**CREATE**:
```
user.is_super_admin = true
OR user.role_id in (admin_role_id)
```

**READ**: 
```
user.id = @requestUser.id
OR user.is_super_admin = true
OR user.role_id in (admin_role_id, auditor_role_id)
```

**UPDATE**:
```
user.id = @requestUser.id  (cambios propios)
OR @requestUser.is_super_admin = true
OR (@requestUser.role_id in (admin_role_id) AND user.role_id != admin_role_id)
```

**DELETE**:
```
@requestUser.is_super_admin = true
AND user.id != @requestUser.id
```

### Ejemplo JSON:
```json
{
  "id": "user_001",
  "email": "juan.perez@empresa.com",
  "username": "jperez",
  "password_hash": "$2a$12$...",
  "first_name": "Juan",
  "last_name": "Pérez",
  "role_id": "role_observer",
  "is_active": true,
  "is_super_admin": false,
  "phone": "+57-300-555-0123",
  "department": "Operaciones",
  "last_login": "2026-05-18T14:30:00.000Z",
  "login_attempts": 0,
  "mfa_enabled": false,
  "created_by_id": "user_admin",
  "created_at": "2026-01-15T08:00:00.000Z",
  "updated_at": "2026-05-18T14:30:00.000Z"
}
```

---

## 1.2 Colección: `roles`

**Descripción**: Define los roles del sistema con permisos asociados.

### Campos:

| Campo | Tipo | Descripción | Required |
|-------|------|-------------|----------|
| `id` | String (UID) | ID único | ✓ |
| `name` | String | Nombre del rol | ✓ |
| `slug` | String | Slug único | ✓ |
| `description` | String | Descripción | ✓ |
| `permissions` | JSON Array | Lista de permisos | ✓ |
| `level` | Int | Nivel jerárquico | ✓ |
| `is_system_role` | Bool | Rol del sistema | ✓ |
| `created_at` | DateTime | Fecha creación | ✓ |
| `updated_at` | DateTime | Fecha actualización | ✓ |

### Permisos Posibles:
```
users.view, users.create, users.edit, users.delete
roles.view, roles.edit
employees.view, employees.create, employees.edit
observations.view, observations.create, observations.edit, observations.delete
behaviors.view, behaviors.create
incidents.view, incidents.create, incidents.edit
corrective_actions.view, corrective_actions.create, corrective_actions.edit
reports.generate, reports.export
analytics.view
audit_logs.view
settings.edit
```

### Reglas de Acceso:

**CREATE/UPDATE/DELETE**: Solo super_admin
**READ**: Todos los usuarios autenticados

### Datos Iniciales (Insert):
```json
[
  {
    "id": "role_admin",
    "name": "Administrador",
    "slug": "admin",
    "description": "Acceso total al sistema",
    "permissions": ["users.view", "users.create", "users.edit", "users.delete", "roles.view", "roles.edit", "employees.view", "employees.create", "employees.edit", "observations.view", "observations.create", "observations.edit", "observations.delete", "behaviors.view", "behaviors.create", "incidents.view", "incidents.create", "incidents.edit", "corrective_actions.view", "corrective_actions.create", "corrective_actions.edit", "reports.generate", "reports.export", "analytics.view", "audit_logs.view", "settings.edit"],
    "level": 1,
    "is_system_role": true,
    "created_at": "2026-01-01T00:00:00.000Z",
    "updated_at": "2026-01-01T00:00:00.000Z"
  },
  {
    "id": "role_supervisor",
    "name": "Supervisor",
    "slug": "supervisor",
    "description": "Supervisa observadores y revisa observaciones",
    "permissions": ["users.view", "employees.view", "employees.create", "employees.edit", "observations.view", "observations.create", "observations.edit", "behaviors.view", "behaviors.create", "incidents.view", "incidents.create", "incidents.edit", "corrective_actions.view", "corrective_actions.create", "corrective_actions.edit", "reports.generate", "analytics.view", "audit_logs.view"],
    "level": 2,
    "is_system_role": true,
    "created_at": "2026-01-01T00:00:00.000Z",
    "updated_at": "2026-01-01T00:00:00.000Z"
  },
  {
    "id": "role_observer",
    "name": "Observador SBC",
    "slug": "observer",
    "description": "Realiza observaciones y registra conductas",
    "permissions": ["employees.view", "observations.view", "observations.create", "observations.edit", "behaviors.view", "behaviors.create", "reports.generate"],
    "level": 3,
    "is_system_role": true,
    "created_at": "2026-01-01T00:00:00.000Z",
    "updated_at": "2026-01-01T00:00:00.000Z"
  },
  {
    "id": "role_worker",
    "name": "Trabajador",
    "slug": "worker",
    "description": "Acceso limitado a sus propias observaciones",
    "permissions": ["observations.view"],
    "level": 4,
    "is_system_role": true,
    "created_at": "2026-01-01T00:00:00.000Z",
    "updated_at": "2026-01-01T00:00:00.000Z"
  },
  {
    "id": "role_auditor",
    "name": "Auditor",
    "slug": "auditor",
    "description": "Acceso lectura solo para auditoría",
    "permissions": ["users.view", "employees.view", "observations.view", "incidents.view", "corrective_actions.view", "reports.generate", "analytics.view", "audit_logs.view"],
    "level": 2,
    "is_system_role": true,
    "created_at": "2026-01-01T00:00:00.000Z",
    "updated_at": "2026-01-01T00:00:00.000Z"
  }
]
```

---

## 1.3 Colección: `work_areas`

**Descripción**: Áreas físicas de trabajo donde se realizan observaciones.

### Campos:

| Campo | Tipo | Descripción | Required |
|-------|------|-------------|----------|
| `id` | String (UID) | ID único | ✓ |
| `name` | String | Nombre del área | ✓ |
| `code` | String | Código único | ✓ |
| `description` | String | Descripción | ✗ |
| `location` | String | Ubicación/Piso | ✓ |
| `risk_level` | Select | alto\|medio\|bajo | ✓ |
| `supervisor_id` | Relation | Supervisor responsable | ✗ |
| `is_active` | Bool | Área activa | ✓ |
| `created_at` | DateTime | Fecha creación | ✓ |
| `updated_at` | DateTime | Fecha actualización | ✓ |

### Índices:
```
CREATE UNIQUE INDEX idx_work_areas_code ON work_areas(code);
CREATE INDEX idx_work_areas_risk_level ON work_areas(risk_level);
```

---

## 1.4 Colección: `employees`

**Descripción**: Registro de trabajadores observados en el SBC.

### Campos:

| Campo | Tipo | Descripción | Required |
|-------|------|-------------|----------|
| `id` | String (UID) | ID único | ✓ |
| `document_type` | Select | CC\|CE\|PASSPORT | ✓ |
| `document_number` | String | Número documento | ✓ |
| `first_name` | String | Primer nombre | ✓ |
| `last_name` | String | Apellido | ✓ |
| `email` | Email | Email personal | ✗ |
| `phone` | String | Teléfono | ✗ |
| `position` | String | Cargo | ✓ |
| `department` | String | Departamento | ✓ |
| `work_area_id` | Relation | Área de trabajo | ✓ |
| `user_id` | Relation | Usuario asociado (opcional) | ✗ |
| `hire_date` | Date | Fecha contratación | ✓ |
| `is_active` | Bool | Empleado activo | ✓ |
| `notes` | String | Notas adicionales | ✗ |
| `created_at` | DateTime | Fecha creación | ✓ |
| `updated_at` | DateTime | Fecha actualización | ✓ |

### Índices:
```
CREATE UNIQUE INDEX idx_employees_document ON employees(document_type, document_number);
CREATE INDEX idx_employees_work_area_id ON employees(work_area_id);
CREATE INDEX idx_employees_is_active ON employees(is_active);
```

---

## 1.5 Colección: `behavior_categories`

**Descripción**: Categorías de conductas (seguras/inseguras) predefinidas.

### Campos:

| Campo | Tipo | Descripción | Required |
|-------|------|-------------|----------|
| `id` | String (UID) | ID único | ✓ |
| `name` | String | Nombre categoría | ✓ |
| `code` | String | Código único | ✓ |
| `type` | Select | safe\|unsafe | ✓ |
| `description` | String | Descripción detallada | ✓ |
| `severity` | Select | low\|medium\|high\|critical | ✓ |
| `prevention_tips` | String | Consejos prevención | ✗ |
| `is_active` | Bool | Categoría activa | ✓ |
| `order` | Int | Orden presentación | ✓ |
| `created_at` | DateTime | Fecha creación | ✓ |

### Datos Iniciales:
```json
[
  {
    "id": "cat_001",
    "name": "Uso correcto de EPP",
    "code": "SAFE_EPP",
    "type": "safe",
    "description": "Trabajador utiliza equipo de protección personal completo",
    "severity": "high",
    "prevention_tips": "Verificar que el EPP sea el correcto para la tarea",
    "is_active": true,
    "order": 1
  },
  {
    "id": "cat_002",
    "name": "No usar EPP",
    "code": "UNSAFE_NO_EPP",
    "type": "unsafe",
    "description": "Trabajador realiza tarea sin equipo de protección requerido",
    "severity": "critical",
    "prevention_tips": "Reforzar obligatoriedad del uso de EPP",
    "is_active": true,
    "order": 2
  },
  {
    "id": "cat_003",
    "name": "Postura ergonómica",
    "code": "SAFE_POSTURE",
    "type": "safe",
    "description": "Trabajador mantiene postura correcta",
    "severity": "medium",
    "prevention_tips": "Capacitación en ergonomía",
    "is_active": true,
    "order": 3
  },
  {
    "id": "cat_004",
    "name": "Postura deficiente",
    "code": "UNSAFE_POSTURE",
    "type": "unsafe",
    "description": "Postura incorrecta que puede causar lesión",
    "severity": "medium",
    "prevention_tips": "Ejercicios de estiramientos y correcciones",
    "is_active": true,
    "order": 4
  }
]
```

---

## 1.6 Colección: `observations`

**Descripción**: Observaciones SBC realizadas a trabajadores.

### Campos:

| Campo | Tipo | Descripción | Required |
|-------|------|-------------|----------|
| `id` | String (UID) | ID único | ✓ |
| `observation_number` | String | Número secuencial | ✓ |
| `employee_id` | Relation | Empleado observado | ✓ |
| `observer_id` | Relation | Usuario observador | ✓ |
| `work_area_id` | Relation | Área donde se realiza | ✓ |
| `observation_date` | DateTime | Fecha/hora observación | ✓ |
| `task_description` | String | Descripción tarea realizada | ✓ |
| `duration_minutes` | Int | Duración observación (min) | ✓ |
| `has_safe_behaviors` | Bool | Contiene conductas seguras | ✓ |
| `has_unsafe_behaviors` | Bool | Contiene conductas inseguras | ✓ |
| `overall_status` | Select | safe\|unsafe\|mixed | ✓ |
| `comments` | String | Comentarios y observaciones | ✗ |
| `recommendations` | String | Recomendaciones | ✗ |
| `requires_action` | Bool | Requiere acción correctiva | ✓ |
| `requires_incident_report` | Bool | Requiere reporte incidente | ✓ |
| `created_by_id` | Relation | Creado por (auditoría) | ✓ |
| `created_at` | DateTime | Fecha creación | ✓ |
| `updated_at` | DateTime | Fecha actualización | ✓ |

### Índices:
```
CREATE INDEX idx_observations_employee_id ON observations(employee_id);
CREATE INDEX idx_observations_observer_id ON observations(observer_id);
CREATE INDEX idx_observations_work_area_id ON observations(work_area_id);
CREATE INDEX idx_observations_observation_date ON observations(observation_date);
CREATE INDEX idx_observations_overall_status ON observations(overall_status);
CREATE INDEX idx_observations_created_at ON observations(created_at);
```

### Reglas de Acceso:

**CREATE**:
```
@requestUser.role_id in (observer_role_id, supervisor_role_id, admin_role_id)
```

**READ**:
```
@requestUser.is_super_admin = true
OR @requestUser.role_id in (admin_role_id, supervisor_role_id, auditor_role_id)
OR @requestUser.id = observation.observer_id
OR @requestUser.id = observation.employee_id
```

---

## 1.7 Colección: `behaviors`

**Descripción**: Conductas específicas registradas en una observación.

### Campos:

| Campo | Tipo | Descripción | Required |
|-------|------|-------------|----------|
| `id` | String (UID) | ID único | ✓ |
| `observation_id` | Relation | Observación asociada | ✓ |
| `behavior_category_id` | Relation | Categoría conducta | ✓ |
| `type` | Select | safe\|unsafe | ✓ |
| `sequence` | Int | Orden en observación | ✓ |
| `description` | String | Detalles específicos | ✓ |
| `timestamp` | DateTime | Momento dentro observación | ✓ |
| `evidence_photo_id` | Relation | Foto evidencia (opcional) | ✗ |
| `correction_applied` | Bool | Corrección inmediata | ✓ |
| `correction_details` | String | Detalles corrección | ✗ |
| `created_at` | DateTime | Fecha creación | ✓ |

### Índices:
```
CREATE INDEX idx_behaviors_observation_id ON behaviors(observation_id);
CREATE INDEX idx_behaviors_behavior_category_id ON behaviors(behavior_category_id);
CREATE INDEX idx_behaviors_type ON behaviors(type);
```

### Ejemplo JSON:
```json
{
  "id": "beh_001",
  "observation_id": "obs_001",
  "behavior_category_id": "cat_001",
  "type": "safe",
  "sequence": 1,
  "description": "Trabajador usa casco de seguridad correctamente ajustado",
  "timestamp": "2026-05-18T10:15:30.000Z",
  "evidence_photo_id": "att_001",
  "correction_applied": false,
  "created_at": "2026-05-18T10:20:00.000Z"
}
```

---

## 1.8 Colección: `incidents`

**Descripción**: Incidentes de seguridad o casi-accidentes reportados.

### Campos:

| Campo | Tipo | Descripción | Required |
|-------|------|-------------|----------|
| `id` | String (UID) | ID único | ✓ |
| `incident_number` | String | Número secuencial | ✓ |
| `observation_id` | Relation | Observación origen | ✗ |
| `employee_id` | Relation | Empleado involucrado | ✓ |
| `work_area_id` | Relation | Área donde ocurrió | ✓ |
| `incident_date` | DateTime | Fecha/hora incidente | ✓ |
| `incident_type` | Select | accident\|near_miss\|hazard | ✓ |
| `severity` | Select | low\|medium\|high\|critical | ✓ |
| `description` | String | Descripción detallada | ✓ |
| `injuries` | String | Lesiones (si aplica) | ✗ |
| `property_damage` | String | Daño a equipos | ✗ |
| `root_cause` | String | Causa raíz | ✗ |
| `status` | Select | open\|in_progress\|closed | ✓ |
| `reported_by_id` | Relation | Quién reportó | ✓ |
| `assigned_to_id` | Relation | Asignado a | ✗ |
| `investigation_notes` | String | Notas investigación | ✗ |
| `created_at` | DateTime | Fecha creación | ✓ |
| `updated_at` | DateTime | Fecha actualización | ✓ |
| `closed_at` | DateTime | Fecha cierre | ✗ |

### Índices:
```
CREATE INDEX idx_incidents_employee_id ON incidents(employee_id);
CREATE INDEX idx_incidents_work_area_id ON incidents(work_area_id);
CREATE INDEX idx_incidents_incident_date ON incidents(incident_date);
CREATE INDEX idx_incidents_severity ON incidents(severity);
CREATE INDEX idx_incidents_status ON incidents(status);
```

---

## 1.9 Colección: `corrective_actions`

**Descripción**: Acciones correctivas derivadas de observaciones e incidentes.

### Campos:

| Campo | Tipo | Descripción | Required |
|-------|------|-------------|----------|
| `id` | String (UID) | ID único | ✓ |
| `action_number` | String | Número secuencial | ✓ |
| `observation_id` | Relation | Observación origen (opt) | ✗ |
| `incident_id` | Relation | Incidente origen (opt) | ✗ |
| `action_type` | Select | corrective\|preventive\|improvement | ✓ |
| `description` | String | Descripción acción | ✓ |
| `responsible_id` | Relation | Responsable ejecución | ✓ |
| `priority` | Select | low\|medium\|high\|critical | ✓ |
| `due_date` | Date | Fecha vencimiento | ✓ |
| `status` | Select | open\|in_progress\|completed\|overdue | ✓ |
| `completion_date` | DateTime | Fecha completación | ✗ |
| `evidence` | String | Evidencia completación | ✗ |
| `effectiveness_check` | String | Validación efectividad | ✗ |
| `created_by_id` | Relation | Creado por | ✓ |
| `created_at` | DateTime | Fecha creación | ✓ |
| `updated_at` | DateTime | Fecha actualización | ✓ |

### Índices:
```
CREATE INDEX idx_corrective_actions_status ON corrective_actions(status);
CREATE INDEX idx_corrective_actions_due_date ON corrective_actions(due_date);
CREATE INDEX idx_corrective_actions_responsible_id ON corrective_actions(responsible_id);
```

---

## 1.10 Colección: `attachments`

**Descripción**: Archivos adjuntos (fotos, documentos) de observaciones.

### Campos:

| Campo | Tipo | Descripción | Required |
|-------|------|-------------|----------|
| `id` | String (UID) | ID único | ✓ |
| `observation_id` | Relation | Observación asociada | ✗ |
| `incident_id` | Relation | Incidente asociado (opt) | ✗ |
| `corrective_action_id` | Relation | Acción correctiva (opt) | ✗ |
| `file_name` | String | Nombre original archivo | ✓ |
| `file_type` | Select | image\|video\|document\|other | ✓ |
| `mime_type` | String | MIME type | ✓ |
| `file_size` | Int | Tamaño en bytes | ✓ |
| `file_path` | String | Ruta almacenamiento | ✓ |
| `description` | String | Descripción archivo | ✗ |
| `uploaded_by_id` | Relation | Quién subió | ✓ |
| `created_at` | DateTime | Fecha subida | ✓ |

### Restricciones de Seguridad:
```
MAX_FILE_SIZE: 10MB para imágenes, 50MB para videos
ALLOWED_TYPES: jpg, png, gif, pdf, xls, xlsx, doc, docx, mp4
VIRUS_SCAN: Antes de almacenar
ENCRYPTION: Archivos sensibles con AES-256
```

---

## 1.11 Colección: `notifications`

**Descripción**: Notificaciones internas del sistema.

### Campos:

| Campo | Tipo | Descripción | Required |
|-------|------|-------------|----------|
| `id` | String (UID) | ID único | ✓ |
| `recipient_id` | Relation | Usuario destinatario | ✓ |
| `notification_type` | Select | info\|warning\|critical | ✓ |
| `title` | String | Título notificación | ✓ |
| `message` | String | Cuerpo mensaje | ✓ |
| `related_entity_type` | String | Tipo entidad (observation, incident) | ✗ |
| `related_entity_id` | String | ID entidad relacionada | ✗ |
| `is_read` | Bool | Notificación leída | ✓ |
| `read_at` | DateTime | Fecha lectura | ✗ |
| `action_url` | String | URL acción asociada | ✗ |
| `created_at` | DateTime | Fecha creación | ✓ |

### Índices:
```
CREATE INDEX idx_notifications_recipient_id ON notifications(recipient_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);
```

---

## 1.12 Colección: `audit_logs`

**Descripción**: Registro completo de auditoría de todos los cambios del sistema.

### Campos:

| Campo | Tipo | Descripción | Required |
|-------|------|-------------|----------|
| `id` | String (UID) | ID único | ✓ |
| `user_id` | Relation | Usuario que realizó acción | ✓ |
| `action` | String | Acción realizada (create, update, delete, login) | ✓ |
| `entity_type` | String | Tipo entidad afectada | ✓ |
| `entity_id` | String | ID entidad afectada | ✓ |
| `old_values` | JSON | Valores anteriores | ✗ |
| `new_values` | JSON | Valores nuevos | ✗ |
| `change_summary` | String | Resumen cambios | ✗ |
| `ip_address` | String | IP origen | ✓ |
| `user_agent` | String | User Agent navegador | ✓ |
| `status` | Select | success\|failure | ✓ |
| `error_message` | String | Mensaje error (si aplica) | ✗ |
| `created_at` | DateTime | Fecha acción | ✓ |

### Índices:
```
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_entity_type ON audit_logs(entity_type);
CREATE INDEX idx_audit_logs_entity_id ON audit_logs(entity_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
```

### Ejemplo JSON:
```json
{
  "id": "audit_001",
  "user_id": "user_001",
  "action": "update",
  "entity_type": "observations",
  "entity_id": "obs_001",
  "old_values": {"overall_status": "mixed"},
  "new_values": {"overall_status": "safe"},
  "change_summary": "Changed overall status from mixed to safe",
  "ip_address": "192.168.1.100",
  "user_agent": "Mozilla/5.0...",
  "status": "success",
  "created_at": "2026-05-18T14:35:00.000Z"
}
```

---

## 1.13 Colección: `kpis`

**Descripción**: Indicadores clave de desempeño de seguridad.

### Campos:

| Campo | Tipo | Descripción | Required |
|-------|------|-------------|----------|
| `id` | String (UID) | ID único | ✓ |
| `kpi_name` | String | Nombre del KPI | ✓ |
| `kpi_code` | String | Código único | ✓ |
| `period` | Select | daily\|weekly\|monthly\|yearly | ✓ |
| `period_start` | Date | Inicio período | ✓ |
| `period_end` | Date | Fin período | ✓ |
| `work_area_id` | Relation | Área (null = global) | ✗ |
| `total_observations` | Int | Total observaciones | ✓ |
| `safe_behaviors_count` | Int | Conductas seguras | ✓ |
| `unsafe_behaviors_count` | Int | Conductas inseguras | ✓ |
| `safe_percentage` | Float | % conductas seguras | ✓ |
| `unsafe_percentage` | Float | % conductas inseguras | ✓ |
| `incidents_count` | Int | Total incidentes | ✓ |
| `high_severity_incidents` | Int | Incidentes críticos | ✓ |
| `corrective_actions_open` | Int | Acciones abiertas | ✓ |
| `corrective_actions_overdue` | Int | Acciones vencidas | ✓ |
| `corrective_actions_completed` | Int | Acciones completadas | ✓ |
| `observation_frequency` | Float | Observaciones/día | ✓ |
| `employees_observed` | Int | Empleados únicos | ✓ |
| `trend_vs_previous` | String | Tendencia vs período anterior | ✗ |
| `calculated_at` | DateTime | Fecha cálculo | ✓ |

### Índices:
```
CREATE INDEX idx_kpis_period ON kpis(period, period_start);
CREATE INDEX idx_kpis_work_area_id ON kpis(work_area_id);
```

---

## 1.14 Colección: `reports`

**Descripción**: Reportes generados (PDF, Excel).

### Campos:

| Campo | Tipo | Descripción | Required |
|-------|------|-------------|----------|
| `id` | String (UID) | ID único | ✓ |
| `report_name` | String | Nombre reporte | ✓ |
| `report_type` | Select | observations\|incidents\|kpis\|monthly\|custom | ✓ |
| `period_start` | Date | Período inicio | ✓ |
| `period_end` | Date | Período fin | ✓ |
| `work_area_id` | Relation | Área (null = todas) | ✗ |
| `format` | Select | pdf\|excel\|json | ✓ |
| `file_path` | String | Ruta archivo generado | ✓ |
| `file_size` | Int | Tamaño bytes | ✓ |
| `generated_by_id` | Relation | Generado por usuario | ✓ |
| `status` | Select | pending\|generating\|completed\|failed | ✓ |
| `error_message` | String | Error si aplica | ✗ |
| `created_at` | DateTime | Fecha creación | ✓ |
| `generated_at` | DateTime | Fecha generación | ✗ |

### Índices:
```
CREATE INDEX idx_reports_report_type ON reports(report_type);
CREATE INDEX idx_reports_generated_by_id ON reports(generated_by_id);
CREATE INDEX idx_reports_created_at ON reports(created_at);
```

---

## 2. RELACIONES RESUMEN

| De | A | Tipo | Cardinalidad |
|----|---|------|-------------|
| `users` | `roles` | many-to-one | N:1 |
| `users` | `users` | many-to-one | (created_by) N:1 |
| `employees` | `work_areas` | many-to-one | N:1 |
| `employees` | `users` | one-to-one | 1:1 (opcional) |
| `observations` | `employees` | many-to-one | N:1 |
| `observations` | `users` | many-to-one | N:1 (observer) |
| `observations` | `work_areas` | many-to-one | N:1 |
| `behaviors` | `observations` | many-to-one | N:1 |
| `behaviors` | `behavior_categories` | many-to-one | N:1 |
| `behaviors` | `attachments` | many-to-one | N:1 (evidence) |
| `incidents` | `observations` | many-to-one | N:1 (origen) |
| `incidents` | `employees` | many-to-one | N:1 |
| `incidents` | `work_areas` | many-to-one | N:1 |
| `incidents` | `users` | many-to-one | N:1 (reported_by, assigned_to) |
| `corrective_actions` | `observations` | many-to-one | N:1 |
| `corrective_actions` | `incidents` | many-to-one | N:1 |
| `corrective_actions` | `users` | many-to-one | N:1 (responsible, created_by) |
| `attachments` | `observations` | many-to-one | N:1 |
| `attachments` | `incidents` | many-to-one | N:1 |
| `attachments` | `corrective_actions` | many-to-one | N:1 |
| `attachments` | `users` | many-to-one | N:1 (uploaded_by) |
| `notifications` | `users` | many-to-one | N:1 |
| `audit_logs` | `users` | many-to-one | N:1 |
| `kpis` | `work_areas` | many-to-one | N:1 |
| `reports` | `users` | many-to-one | N:1 |
| `reports` | `work_areas` | many-to-one | N:1 |

---

## 3. VALIDACIONES Y RESTRICCIONES

### Por Colección:

**users**:
- Email: Formato válido, único en BD
- Password: Mín 8 chars, mayús, minús, número, símbolo
- Username: 3-50 chars, solo alfanuméricos y guión bajo
- Email y username: únicos

**employees**:
- Document number: Único combinado con tipo
- Email: Formato válido
- Hire date: No puede ser futura

**observations**:
- Duration: Mín 5 minutos, máx 480 minutos
- Status: Si has_safe_behaviors es true Y has_unsafe_behaviors es false → overall_status debe ser "safe"
- Si has_unsafe_behaviors es true → overall_status debe incluir "unsafe"

**behaviors**:
- Tipo debe coincidir con categoría
- Timestamp debe estar dentro de rango observación

**incidents**:
- Incident date: No puede ser futura
- Si severity es critical → requires_investigation debe ser true

**corrective_actions**:
- Due date: Mayor que created_at
- Si status = completed → completion_date es obligatorio
- Si status = completed → evidence es obligatorio

**attachments**:
- File size: Máx 10MB (imagen), 50MB (video)
- MIME type debe ser válido

---

## 4. HOOKS POCKETBASE RECOMENDADOS

Se detallarán en siguiente sección.
