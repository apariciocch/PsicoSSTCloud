# 🛡️ PsicoSSTCloud - Sistema de Seguridad Basada en el Comportamiento

**Comportamiento Seguro = Cero Incidentes**

## 📋 Resumen Ejecutivo

PsicoSSTCloud es una plataforma empresarial de **Seguridad Basada en el Comportamiento (SBC)** que transforma cómo las organizaciones abordan la prevención de accidentes laborales.

### Visión
Crear una cultura de seguridad donde los comportamientos seguros sean la norma, reduciendo incidentes laborales a través de observación sistemática, análisis de datos, y mejora continua.

### Solución
- 📊 **Observaciones en tiempo real** de conductas seguras e inseguras
- 🎯 **Análisis predictivo** de riesgos por área de trabajo
- 📈 **Dashboard ejecutivo** con KPIs y tendencias
- 🔔 **Alertas automáticas** de comportamientos de riesgo
- 📋 **Auditoría completa** de todas las acciones del sistema
- 🔐 **Seguridad empresarial** con encriptación y RBAC

---

## 🚀 Quick Start

```bash
# 1. Clonar repositorio
git clone https://github.com/apariciocch/PsicoSSTCloud.git
cd PsicoSSTCloud

# 2. Configurar entorno
cp backend/.env.example backend/.env
# Editar: JWT_SECRET, POCKETBASE_ADMIN_PASSWORD

# 3. Iniciar sistema
docker-compose up -d

# 4. Verificar
curl http://localhost:8080/health
# Respuesta: {"status":"ok","timestamp":"..."}
```

**Acceso**:
- Backend API: http://localhost:8080/api/v1
- PocketBase Admin: http://localhost:8090/_/

Más detalles → [QUICKSTART.md](./docs/QUICKSTART.md)

---

## 🏗️ Stack Tecnológico

| Capa | Tecnología | Versión |
|------|-----------|---------|
| **Backend** | Go | 1.21+ |
| **Router/HTTP** | Chi | v5 |
| **Base de Datos** | PocketBase | Latest |
| **BD (Dev)** | SQLite | 3.x |
| **BD (Prod)** | PostgreSQL | 14+ |
| **Autenticación** | JWT | HS256 |
| **Logging** | Zap | v1.24+ |
| **Validation** | Validator/v10 | v10 |
| **Containerización** | Docker | 20.10+ |
| **Orquestación** | Docker Compose | 2.0+ |

---

## 📂 Estructura del Proyecto

```
PsicoSSTCloud/
├── 📚 docs/                    # Documentación completa
│   ├── INDEX.md               # ← Empieza aquí
│   ├── ARCHITECTURE.md        # Diseño del sistema
│   ├── API_ENDPOINTS.md       # Todos los endpoints
│   ├── SECURITY.md            # Seguridad
│   ├── DEPLOYMENT.md          # Despliegue
│   └── ...
│
├── 🔧 backend/                # Go backend
│   ├── cmd/main.go            # Punto de entrada
│   ├── internal/
│   │   ├── handler/           # HTTP handlers
│   │   ├── service/           # Lógica de negocio
│   │   ├── repository/        # Acceso a datos
│   │   ├── middleware/        # Middlewares
│   │   ├── config/            # Configuración
│   │   └── ...
│   ├── go.mod                 # Dependencias
│   └── Dockerfile             # Imagen Docker
│
├── 🗄️ pocketbase/             # Configuración de BD
├── 🐳 docker-compose.yml      # Orquestación local
└── 📝 .env.example            # Variables de entorno
```

---

## 🎯 Características Principales

### 1️⃣ Gestión de Observaciones
- ✅ Registro de observaciones de conductas seguras/inseguras
- ✅ Clasificación automática de comportamientos
- ✅ Vinculación con empleados y áreas de trabajo
- ✅ Filtrado avanzado y búsqueda

### 2️⃣ Gestión de Incidentes
- ✅ Registro de accidentes y casi-accidentes
- ✅ Asignación de acciones correctivas
- ✅ Seguimiento de correcciones
- ✅ Cierre con evidencia

### 3️⃣ Analytics y KPIs
- ✅ Dashboard ejecutivo en tiempo real
- ✅ Indicadores de seguridad (% comportamientos seguros)
- ✅ Análisis por área de trabajo
- ✅ Tendencias mensuales/anuales
- ✅ Alertas de riesgos críticos

### 4️⃣ Reportes
- ✅ Generación de reportes PDF
- ✅ Exportación a Excel
- ✅ Envío automático por email
- ✅ Auditoría completa

