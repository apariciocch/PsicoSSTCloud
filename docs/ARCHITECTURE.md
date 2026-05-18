# Arquitectura del Sistema SBC - PsicoSST Cloud

## 1. Visión General

Sistema de gestión integral de Seguridad Basada en el Comportamiento (SBC) con arquitectura modular, escalable y segura, construido con Go, PocketBase y principios de Clean Architecture.

### Objetivos Arquitectónicos:
- **Modularidad**: Capas independientes y reutilizables
- **Escalabilidad**: Preparado para crecimiento empresarial
- **Seguridad**: Autenticación JWT, encriptación, auditoría completa
- **Mantenibilidad**: Código limpio, documentado y testeable
- **Performance**: Queries optimizadas, índices estratégicos
- **Observabilidad**: Logs estructurados, trazabilidad de cambios

---

## 2. Arquitectura de Capas

```
┌─────────────────────────────────────────────────────────┐
│                    Capa de Presentación                 │
│              (Frontend - No incluido en scope)          │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│                 Capa API REST (Routers)                 │
│  Endpoints, Validación de entrada, Serialización JSON  │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│              Capa de Middlewares                        │
│     JWT Auth, Rate Limiting, CORS, Logging, RBAC       │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│          Capa de Controladores (Use Cases)              │
│     Orquestación de lógica, validación de negocio      │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│          Capa de Servicios (Business Logic)             │
│     Reglas de negocio, cálculos, transformaciones      │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│        Capa de Repositorio (Data Access)                │
│       Queries PocketBase, abstracción de datos         │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│            Capa de Modelos/Entidades                    │
│              Estructuras de datos, DTOs                │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│             PocketBase + SQLite/PostgreSQL              │
│         Persistencia, relaciones, índices              │
└─────────────────────────────────────────────────────────┘
```

---

## 3. Estructura de Carpetas del Backend

```
backend/
├── cmd/
│   └── main.go                    # Punto de entrada
├── internal/
│   ├── config/
│   │   ├── config.go              # Configuración env
│   │   └── constants.go           # Constantes del sistema
│   ├── domain/
│   │   ├── entities/              # Modelos de negocio
│   │   ├── interfaces/            # Contratos
│   │   └── errors/                # Errores personalizados
│   ├── usecase/
│   │   ├── auth/                  # Casos de uso autenticación
│   │   ├── user/                  # Casos de uso usuarios
│   │   ├── employee/              # Casos de uso empleados
│   │   ├── observation/           # Casos de uso observaciones
│   │   ├── behavior/              # Casos de uso conductas
│   │   ├── incident/              # Casos de uso incidentes
│   │   ├── report/                # Casos de uso reportes
│   │   └── analytics/             # Casos de uso analytics
│   ├── handler/
│   │   ├── auth/                  # Controladores auth
│   │   ├── user/                  # Controladores usuarios
│   │   ├── employee/              # Controladores empleados
│   │   ├── observation/           # Controladores observaciones
│   │   ├── behavior/              # Controladores conductas
│   │   ├── incident/              # Controladores incidentes
│   │   ├── report/                # Controladores reportes
│   │   └── analytics/             # Controladores analytics
│   ├── repository/
│   │   ├── user/                  # Repo usuarios
│   │   ├── employee/              # Repo empleados
│   │   ├── observation/           # Repo observaciones
│   │   ├── behavior/              # Repo conductas
│   │   └── pocketbase.go          # Cliente PocketBase
│   ├── service/
│   │   ├── auth/                  # Servicio autenticación
│   │   ├── jwt/                   # Servicio JWT
│   │   ├── password/              # Servicio hash contraseñas
│   │   ├── observation/           # Servicio observaciones
│   │   ├── analytics/             # Servicio analytics
│   │   └── email/                 # Servicio notificaciones
│   ├── middleware/
│   │   ├── auth.go                # JWT middleware
│   │   ├── rbac.go                # Control permisos
│   │   ├── logger.go              # Logging middleware
│   │   ├── cors.go                # CORS
│   │   └── ratelimit.go           # Rate limiting
│   ├── router/
│   │   └── routes.go              # Definición de rutas
│   ├── dto/                       # Data Transfer Objects
│   │   ├── auth.go
│   │   ├── user.go
│   │   ├── observation.go
│   │   └── ...
│   └── pkg/
│       ├── logger/                # Logger estructurado
│       ├── validator/             # Validaciones
│       ├── response/              # Formatos respuesta
│       ├── utils/                 # Utilidades
│       └── crypto/                # Encriptación
├── migrations/
│   └── pocketbase_init.go        # Inicialización PB
├── tests/
│   ├── unit/
│   ├── integration/
│   └── fixtures/
├── .env.example                   # Variables de entorno
├── .dockerignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

---

## 4. Componentes Clave

### 4.1 Sistema de Autenticación
- **JWT**: Token de corta duración (15 min), Refresh token (7 días)
- **Roles y Permisos**: RBAC granular por endpoint
- **Auditoría**: Registro de login, logout, cambios de contraseña

### 4.2 Validaciones
- **Request Validation**: DTO con validación de etiquetas
- **Data Integrity**: Restricciones en base de datos
- **Business Rules**: Validación en capa de servicio

### 4.3 Logging y Monitoreo
- **Structured Logging**: JSON con contexto
- **Request/Response Logging**: ID de trace para seguimiento
- **Audit Trail**: Todos los cambios registrados

### 4.4 Seguridad
- **Password Hashing**: bcrypt con salt
- **Input Sanitization**: Protección contra inyección SQL
- **CORS**: Configuración restrictiva
- **Rate Limiting**: Protección contra fuerza bruta
- **Encriptación de datos sensibles**: AES-256 para campos críticos

---

## 5. Flujos de Datos Principales

### 5.1 Autenticación
```
1. POST /api/auth/login (email, password)
2. Servicio Auth valida credenciales
3. JWT Service genera tokens
4. Response: accessToken, refreshToken, user
5. Cliente almacena y envía accessToken en headers
6. Middleware JWT valida en cada request
```

### 5.2 Creación de Observación
```
1. POST /api/observations (datos observación)
2. Handler valida entrada
3. Controlador verifica permisos
4. Servicio aplica reglas de negocio
5. Repositorio persiste en PocketBase
6. Audit log registra creación
7. Response: observación creada con ID
8. Evento desencadena análisis automático
```

### 5.3 Generación de Dashboard
```
1. GET /api/analytics/dashboard?period=month
2. Handler valida parámetros
3. Controlador obtiene datos del servicio
4. Servicio ejecuta queries paralelas:
   - Observaciones por período
   - Conductas seguras/inseguras
   - Incidentes agrupados
