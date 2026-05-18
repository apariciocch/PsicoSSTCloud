# API REST Endpoints - SBC Backend

## Versión: v1

Base URL: `https://api.empresa.com/api/v1`

---

## 1. AUTENTICACIÓN

### POST /auth/login
Autentica un usuario y retorna tokens JWT.

**Request:**
```json
{
  "email": "usuario@empresa.com",
  "password": "SecurePassword123!"
}
```

**Response (200):**
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "expires_in": 900,
  "user": {
    "id": "user_123",
    "email": "usuario@empresa.com",
    "username": "usuario",
    "first_name": "Juan",
    "last_name": "Pérez",
    "role_id": "role_observer",
    "is_active": true,
    "is_super_admin": false,
    "created_at": "2026-01-15T08:00:00Z"
  }
}
```

**Errores:**
- `401 AUTH_001`: Credenciales inválidas
- `400 VALIDATION_ERROR`: Campos faltantes

---

### POST /auth/refresh
Renueva el access token usando un refresh token.

**Request:**
```json
{
  "refresh_token": "eyJhbGc..."
}
```

**Response (200):**
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "expires_in": 900,
  "user": { ... }
}
```

---

### POST /auth/logout
Cierra la sesión del usuario.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200):**
```json
{
  "message": "Logged out successfully"
}
```

---

### POST /auth/change-password
Cambia la contraseña del usuario.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request:**
```json
{
  "current_password": "OldPassword123!",
  "new_password": "NewPassword456!",
  "confirm_password": "NewPassword456!"
}
```

**Response (200):**
```json
{
  "message": "Password changed successfully"
}
```

---

## 2. USUARIOS

### GET /users
Obtiene lista de usuarios (solo admin).

**Headers:**
```
Authorization: Bearer <access_token>
```

**Query Parameters:**
- `page` (int): Número de página (default: 1)
- `page_size` (int): Registros por página (default: 20, máx: 100)
- `role_id` (string): Filtrar por rol
- `is_active` (boolean): Filtrar por estado

**Response (200):**
```json
{
  "data": [
    {
      "id": "user_123",
      "email": "usuario@empresa.com",
      "username": "usuario",
      "first_name": "Juan",
      "last_name": "Pérez",
      "role_id": "role_observer",
      "is_active": true,
      "created_at": "2026-01-15T08:00:00Z"
    }
  ],
  "page": 1,
  "page_size": 20,
  "total": 150,
  "total_pages": 8,
  "has_next": true,
  "has_prev": false
}
```

---

### POST /users
Crea un nuevo usuario.

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request:**
```json
{
  "email": "nuevo@empresa.com",
  "username": "nuevo",
  "password": "SecurePassword123!",
  "first_name": "Carlos",
  "last_name": "López",
  "role_id": "role_observer",
  "phone": "+57-300-555-0123",
  "department": "Operaciones"
}
```

**Response (201):**
```json
{
  "id": "user_124",
  "email": "nuevo@empresa.com",
  "username": "nuevo",
  "first_name": "Carlos",
  "last_name": "López",
  "role_id": "role_observer",
  "is_active": true,
  "created_at": "2026-05-18T14:30:00Z"
}
```

---

### GET /users/{id}
Obtiene detalles de un usuario específico.

---

### PUT /users/{id}
Actualiza datos de un usuario.

**Request:**
```json
{
  "first_name": "Carlos",
  "last_name": "López González",
  "phone": "+57-300-666-0123",
  "department": "Recursos Humanos"
}
```

---

### DELETE /users/{id}
Elimina un usuario (solo admin).

---

## 3. EMPLEADOS

### GET /employees
Obtiene lista de empleados.

**Query Parameters:**
- `page` (int): Número de página
- `page_size` (int): Registros por página
- `work_area_id` (string): Filtrar por área
- `is_active` (boolean): Filtrar por estado

---

### POST /employees
Crea un nuevo empleado.

