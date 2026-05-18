# Quick Start Guide - SBC Backend

## 🚀 Inicio Rápido (5 minutos)

### 1. Prerrequisitos
```bash
# Verificar Docker y Docker Compose
docker --version   # v20.10+
docker-compose --version  # v2.0+

# Clonar repo
git clone https://github.com/apariciocch/PsicoSSTCloud.git
cd PsicoSSTCloud
```

### 2. Configuración Inicial
```bash
# Copiar archivo de entorno
cp backend/.env.example backend/.env

# Generar secreto JWT seguro (32+ caracteres)
openssl rand -base64 32

# Editar backend/.env y reemplazar JWT_SECRET
nano backend/.env
```

### 3. Iniciar Sistema
```bash
# Iniciar servicios (PocketBase + Backend)
docker-compose up -d

# Verificar que todo esté corriendo
docker-compose ps

# Ver logs
docker-compose logs -f backend
```

### 4. Verificar Funcionamiento
```bash
# Health check
curl http://localhost:8080/health

# Intentar login (error esperado, es normal)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password"}'

# Respuesta:
# {"error":"USER_NOT_FOUND","message":"user not found"}
```

### 5. Crear Usuario de Prueba en PocketBase

**Acceder a Admin Panel:**
- URL: http://localhost:8090/_/
- Email: test@example.com (crear nuevo)
- Password: (cualquiera)

**Crear colección "users":**
1. Ir a "Collections"
2. Click "New collection"
3. Nombre: `users`
4. Campos:
   - `email` (text, required, unique)
   - `username` (text, required, unique)
   - `password_hash` (text, required)
   - `role_id` (text)
   - `is_active` (bool, default: true)
   - `failed_attempts` (number, default: 0)
   - `locked_until` (date)
   - `created_at` (auto)
   - `updated_at` (auto)

5. Click Save

### 6. Test Básico de API

```bash
# 1. Login (primero crear usuario en PocketBase)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "tu_email@example.com",
    "password": "tu_password"
  }'

# Respuesta:
# {
#   "access_token": "eyJhbGciOiJIUzI1NiIs...",
#   "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
#   "user": {
#     "id": "user_123",
#     "email": "tu_email@example.com",
#     "username": "usuario"
#   }
# }

# 2. Usar token en próximos requests
TOKEN="eyJhbGciOiJIUzI1NiIs..."

curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer $TOKEN"
```

---

## 📚 Endpoints Principales

### Autenticación
```bash
POST /api/v1/auth/login              # Login
POST /api/v1/auth/logout             # Logout
POST /api/v1/auth/refresh            # Renovar token
POST /api/v1/auth/change-password    # Cambiar contraseña
```

### Usuarios
```bash
GET    /api/v1/users                 # Listar usuarios
POST   /api/v1/users                 # Crear usuario
GET    /api/v1/users/:id             # Obtener usuario
PUT    /api/v1/users/:id             # Actualizar usuario
DELETE /api/v1/users/:id             # Eliminar usuario
```

### Empleados
```bash
GET    /api/v1/employees             # Listar empleados
POST   /api/v1/employees             # Crear empleado
GET    /api/v1/employees/:id         # Obtener empleado
PUT    /api/v1/employees/:id         # Actualizar empleado
```

### Observaciones
```bash
GET    /api/v1/observations          # Listar observaciones
POST   /api/v1/observations          # Crear observación
GET    /api/v1/observations/:id      # Obtener observación
PUT    /api/v1/observations/:id      # Actualizar observación
```

### Analytics
```bash
GET    /api/v1/analytics/dashboard   # Dashboard completo
GET    /api/v1/analytics/kpis        # Métricas KPI
GET    /api/v1/analytics/trends      # Tendencias
```

### Reportes
```bash
GET    /api/v1/reports               # Listar reportes
POST   /api/v1/reports               # Generar reporte
GET    /api/v1/reports/:id/download  # Descargar reporte
```

---

## 🔧 Comandos Útiles

### Desarrollo
```bash
# Iniciar servicios
docker-compose up -d

# Ver logs en tiempo real
docker-compose logs -f backend

# Detener servicios
docker-compose down

# Limpiar datos de desarrollo
docker-compose down -v

# Reconstruir imágenes
docker-compose build --no-cache
```

### Debugging
```bash
# Ejecutar Go en modo debug
dlv debug ./cmd/main.go

# Conectar debugger en VSCode (ver launch.json)

# Ver variables de entorno
docker-compose exec backend env | grep -E "JWT|POCKETBASE"
```

### Base de Datos
```bash
# Acceder a SQLite directamente
docker exec -it pocketbase sqlite3 /pb_data/data.db

# Ver tablas
.tables

# Ver schema de tabla
.schema users

# Salir
.quit
```

### Testing
```bash
# Ejecutar tests unitarios
go test -v ./...

# Con cobertura
go test -cover ./...

# Tests específicos
go test -run TestLogin -v ./internal/service/auth
```

---

## ⚙️ Configuración Comunes