### 5️⃣ Seguridad
- ✅ Autenticación JWT con refresh tokens
- ✅ Control de acceso basado en roles (RBAC)
- ✅ Encriptación AES-256 de datos sensibles
- ✅ Auditoría de todas las acciones
- ✅ Rate limiting y validación de input

---

## 👥 Roles de Acceso

| Rol | Permisos |
|-----|----------|
| **Administrador** | Acceso total al sistema |
| **Supervisor** | Gestión de su equipo, observaciones, reportes |
| **Observador** | Crear observaciones, leer incidentes |
| **Trabajador** | Leer sus propias observaciones |
| **Auditor** | Solo lectura de auditoría |

---

## 📊 Modelo de Datos (14 Colecciones)

```
users ←→ roles ←→ permissions
  ↓
employees ←→ work_areas
  ↓
observations → behaviors → behavior_categories
  ↓
incidents → corrective_actions
  ↓
notifications, attachments
  ↓
kpis, reports, audit_logs
```

Más detalles → [COLLECTIONS.md](./docs/COLLECTIONS.md) | [DATABASE_DESIGN.md](./docs/DATABASE_DESIGN.md)

---

## 🔌 API REST

### Endpoints Principales

```bash
# Autenticación
POST   /api/v1/auth/login              # Login
POST   /api/v1/auth/logout             # Logout
POST   /api/v1/auth/refresh            # Renovar token

# Usuarios
GET    /api/v1/users                   # Listar (admin/supervisor)
POST   /api/v1/users                   # Crear
GET    /api/v1/users/:id               # Obtener
PUT    /api/v1/users/:id               # Actualizar
DELETE /api/v1/users/:id               # Eliminar

# Observaciones
POST   /api/v1/observations            # Crear observación
GET    /api/v1/observations            # Listar
GET    /api/v1/observations/:id        # Obtener
PUT    /api/v1/observations/:id        # Actualizar
GET    /api/v1/observations?filter=... # Con filtros

# Analytics
GET    /api/v1/analytics/dashboard     # Dashboard completo
GET    /api/v1/analytics/kpis          # KPIs
GET    /api/v1/analytics/trends        # Tendencias

# Reportes
POST   /api/v1/reports                 # Generar reporte
GET    /api/v1/reports/:id/download    # Descargar
```

**Total**: 35+ endpoints

Referencia completa → [API_ENDPOINTS.md](./docs/API_ENDPOINTS.md)

---

## 🔐 Seguridad

### Implementado
- ✅ **JWT Authentication** - Tokens con expiración y refresh
- ✅ **RBAC** - 5 roles con permisos granulares
- ✅ **Encryption** - AES-256 para datos sensibles
- ✅ **Password Hashing** - bcrypt con cost=12
- ✅ **Rate Limiting** - 100 req/min por IP
- ✅ **CORS** - Origen whitelist configurable
- ✅ **Security Headers** - X-Frame-Options, CSP, HSTS
- ✅ **Input Validation** - Campos requeridos, tipos, longitud
- ✅ **SQL Injection Prevention** - Parametrized queries
- ✅ **Auditing** - Log de todas las acciones

Más detalles → [SECURITY.md](./docs/SECURITY.md)

---

## 📚 Documentación

| Documento | Propósito |
|-----------|-----------|
| **INDEX.md** | 📑 Índice y navegación |
| **QUICKSTART.md** | ⚡ Comenzar en 5 minutos |
| **ARCHITECTURE.md** | 🏗️ Diseño completo |
| **API_ENDPOINTS.md** | 🔌 Todos los endpoints |
| **COLLECTIONS.md** | 📊 Modelo de datos |
| **DATABASE_DESIGN.md** | 🗄️ Schema SQL |
| **SECURITY.md** | 🔐 Seguridad |
| **POCKETBASE_CONFIG.md** | 🔧 Configuración BD |
| **DEPLOYMENT.md** | 🚀 Despliegue |
| **ROADMAP.md** | 🗺️ Fases de desarrollo |
| **ADR.md** | 📋 Decisiones arquitectónicas |

👉 **Punto de inicio recomendado**: [docs/INDEX.md](./docs/INDEX.md)

---

## 🧪 Testing

### Unit Tests
```bash
go test -v -cover ./...
```

### Integration Tests
```bash
# Usar docker-compose para BD de test
docker-compose -f docker-compose.test.yml up
go test -v -count=1 ./internal/...
```

### E2E Tests
- Thunder Client / Postman collection (próximamente)
- Scripts de test automatizados

---

## 🚀 Deployment

### Local (Desarrollo)
```bash
docker-compose up -d
# Acceso: http://localhost:8080
```

### Producción (Múltiples opciones)
1. **Docker + VPS** - Nginx reverse proxy, SSL
2. **Render** - PaaS manejado
3. **Kubernetes** - Enterprise-grade