**Request:**
```json
{
  "document_type": "CC",
  "document_number": "1234567890",
  "first_name": "María",
  "last_name": "González",
  "email": "maria.gonzalez@empresa.com",
  "phone": "+57-300-555-0123",
  "position": "Operario",
  "department": "Producción",
  "work_area_id": "area_001",
  "hire_date": "2025-01-15T00:00:00Z"
}
```

---

## 4. OBSERVACIONES

### GET /observations
Obtiene lista de observaciones.

**Query Parameters:**
- `page` (int): Número de página
- `page_size` (int): Registros por página
- `employee_id` (string): Filtrar por empleado
- `work_area_id` (string): Filtrar por área
- `observer_id` (string): Filtrar por observador
- `overall_status` (string): safe, unsafe, mixed
- `start_date` (string): ISO 8601 date
- `end_date` (string): ISO 8601 date

---

### POST /observations
Crea una nueva observación.

**Request:**
```json
{
  "employee_id": "emp_123",
  "work_area_id": "area_001",
  "task_description": "Operación de maquinaria pesada",
  "duration_minutes": 30,
  "has_safe_behaviors": true,
  "has_unsafe_behaviors": false,
  "comments": "Observación realizada sin incidentes",
  "requires_action": false,
  "requires_incident_report": false,
  "behaviors": [
    {
      "behavior_category_id": "cat_001",
      "type": "safe",
      "sequence": 1,
      "description": "Utilizó casco de seguridad correctamente",
      "timestamp": "2026-05-18T10:15:00Z",
      "correction_applied": false
    }
  ]
}
```

**Response (201):**
```json
{
  "id": "obs_001",
  "observation_number": "OBS-2026-00001",
  "employee_id": "emp_123",
  "observer_id": "user_123",
  "work_area_id": "area_001",
  "overall_status": "safe",
  "safe_count": 2,
  "unsafe_count": 0,
  "task_description": "Operación de maquinaria pesada",
  "duration_minutes": 30,
  "created_at": "2026-05-18T10:20:00Z"
}
```

---

### GET /observations/{id}
Obtiene detalles de una observación específica.

---

### PUT /observations/{id}
Actualiza una observación existente.

---

### DELETE /observations/{id}
Elimina una observación (solo creador o admin).

---

## 5. INCIDENTES

### GET /incidents
Obtiene lista de incidentes.

**Query Parameters:**
- `page`, `page_size`
- `employee_id`, `work_area_id`
- `severity`: low, medium, high, critical
- `status`: open, in_progress, closed
- `start_date`, `end_date`

---

### POST /incidents
Crea un nuevo incidente.

**Request:**
```json
{
  "observation_id": "obs_001",
  "employee_id": "emp_123",
  "work_area_id": "area_001",
  "incident_date": "2026-05-18T10:15:00Z",
  "incident_type": "accident",
  "severity": "high",
  "description": "El trabajador se golpeó la mano en la maquinaria",
  "injuries": "Contusión en mano derecha",
  "property_damage": "Ninguno",
  "requires_incident_report": true
}
```

---

## 6. ACCIONES CORRECTIVAS

### GET /corrective-actions
Obtiene lista de acciones correctivas.

**Query Parameters:**
- `status`: open, in_progress, completed, overdue
- `priority`: low, medium, high, critical
- `due_date` (range filter)

---

### POST /corrective-actions
Crea una nueva acción correctiva.

**Request:**
```json
{
  "observation_id": "obs_001",
  "action_type": "corrective",
  "description": "Capacitación en uso de EPP",
  "responsible_id": "user_supervisor",
  "priority": "high",
  "due_date": "2026-05-25T00:00:00Z"
}
```

---

## 7. ANALYTICS

### GET /analytics/dashboard
Obtiene el dashboard estadístico.

**Query Parameters:**
- `period` (string): day, week, month, year (default: month)
- `work_area_id` (string): Filtrar por área (opcional, null = todas)