### JWT Secret
```bash
# Generar secret nuevo (recomendado para cada entorno)
openssl rand -base64 32

# Copiar a .env
JWT_SECRET=tu_secreto_aqui
```

### CORS Origins
```
# En .env, agregar origins permitidos
CORS_ORIGINS=http://localhost:3000,http://localhost:8080,https://app.empresa.com
```

### Rate Limiting
```
# En .env
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW_SECONDS=60
```

### Niveles de Log
```
LOG_LEVEL=debug    # Máximo detalle (desarrollo)
LOG_LEVEL=info     # Info + warnings + errors (default)
LOG_LEVEL=warn     # Solo warnings + errors
LOG_LEVEL=error    # Solo errors (producción)
```

---

## 📋 Estructura de Directorios

```
backend/
├── cmd/
│   └── main.go                   # Punto de entrada
├── internal/
│   ├── config/
│   │   ├── config.go            # Variables de entorno
│   │   └── constants.go         # Constantes del sistema
│   ├── domain/
│   │   └── entities/
│   │       └── models.go        # Structs de modelos
│   ├── handler/                 # HTTP handlers
│   │   ├── auth.go
│   │   ├── user.go
│   │   └── ...
│   ├── service/                 # Lógica de negocio
│   │   ├── interfaces.go
│   │   ├── auth/service.go
│   │   └── ...
│   ├── repository/              # Acceso a datos
│   │   ├── pocketbase.go
│   │   └── ...
│   ├── middleware/              # Middlewares HTTP
│   │   ├── auth.go
│   │   ├── ratelimit.go
│   │   └── ...
│   ├── dto/                     # Data Transfer Objects
│   │   └── dtos.go
│   ├── pkg/                     # Paquetes reutilizables
│   │   ├── logger/
│   │   ├── response/
│   │   └── validator/
│   └── router/
│       └── routes.go            # Definición de rutas
├── .env.example                 # Variables ejemplo
├── go.mod                       # Dependencias
└── Dockerfile                   # Imagen Docker
```

---

## 🧪 Ejemplos de Requests

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "SecurePassword123!"
  }'
```

Respuesta esperada (201):
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "user_abc123",
    "email": "admin@example.com",
    "username": "admin",
    "role_id": "role_admin"
  }
}
```

### Crear Observación
```bash
curl -X POST http://localhost:8080/api/v1/observations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "employee_id": "emp_123",
    "observer_id": "user_456",
    "work_area_id": "area_789",
    "observation_date": "2026-05-18T14:30:00Z",
    "duration_minutes": 30,
    "behaviors": [
      {
        "behavior_category_id": "cat_001",
        "type": "safe",
        "description": "Utilizó EPP correctamente"
      },
      {
        "behavior_category_id": "cat_002",
        "type": "unsafe",
        "description": "No usó casco en área designada"
      }
    ]
  }'
```

### Obtener Dashboard
```bash
curl -X GET "http://localhost:8080/api/v1/analytics/dashboard?work_area_id=area_789" \
  -H "Authorization: Bearer $TOKEN"
```

---

## 🐛 Troubleshooting

### Error: "Cannot connect to PocketBase"
```bash
# Verificar que PocketBase esté corriendo
docker ps | grep pocketbase

# Reiniciar
docker restart pocketbase

# Ver logs
docker logs pocketbase
```

### Error: "Invalid JWT token"
```bash
# Verificar que el token no esté expirado
# Access token: 15 minutos
# Refresh token: 7 días

# Generar token nuevo
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "tu_refresh_token"}'
```

### Error: "Permission Denied"
```bash
# Verificar que el usuario tenga el rol correcto en PocketBase
# Admin debe tener "admin" role para crear usuarios
# Observer puede solo leer observaciones, no crear
```

### Error: "Rate limit exceeded"
```bash
# Esperar 60 segundos o cambiar RATE_LIMIT_WINDOW_SECONDS en .env
# Default: 100 requests por minuto por IP
```

---

## 📖 Recursos Adicionales

### Documentación
- [ARCHITECTURE.md](ARCHITECTURE.md) - Visión general
- [API_ENDPOINTS.md](API_ENDPOINTS.md) - Todos los endpoints
- [SECURITY.md](SECURITY.md) - Seguridad
- [DEPLOYMENT.md](DEPLOYMENT.md) - Despliegue
- [POCKETBASE_CONFIG.md](POCKETBASE_CONFIG.md) - Base de datos

### Dependencias
- Go Chi Router: https://github.com/go-chi/chi
- PocketBase: https://pocketbase.io
- JWT Go: https://github.com/golang-jwt/jwt
- Zap Logger: https://github.com/uber-go/zap

### Comunidad
- GitHub Issues: https://github.com/apariciocch/PsicoSSTCloud/issues
- Discussions: https://github.com/apariciocch/PsicoSSTCloud/discussions

---

**Versión**: 1.0  
**Última actualización**: Mayo 18, 2026  
**Mantenedor**: DevOps Team
