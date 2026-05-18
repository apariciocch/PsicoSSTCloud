# Diseño de Base de Datos - Esquema Relacional

## 1. DIAGRAMA RELACIONAL (ER-UML)

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        DISEÑO ENTIDAD-RELACIÓN                         │
└─────────────────────────────────────────────────────────────────────────┘

                           ┌──────────────┐
                           │    roles     │
                           ├──────────────┤
                           │ id (PK)      │
                           │ name         │
                           │ slug         │
                           │ permissions  │
                           │ level        │
                           └───────┬──────┘
                                   │
                                   │ 1:N
                                   │
                           ┌───────▼──────────┐
                           │     users        │
                           ├──────────────────┤
                           │ id (PK)          │
                           │ email (UNIQUE)   │
                           │ username (UNIQUE)│
                           │ password_hash    │
                           │ first_name       │
                           │ last_name        │
                           │ role_id (FK)     │◄─┐
                           │ created_by_id (FK)  │
                           │ is_active        │  │
                           │ last_login       │  └─ SELF-JOIN
                           └────┬─────────────┘
                                │
                    ┌───────────┼───────────┬──────────┐
                    │           │           │          │
                    │ 1:N       │ 1:N      │ 1:N     │ 1:N
                    │           │           │          │
          ┌─────────▼────┐    ┌─▼────────┐ │    ┌─────▼──────┐
          │ notifications│    │observations │    │ audit_logs  │
          ├──────────────┤    ├────────────┤ │    ├─────────────┤
          │ id (PK)      │    │ id (PK)    │ │    │ id (PK)    │
          │recipient_id..◄───┤observer_id │ │    │ user_id (FK)◄───┐
          │ title        │    │created_by.│ │    │ action      │   │
          │ message      │    │is_read     │ │    │ entity_type │   │
          │ created_at   │    └─┬──────────┘ │    │ created_at  │   │
          └──────────────┘      │            │    └─────────────┘   │
                                │            │                      │
                                │ 1:N        │                      │
                                │            │                      │
                                │     ┌──────▼────────┐             │
                                │     │  work_areas   │             │
                                │     ├───────────────┤             │
                                │     │ id (PK)       │             │
                                │     │ name          │             │
                                │     │ code (UNIQUE) │             │
                                │     │ risk_level    │             │
                                │     │ supervisor_id ◄─────────────┘
                                │     └─┬──────────────┘
                                │       │
                                │       │ 1:N
                                │       │
          ┌─────────────────────┼───────▼─────────────────────┐
          │                     │                             │
          │              ┌──────▼──────────┐                  │
          │              │   employees     │                  │
          │              ├─────────────────┤                  │
          │              │ id (PK)         │                  │
          │              │ document_type   │                  │
          │              │ document_number │                  │
          │              │ first_name      │                  │
          │              │ last_name       │                  │
          │              │ position        │                  │
          │              │ work_area_id (FK)                  │
          │              │ user_id (FK opt)◄───────────────┐  │
          │              └┬────────────────┘                │  │
          │               │                                 │  │
          │               │ 1:N                            │  │
          │               │                                 │  │
          │      ┌────────▼──────────────┐                 │  │
          │      │  observations         │                 │  │
          │      ├───────────────────────┤                 │  │
          │      │ id (PK)               │                 │  │
          │      │ observation_number    │                 │  │
          │      │ employee_id (FK)  ────┼─────────────────┘  │
          │      │ observer_id (FK)  ────┼─────────┐           │
          │      │ work_area_id (FK) ────┼─────────┼───────────┘
          │      │ overall_status        │         │
          │      │ has_safe_behaviors    │         │
          │      │ has_unsafe_behaviors  │         │
          │      │ observation_date      │         │
          │      └┬──────────────────────┘         │
          │       │                                │
          │       │ 1:N                           │
          │       │                                │
          │   ┌───▼────────────────────────────────┼─────┐
          │   │                                    │     │
          │   │                   ┌────────────────▼──┐  │
          │   │                   │ behavior_categories│  │
          │   │                   ├──────────────────┤  │
          │   │                   │ id (PK)          │  │
          │   │                   │ name             │  │
          │   │                   │ code (UNIQUE)    │  │
          │   │                   │ type (safe/unsafe)  │
          │   │                   │ severity         │  │
          │   │                   └──────────────────┘  │
          │   │                          ▲              │
          │   │                          │ 1:N          │
          │   │                          │              │
          │   │                    ┌─────┴──────────┐   │
          │   │                    │    behaviors   │   │
          │   │                    ├────────────────┤   │
          │   │                    │ id (PK)        │   │
          │   │                    │observation_id..────┼─┘
          │   │                    │category_id (FK)◄───┘
          │   │                    │ type           │
          │   │                    │ evidence_id (FK)
          │   │                    └┬───────────────┘
          │   │                     │
          │   │                     │ 1:N
          │   │                     │
          │   │            ┌────────▼──────────┐
          │   │            │   attachments    │
          │   │            ├──────────────────┤
          │   │            │ id (PK)          │
          │   │            │observation_id (FK)
          │   │            │ incident_id (FK opt)
          │   │            │corrective_act(FK opt)
          │   │            │ file_name        │
          │   │            │ file_type        │
          │   │            │ file_path        │
          │   │            │ uploaded_by (FK) │
          │   │            └──────────────────┘
          │   │
          │   │    ┌──────────────────┐
          │   └───►│    incidents     │
          │        ├──────────────────┤
          │        │ id (PK)          │
          │        │ incident_number  │
          │        │ observation_id (FK opt)
          │        │ employee_id (FK) │
          │        │ work_area_id (FK)│
          │        │ incident_type    │
          │        │ severity         │
          │        │ status           │
          │        │ incident_date    │
          │        └┬─────────────────┘
          │         │
          │         │ 1:N
          │         │
          │    ┌────▼──────────────────┐
          │    │ corrective_actions   │
          │    ├─────────────────────┤
          │    │ id (PK)             │
          │    │ action_number       │
          │    │observation_id(FK opt)
          │    │ incident_id (FK opt)│
          │    │ action_type         │
          │    │ responsible_id (FK) │
          │    │ due_date            │
          │    │ status              │
          │    │ completion_date     │
          │    └─────────────────────┘
          │
          │   ┌──────────────┐
          └──►│     kpis     │
              ├──────────────┤
              │ id (PK)      │
              │ kpi_name     │
              │ period       │
              │ work_area_id │
              │ safe_percent │
              │ unsafe_prcnt │
              │ incidents    │
              │ calculated_at│
              └──────────────┘