5. Servicio calcula KPIs y tendencias
6. Response: JSON con estadísticas
```

---

## 6. Patrones de Diseño Utilizados

1. **Inyección de Dependencias**: Configuración centralizada
2. **Repository Pattern**: Abstracción de datos
3. **Service Pattern**: Lógica de negocio centralizada
4. **DTO Pattern**: Separación de representación
5. **Middleware Pattern**: Concerns transversales
6. **Factory Pattern**: Creación de entidades
7. **Observer Pattern**: Eventos y notificaciones

---

## 7. Configuración Ambiente

### Variables de Entorno (.env)
```
# Servidor
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
ENVIRONMENT=development|production

# PocketBase
POCKETBASE_URL=http://localhost:8090
POCKETBASE_ADMIN_EMAIL=admin@example.com
POCKETBASE_ADMIN_PASSWORD=secure_password

# JWT
JWT_SECRET=your_super_secret_key_min_32_chars
JWT_EXPIRATION=900          # 15 minutos
JWT_REFRESH_EXPIRATION=604800  # 7 días

# Base de datos
DATABASE_URL=file:./pb_data/data.db  # SQLite
# DATABASE_URL=postgres://user:pass@host/db  # PostgreSQL

# Seguridad
CORS_ORIGINS=http://localhost:3000,https://app.example.com
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60

# Logs
LOG_LEVEL=info
LOG_FORMAT=json

# Email (para notificaciones)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=app@example.com
SMTP_PASSWORD=app_password
```

---

## 8. Ciclo de Despliegue

```
Desarrollo Local
    ↓
Git Push (main)
    ↓
CI/CD Pipeline (GitHub Actions)
    ↓
Tests + Lint + Coverage
    ↓
Build Docker Image
    ↓
Deploy a Staging (Render)
    ↓
E2E Tests
    ↓
Deploy a Production
    ↓
Monitoreo y Alerts
```

---

## 9. Performance y Escalabilidad

### Índices Críticos en PocketBase
- `users.email` (único)
- `observations.employee_id`
- `observations.created_at`
- `observations.area_id`
- `behaviors.observation_id`
- `audit_logs.user_id`
- `audit_logs.created_at`

### Caché Implementado
- Permisos de usuario: Caché en memoria (invalidación 1 hora)
- Roles: Caché estático
- Configuración: Caché al iniciar

### Paginación
- Default: 20 registros
- Máximo: 100 registros
- Offset-based para compatibilidad

---

## 10. Próximas Evoluciones

- [ ] Microservicios para analytics
- [ ] Message Queue (RabbitMQ) para eventos asincronos
- [ ] WebSockets para notificaciones reales
- [ ] Caché distribuida (Redis)
- [ ] API GraphQL como alternativa
- [ ] Mobile app native