**Response (200):**
```json
{
  "period": "month",
  "period_start": "2026-05-01T00:00:00Z",
  "period_end": "2026-05-31T23:59:59Z",
  "total_observations": 145,
  "safe_behaviors_percent": 78.5,
  "unsafe_behaviors_percent": 21.5,
  "total_incidents": 8,
  "critical_incidents": 1,
  "overdue_actions": 3,
  "completed_actions": 12,
  "open_actions": 7,
  "trend": "up",
  "incidents_by_area": [
    {
      "area_id": "area_001",
      "area_name": "Producción",
      "count": 5,
      "severity": "high"
    }
  ],
  "critical_areas": [
    {
      "area_id": "area_001",
      "area_name": "Producción",
      "unsafe_percent": 35.2,
      "incident_count": 5,
      "risk_level": "alto"
    }
  ],
  "most_observed_areas": [
    {
      "area_id": "area_001",
      "area_name": "Producción",
      "observation_count": 87,
      "safe_percent": 72.4,
      "incident_count": 5
    }
  ],
  "top_behaviors": [
    {
      "behavior_id": "cat_001",
      "behavior_name": "No usar EPP",
      "type": "unsafe",
      "occurrence_count": 22,
      "severity": "critical"
    }
  ]
}
```

---

### GET /analytics/kpis
Obtiene KPIs agregados.

---

### GET /analytics/trends
Obtiene tendencias históricas.

---

## 8. REPORTES

### POST /reports
Genera un nuevo reporte.

**Request:**
```json
{
  "report_type": "observations",
  "period_start": "2026-05-01",
  "period_end": "2026-05-31",
  "work_area_id": "area_001",
  "format": "pdf"
}
```

**Response (202):**
```json
{
  "id": "rep_001",
  "status": "pending",
  "report_type": "observations",
  "created_at": "2026-05-18T14:30:00Z"
}
```

---

### GET /reports/{id}/download
Descarga un reporte generado.

---

## 9. AUDIT LOGS

### GET /audit-logs
Obtiene logs de auditoría.

**Query Parameters:**
- `page`, `page_size`
- `user_id`: Filtrar por usuario
- `entity_type`: Tipo de entidad
- `action`: create, update, delete, login
- `start_date`, `end_date`

---

## 10. NOTIFICACIONES

### GET /notifications
Obtiene notificaciones del usuario.

---

### GET /notifications/unread
Obtiene notificaciones no leídas.

---

### PUT /notifications/{id}/read
Marca una notificación como leída.

---

### PUT /notifications/mark-all-as-read
Marca todas las notificaciones como leídas.

---

## Códigos de Error Estándar

| Código | HTTP Status | Descripción |
|--------|-------------|------------|
| `AUTH_001` | 401 | Credenciales inválidas |
| `AUTH_002` | 401 | Token expirado |
| `AUTH_003` | 401 | Formato de token inválido |
| `AUTH_004` | 403 | Acceso denegado / Permisos insuficientes |
| `VALIDATION_ERROR` | 400 | Error de validación |
| `NOT_FOUND` | 404 | Recurso no encontrado |
| `ALREADY_EXISTS` | 409 | El recurso ya existe |
| `SYSTEM_001` | 500 | Error interno del servidor |
| `DB_001` | 500 | Error de base de datos |

---

## Headers Requeridos

Todos los endpoints (excepto `/auth/login` y `/auth/refresh`) requieren:

```
Authorization: Bearer <access_token>
Content-Type: application/json
```

---

## Rate Limiting

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1716036600
```

Si se excede el límite:
```
HTTP/1.1 429 Too Many Requests
Retry-After: 60
```

---

## Paginación

**Query Parameters:**
```
?page=1&page_size=20
```

**Response Fields:**
- `page`: Página actual
- `page_size`: Registros por página
- `total`: Total de registros
- `total_pages`: Número total de páginas
- `has_next`: Indica si hay siguiente página
- `has_prev`: Indica si hay página anterior

---

## Filtros Comunes

**Date Range:**
```
?start_date=2026-05-01&end_date=2026-05-31
```

**Status Filter:**
```
?status=open&status=in_progress
```

**Sorting (cuando esté implementado):**
```
?sort=created_at&order=desc
```