```

---

## 2. TABLAS NORMALIZADAS (3FN)

### 2.1 Tabla: users
```sql
CREATE TABLE users (
  id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(6)))) || '_' || datetime('now', 'subsec'),
  email TEXT UNIQUE NOT NULL,
  username TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  role_id TEXT NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT 1,
  is_super_admin BOOLEAN NOT NULL DEFAULT 0,
  phone TEXT,
  department TEXT,
  last_login DATETIME,
  login_attempts INTEGER NOT NULL DEFAULT 0,
  locked_until DATETIME,
  mfa_enabled BOOLEAN NOT NULL DEFAULT 0,
  created_by_id TEXT,
  created_at DATETIME NOT NULL DEFAULT (datetime('now')),
  updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
  
  FOREIGN KEY (role_id) REFERENCES roles(id),
  FOREIGN KEY (created_by_id) REFERENCES users(id)
);
```

### 2.2 Tabla: roles
```sql
CREATE TABLE roles (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  slug TEXT UNIQUE NOT NULL,
  description TEXT NOT NULL,
  permissions TEXT NOT NULL,  -- JSON array
  level INTEGER NOT NULL,
  is_system_role BOOLEAN NOT NULL DEFAULT 1,
  created_at DATETIME NOT NULL DEFAULT (datetime('now')),
  updated_at DATETIME NOT NULL DEFAULT (datetime('now'))
);
```

### 2.3 Tabla: work_areas
```sql
CREATE TABLE work_areas (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  code TEXT UNIQUE NOT NULL,
  description TEXT,
  location TEXT NOT NULL,
  risk_level TEXT NOT NULL CHECK (risk_level IN ('alto', 'medio', 'bajo')),
  supervisor_id TEXT,
  is_active BOOLEAN NOT NULL DEFAULT 1,
  created_at DATETIME NOT NULL DEFAULT (datetime('now')),
  updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
  
  FOREIGN KEY (supervisor_id) REFERENCES users(id)
);
```

### 2.4 Tabla: employees
```sql
CREATE TABLE employees (
  id TEXT PRIMARY KEY,
  document_type TEXT NOT NULL CHECK (document_type IN ('CC', 'CE', 'PASSPORT')),
  document_number TEXT NOT NULL,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  email TEXT,
  phone TEXT,
  position TEXT NOT NULL,
  department TEXT NOT NULL,
  work_area_id TEXT NOT NULL,
  user_id TEXT,
  hire_date DATE NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT 1,
  notes TEXT,
  created_at DATETIME NOT NULL DEFAULT (datetime('now')),
  updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
  
  FOREIGN KEY (work_area_id) REFERENCES work_areas(id),
  FOREIGN KEY (user_id) REFERENCES users(id),
  UNIQUE (document_type, document_number)
);
```

### 2.5 Tabla: observations
```sql
CREATE TABLE observations (
  id TEXT PRIMARY KEY,
  observation_number TEXT UNIQUE NOT NULL,
  employee_id TEXT NOT NULL,
  observer_id TEXT NOT NULL,
  work_area_id TEXT NOT NULL,
  observation_date DATETIME NOT NULL,
  task_description TEXT NOT NULL,
  duration_minutes INTEGER NOT NULL,
  has_safe_behaviors BOOLEAN NOT NULL,
  has_unsafe_behaviors BOOLEAN NOT NULL,
  overall_status TEXT NOT NULL CHECK (overall_status IN ('safe', 'unsafe', 'mixed')),
  comments TEXT,
  recommendations TEXT,
  requires_action BOOLEAN NOT NULL DEFAULT 0,
  requires_incident_report BOOLEAN NOT NULL DEFAULT 0,
  created_by_id TEXT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT (datetime('now')),
  updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
  
  FOREIGN KEY (employee_id) REFERENCES employees(id),
  FOREIGN KEY (observer_id) REFERENCES users(id),
  FOREIGN KEY (work_area_id) REFERENCES work_areas(id),
  FOREIGN KEY (created_by_id) REFERENCES users(id)
);
```

### 2.6 Tabla: behavior_categories
```sql
CREATE TABLE behavior_categories (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  code TEXT UNIQUE NOT NULL,
  type TEXT NOT NULL CHECK (type IN ('safe', 'unsafe')),
  description TEXT NOT NULL,
  severity TEXT NOT NULL CHECK (severity IN ('low', 'medium', 'high', 'critical')),
  prevention_tips TEXT,
  is_active BOOLEAN NOT NULL DEFAULT 1,
  order INTEGER NOT NULL,
  created_at DATETIME NOT NULL DEFAULT (datetime('now'))
);
```

### 2.7 Tabla: behaviors
```sql
CREATE TABLE behaviors (
  id TEXT PRIMARY KEY,
  observation_id TEXT NOT NULL,
  behavior_category_id TEXT NOT NULL,
  type TEXT NOT NULL CHECK (type IN ('safe', 'unsafe')),
  sequence INTEGER NOT NULL,
  description TEXT NOT NULL,
  timestamp DATETIME NOT NULL,
  evidence_photo_id TEXT,
  correction_applied BOOLEAN NOT NULL DEFAULT 0,
  correction_details TEXT,
  created_at DATETIME NOT NULL DEFAULT (datetime('now')),
  
  FOREIGN KEY (observation_id) REFERENCES observations(id),
  FOREIGN KEY (behavior_category_id) REFERENCES behavior_categories(id),
  FOREIGN KEY (evidence_photo_id) REFERENCES attachments(id)
);
```

### 2.8 Tabla: incidents
```sql
CREATE TABLE incidents (
  id TEXT PRIMARY KEY,
  incident_number TEXT UNIQUE NOT NULL,
  observation_id TEXT,
  employee_id TEXT NOT NULL,
  work_area_id TEXT NOT NULL,
  incident_date DATETIME NOT NULL,
  incident_type TEXT NOT NULL CHECK (incident_type IN ('accident', 'near_miss', 'hazard')),
  severity TEXT NOT NULL CHECK (severity IN ('low', 'medium', 'high', 'critical')),
  description TEXT NOT NULL,
  injuries TEXT,
  property_damage TEXT,
  root_cause TEXT,
  status TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'in_progress', 'closed')),
  reported_by_id TEXT NOT NULL,
  assigned_to_id TEXT,
  investigation_notes TEXT,
  created_at DATETIME NOT NULL DEFAULT (datetime('now')),
  updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
  closed_at DATETIME,
  
  FOREIGN KEY (observation_id) REFERENCES observations(id),
  FOREIGN KEY (employee_id) REFERENCES employees(id),
  FOREIGN KEY (work_area_id) REFERENCES work_areas(id),
  FOREIGN KEY (reported_by_id) REFERENCES users(id),
  FOREIGN KEY (assigned_to_id) REFERENCES users(id)
);
```

### 2.9 Tabla: corrective_actions
```sql
CREATE TABLE corrective_actions (
  id TEXT PRIMARY KEY,
  action_number TEXT UNIQUE NOT NULL,
  observation_id TEXT,
  incident_id TEXT,
  action_type TEXT NOT NULL CHECK (action_type IN ('corrective', 'preventive', 'improvement')),
  description TEXT NOT NULL,
  responsible_id TEXT NOT NULL,
  priority TEXT NOT NULL CHECK (priority IN ('low', 'medium', 'high', 'critical')),
  due_date DATE NOT NULL,
  status TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'in_progress', 'completed', 'overdue')),
  completion_date DATETIME,
  evidence TEXT,
  effectiveness_check TEXT,
  created_by_id TEXT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT (datetime('now')),
  updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
  
  FOREIGN KEY (observation_id) REFERENCES observations(id),
  FOREIGN KEY (incident_id) REFERENCES incidents(id),
  FOREIGN KEY (responsible_id) REFERENCES users(id),
  FOREIGN KEY (created_by_id) REFERENCES users(id)
);
```

### 2.10 Tabla: attachments
```sql
CREATE TABLE attachments (
  id TEXT PRIMARY KEY,
  observation_id TEXT,
  incident_id TEXT,
  corrective_action_id TEXT,
  file_name TEXT NOT NULL,
  file_type TEXT NOT NULL CHECK (file_type IN ('image', 'video', 'document', 'other')),
  mime_type TEXT NOT NULL,
  file_size INTEGER NOT NULL,
  file_path TEXT NOT NULL,
  description TEXT,
  uploaded_by_id TEXT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT (datetime('now')),
  
  FOREIGN KEY (observation_id) REFERENCES observations(id),
  FOREIGN KEY (incident_id) REFERENCES incidents(id),
  FOREIGN KEY (corrective_action_id) REFERENCES corrective_actions(id),
  FOREIGN KEY (uploaded_by_id) REFERENCES users(id)
);
```

### 2.11 Tabla: notifications
```sql
CREATE TABLE notifications (
  id TEXT PRIMARY KEY,
  recipient_id TEXT NOT NULL,
  notification_type TEXT NOT NULL CHECK (notification_type IN ('info', 'warning', 'critical')),
  title TEXT NOT NULL,
  message TEXT NOT NULL,
  related_entity_type TEXT,
  related_entity_id TEXT,
  is_read BOOLEAN NOT NULL DEFAULT 0,
  read_at DATETIME,
  action_url TEXT,
  created_at DATETIME NOT NULL DEFAULT (datetime('now')),
  
  FOREIGN KEY (recipient_id) REFERENCES users(id)
);
```

### 2.12 Tabla: audit_logs
```sql
CREATE TABLE audit_logs (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  action TEXT NOT NULL CHECK (action IN ('create', 'update', 'delete', 'login', 'logout')),
  entity_type TEXT NOT NULL,
  entity_id TEXT NOT NULL,
  old_values TEXT,  -- JSON
  new_values TEXT,  -- JSON
  change_summary TEXT,
  ip_address TEXT NOT NULL,
  user_agent TEXT NOT NULL,
  status TEXT NOT NULL CHECK (status IN ('success', 'failure')),
  error_message TEXT,
  created_at DATETIME NOT NULL DEFAULT (datetime('now')),
  
  FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### 2.13 Tabla: kpis
```sql
CREATE TABLE kpis (
  id TEXT PRIMARY KEY,
  kpi_name TEXT NOT NULL,
  kpi_code TEXT UNIQUE NOT NULL,
  period TEXT NOT NULL CHECK (period IN ('daily', 'weekly', 'monthly', 'yearly')),
  period_start DATE NOT NULL,
  period_end DATE NOT NULL,
  work_area_id TEXT,
  total_observations INTEGER NOT NULL,
  safe_behaviors_count INTEGER NOT NULL,
  unsafe_behaviors_count INTEGER NOT NULL,
  safe_percentage REAL NOT NULL,
  unsafe_percentage REAL NOT NULL,
  incidents_count INTEGER NOT NULL,
  high_severity_incidents INTEGER NOT NULL,
  corrective_actions_open INTEGER NOT NULL,
  corrective_actions_overdue INTEGER NOT NULL,
  corrective_actions_completed INTEGER NOT NULL,
  observation_frequency REAL NOT NULL,
  employees_observed INTEGER NOT NULL,
  trend_vs_previous TEXT,
  calculated_at DATETIME NOT NULL DEFAULT (datetime('now')),
  
  FOREIGN KEY (work_area_id) REFERENCES work_areas(id)
);
```

### 2.14 Tabla: reports
```sql
CREATE TABLE reports (
  id TEXT PRIMARY KEY,
  report_name TEXT NOT NULL,
  report_type TEXT NOT NULL CHECK (report_type IN ('observations', 'incidents', 'kpis', 'monthly', 'custom')),
  period_start DATE NOT NULL,
  period_end DATE NOT NULL,
  work_area_id TEXT,
  format TEXT NOT NULL CHECK (format IN ('pdf', 'excel', 'json')),
  file_path TEXT NOT NULL,
  file_size INTEGER NOT NULL,
  generated_by_id TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'generating', 'completed', 'failed')),
  error_message TEXT,
  created_at DATETIME NOT NULL DEFAULT (datetime('now')),
  generated_at DATETIME,
  
  FOREIGN KEY (generated_by_id) REFERENCES users(id),
  FOREIGN KEY (work_area_id) REFERENCES work_areas(id)
);
```

---

## 3. ÍNDICES ESTRATÉGICOS

```sql
-- users
CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE UNIQUE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_role_id ON users(role_id);
CREATE INDEX idx_users_is_active ON users(is_active);
CREATE INDEX idx_users_created_at ON users(created_at);

-- work_areas
CREATE UNIQUE INDEX idx_work_areas_code ON work_areas(code);
CREATE INDEX idx_work_areas_risk_level ON work_areas(risk_level);

-- employees
CREATE UNIQUE INDEX idx_employees_document ON employees(document_type, document_number);
CREATE INDEX idx_employees_work_area_id ON employees(work_area_id);
CREATE INDEX idx_employees_is_active ON employees(is_active);
CREATE INDEX idx_employees_created_at ON employees(created_at);

-- observations
CREATE INDEX idx_observations_employee_id ON observations(employee_id);
CREATE INDEX idx_observations_observer_id ON observations(observer_id);
CREATE INDEX idx_observations_work_area_id ON observations(work_area_id);
CREATE INDEX idx_observations_observation_date ON observations(observation_date DESC);
CREATE INDEX idx_observations_overall_status ON observations(overall_status);
CREATE INDEX idx_observations_created_at ON observations(created_at DESC);

-- behaviors
CREATE INDEX idx_behaviors_observation_id ON behaviors(observation_id);
CREATE INDEX idx_behaviors_behavior_category_id ON behaviors(behavior_category_id);
CREATE INDEX idx_behaviors_type ON behaviors(type);

-- incidents
CREATE INDEX idx_incidents_employee_id ON incidents(employee_id);
CREATE INDEX idx_incidents_work_area_id ON incidents(work_area_id);
CREATE INDEX idx_incidents_incident_date ON incidents(incident_date DESC);
CREATE INDEX idx_incidents_severity ON incidents(severity);
CREATE INDEX idx_incidents_status ON incidents(status);

-- corrective_actions
CREATE INDEX idx_corrective_actions_status ON corrective_actions(status);
CREATE INDEX idx_corrective_actions_due_date ON corrective_actions(due_date);
CREATE INDEX idx_corrective_actions_responsible_id ON corrective_actions(responsible_id);

-- notifications
CREATE INDEX idx_notifications_recipient_id ON notifications(recipient_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at DESC);

-- audit_logs
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_entity_type ON audit_logs(entity_type);
CREATE INDEX idx_audit_logs_entity_id ON audit_logs(entity_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);

-- kpis
CREATE INDEX idx_kpis_period ON kpis(period, period_start);
CREATE INDEX idx_kpis_work_area_id ON kpis(work_area_id);

-- reports
CREATE INDEX idx_reports_report_type ON reports(report_type);
CREATE INDEX idx_reports_generated_by_id ON reports(generated_by_id);
CREATE INDEX idx_reports_created_at ON reports(created_at DESC);
```

---

## 4. CARDINALIDAD Y RELACIONES

| Relación | De | A | Cardinalidad | Tipo | Descripción |
|----------|-----|-----|-------------|------|------------|
| Role Assignment | roles | users | 1:N | Required | Cada usuario tiene exactamente un rol |
| User Creation | users | users | 1:N | Optional | Auditoría de quién creó cada usuario |
| Area Supervision | users | work_areas | 1:N | Optional | Un supervisor puede supervisar múltiples áreas |
| Area Assignment | work_areas | employees | 1:N | Required | Cada empleado pertenece a un área |
| User Binding | users | employees | 1:1 | Optional | Empleado puede tener cuenta de usuario |
| Observations | employees | observations | 1:N | Required | Cada observación es de un empleado |
| Observer | users | observations | 1:N | Required | Cada observación es de un observador |
| Area Observation | work_areas | observations | 1:N | Required | Observaciones se realizan en áreas |
| Behaviors | observations | behaviors | 1:N | Required | Una observación tiene múltiples conductas |
| Category | behavior_categories | behaviors | 1:N | Required | Conductas se clasifican en categorías |
| Evidence | attachments | behaviors | 1:N | Optional | Conductas pueden tener foto evidencia |
| Incident Origin | observations | incidents | 1:N | Optional | Incidente puede derivar de observación |
| Incident Employee | employees | incidents | 1:N | Required | Incidente involucra un empleado |
| Incident Area | work_areas | incidents | 1:N | Required | Incidente ocurre en área |
| Incident Reporter | users | incidents | 1:N | Required | Alguien reporta el incidente |
| Incident Assignee | users | incidents | 1:N | Optional | Incidente asignado a responsable |
| Corrective Origin | observations | corrective_actions | 1:N | Optional | Acción correctiva de observación |
| Corrective Incident | incidents | corrective_actions | 1:N | Optional | Acción correctiva de incidente |
| Action Responsible | users | corrective_actions | 1:N | Required | Usuario responsable de acción |
| Attachments | observations | attachments | 1:N | Optional | Observación puede tener archivos |

---

## 5. INTEGRIDAD REFERENCIAL

Todas las Foreign Keys utilizan:
- **ON DELETE**: RESTRICT (prevenir eliminación si hay referencias activas)
- **ON UPDATE**: CASCADE (actualizar referencias si cambia la clave primaria)
- **DEFER CONSTRAINTS**: Algunos triggers pueden diferir validaciones

---

## 6. CONSTRAINTS Y VALIDACIONES

### Restricciones a Nivel de BD:

1. **users.email**: UNIQUE, NOT NULL, formato email válido
2. **users.username**: UNIQUE, NOT NULL, 3-50 chars
3. **users.password_hash**: NOT NULL, minLength 60 (bcrypt)
4. **observations.observation_number**: UNIQUE, auto-generado
5. **incidents.incident_number**: UNIQUE, auto-generado
6. **corrective_actions.action_number**: UNIQUE, auto-generado
7. **employees.document_type + document_number**: UNIQUE combinado

### Constraints de Negocio:

1. No se puede asignar rol eliminado a usuario
2. Empleado activo debe estar en área activa
3. Observación no puede tener fecha futura
4. Acción correctiva vencida → status = "overdue"
5. Incidente crítico automáticamente asignado a supervisor

---

## 7. VISTA ANALÍTICA (Materialized View)

```sql
CREATE VIEW vw_observation_summary AS
SELECT
  o.id,
  o.observation_number,
  e.first_name || ' ' || e.last_name AS employee_name,
  u.first_name || ' ' || u.last_name AS observer_name,
  wa.name AS work_area_name,
  o.observation_date,
  COUNT(CASE WHEN b.type = 'safe' THEN 1 END) AS safe_count,
  COUNT(CASE WHEN b.type = 'unsafe' THEN 1 END) AS unsafe_count,
  o.overall_status,
  o.requires_action,
  o.created_at
FROM observations o
JOIN employees e ON o.employee_id = e.id
JOIN users u ON o.observer_id = u.id
JOIN work_areas wa ON o.work_area_id = wa.id
LEFT JOIN behaviors b ON o.id = b.observation_id
GROUP BY o.id;
```

---

## 8. ESTRATEGIA DE BACKUP Y RECUPERACIÓN

- **Frecuencia**: Diaria (incremental), semanal (completo)
- **Destino**: AWS S3 / Google Cloud Storage
- **Retención**: 30 días mínimo
- **Testing**: Recuperación de prueba mensual
- **Point-in-Time Recovery**: Disponible últimas 7 días