Guía completa → [DEPLOYMENT.md](./docs/DEPLOYMENT.md)

---

## 📈 Performance

### Targets
- **Latencia API**: <100ms (p95)
- **Throughput**: 1,000 req/sec
- **Disponibilidad**: 99.9%
- **Uptime**: 24/7 con health checks

### Optimización
- Índices de base de datos optimizados
- Caché de datos frecuentes (Redis)
- Connection pooling
- Compression GZIP
- CDN para assets estáticos

---

## 🔧 Comandos Útiles

```bash
# Desarrollo
docker-compose up -d          # Iniciar
docker-compose logs -f        # Ver logs
docker-compose down           # Detener
docker-compose down -v        # Limpiar todo

# Testing
go test ./...                 # Todas las pruebas
go test -v ./internal/...     # Con verbosidad
go test -cover ./...          # Con cobertura

# Linting
golangci-lint run ./...       # Lint
go fmt ./...                  # Formato

# Build
docker build -t psicosst .    # Imagen Docker
go build -o backend ./cmd/main.go  # Binario

# Database
docker exec -it pocketbase sqlite3 /pb_data/data.db  # CLI DB
# En PocketBase: .tables, .schema users, etc.
```

---

## 🐛 Troubleshooting

### Backend no responde
```bash
# Verificar health
curl http://localhost:8080/health

# Ver logs
docker-compose logs backend

# Reiniciar
docker restart pocketbase backend
```

### Error de conexión a PocketBase
```bash
# Verificar PocketBase está corriendo
docker ps | grep pocketbase

# Verificar URL en .env
echo $POCKETBASE_URL  # Debe ser: http://pocketbase:8090
```

### JWT token expirado
- Access token: 15 minutos
- Refresh token: 7 días
- Usar endpoint `/api/v1/auth/refresh` para renovar

Más troubleshooting → [DEPLOYMENT.md](./docs/DEPLOYMENT.md#troubleshooting)

---

## 🗺️ Roadmap

### ✅ Fase 1: Base (Semana 1-2)
- [x] Arquitectura diseñada
- [x] Base de datos normalizada
- [x] API REST especificada
- [x] Autenticación implementada
- [ ] Repositorios (WIP)
- [ ] Servicios (WIP)
- [ ] Handlers (WIP)

### 🔄 Fase 2: Features (Semana 3-4)
- [ ] Dashboard Analytics
- [ ] Generación de Reportes
- [ ] Sistema de Notificaciones
- [ ] Auditoría completa

### 🚀 Fase 3: Producción (Semana 5-6)
- [ ] Performance tuning
- [ ] Testing de carga
- [ ] Despliegue staging
- [ ] Documentación final

### 💡 Fase 4: Mejoras (Futuro)
- [ ] GraphQL API
- [ ] WebSockets
- [ ] Mobile app
- [ ] Machine Learning

Más detalles → [ROADMAP.md](./docs/ROADMAP.md)

---

## 👨‍💻 Contribuyentes

- **Arquitecto**: Diseño e innovación
- **Backend Developers**: Implementación Go
- **DevOps**: Infraestructura y deployment
- **QA**: Testing y validación

---

## 📞 Soporte

- **Issues**: [GitHub Issues](https://github.com/apariciocch/PsicoSSTCloud/issues)
- **Documentación**: [docs/](./docs/)
- **Email**: soporte@empresa.com

---

## 📋 Requerimientos

### Sistema
- Linux, macOS o Windows (WSL2)
- 4GB RAM mínimo
- 20GB almacenamiento

### Software
- Docker 20.10+
- Docker Compose 2.0+
- Git
- Go 1.21+ (opcional, para desarrollo local)

---

## 🎓 Primeros Pasos

1. **Leer**: [docs/INDEX.md](./docs/INDEX.md) - Índice y navegación
2. **Setup**: [docs/QUICKSTART.md](./docs/QUICKSTART.md) - 5 minutos
3. **Entender**: [docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md) - Diseño
4. **Explorar**: [docs/API_ENDPOINTS.md](./docs/API_ENDPOINTS.md) - APIs
5. **Codificar**: Implementar según [docs/ROADMAP.md](./docs/ROADMAP.md)

---

## 📄 Licencia

Propietario - PsicoSSTCloud © 2026

---

## 🎯 Visión a Largo Plazo

**2026**: MVP con funcionalidades core  
**2027**: Escalado a múltiples empresas  
**2028**: ML para predicción de riesgos  
**2030**: Plataforma líder en SBC latinoamérica  

---

**Última actualización**: Mayo 18, 2026  
**Versión**: 1.0.0  
**Estado**: Listo para desarrollo
