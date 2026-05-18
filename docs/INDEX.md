# Índice Completo de Documentación - SBC Backend

## 📑 Guía de Navegación

Esta documentación cubre todos los aspectos del sistema SBC (Seguridad Basada en el Comportamiento) desarrollado con Go, PocketBase y Clean Architecture.

---

## 🎯 Punto de Inicio Recomendado

### Para Diferentes Roles:

**👨‍💼 Product Manager / Stakeholder**
1. [README de proyecto](./README.md) - Visión general
2. [ARCHITECTURE.md](./ARCHITECTURE.md) - Qué construimos
3. [ROADMAP.md](./ROADMAP.md) - Timeline y fases
4. [DEPLOYMENT.md](./DEPLOYMENT.md) - Despliegue

**👨‍💻 Backend Developer**
1. [QUICKSTART.md](./QUICKSTART.md) - Comenzar en 5 minutos
2. [ARCHITECTURE.md](./ARCHITECTURE.md) - Visión general
3. [API_ENDPOINTS.md](./API_ENDPOINTS.md) - Endpoints a implementar
4. [COLLECTIONS.md](./COLLECTIONS.md) - Base de datos
5. [SECURITY.md](./SECURITY.md) - Seguridad
6. [ADR.md](./ADR.md) - Decisiones de arquitectura

**🔧 DevOps / SRE**
1. [DEPLOYMENT.md](./DEPLOYMENT.md) - Despliegue
2. [DOCKER setup](../docker-compose.yml) - Orquestación local
3. [SECURITY.md](./SECURITY.md) - Seguridad en producción
4. [POCKETBASE_CONFIG.md](./POCKETBASE_CONFIG.md) - Base de datos

**🏗️ Software Architect**
1. [ARCHITECTURE.md](./ARCHITECTURE.md) - Diseño completo
2. [ADR.md](./ADR.md) - Decisiones tomadas
3. [DATABASE_DESIGN.md](./DATABASE_DESIGN.md) - Esquema relacional
4. [SECURITY.md](./SECURITY.md) - Seguridad
5. [ROADMAP.md](./ROADMAP.md) - Visión futura

---

## 📚 Documentación Completa

### 1. **QUICKSTART.md** ⚡
**Propósito**: Comenzar rápidamente en 5 minutos  
**Contenido**:
- Instalación de prerrequisitos
- Configuración inicial
- Primeros tests de API
- Troubleshooting básico
- Comandos útiles

**Cuándo leer**: Primera vez en el proyecto

---

### 2. **ARCHITECTURE.md** 🏗️
**Propósito**: Comprensión profunda del diseño del sistema  
**Contenido**:
- Visión y objetivos del proyecto
- Stack tecnológico detallado
- 4 capas de Clean Architecture
- Seguridad y autenticación
- Performance y escalabilidad
- Componentes principales
- Flujos de datos
- Diagramas conceptuales

**Cuándo leer**: Para entender diseño del sistema

---

### 3. **API_ENDPOINTS.md** 🔌
**Propósito**: Referencia completa de todos los endpoints REST  
**Contenido**:
- 35+ endpoints organizados por grupo
- Request/response ejemplos
- Códigos de error específicos
- Parámetros de query
- Headers requeridos
- Autenticación por endpoint
- Rate limits
- Métodos HTTP

**Endpoints Cubiertos**:
- /auth (login, logout, refresh)
- /users (CRUD)
- /employees (CRUD)
- /observations (CRUD + filtros)
- /incidents (CRUD)
- /corrective-actions (management)
- /analytics (dashboard, KPIs)
- /reports (generación, export)
- /audit-logs (búsqueda)
- /notifications (alertas)

**Cuándo leer**: Para integración frontend o testing

---

### 4. **COLLECTIONS.md** 📊
**Propósito**: Especificación completa de las 14 colecciones PocketBase  
**Contenido**:
- Schema de cada colección
- Tipos de campos
- Validaciones
- Relaciones
- Índices
- Ejemplos de datos
- Reglas de acceso RBAC

**Colecciones Cubiertas**:
1. users
2. roles
3. work_areas
4. employees
5. observations
6. behaviors
7. behavior_categories
8. incidents
9. corrective_actions
10. attachments
11. notifications
12. audit_logs
13. kpis
14. reports

**Cuándo leer**: Para entender modelo de datos

---

### 5. **DATABASE_DESIGN.md** 🗄️
**Propósito**: Esquema relacional y DDL SQL  
**Contenido**:
- Normalización (3NF)
- DDL para todas las tablas
- Índices de performance
- Foreign keys y constraints
- Data types y validaciones
- Integridad referencial
- Views SQL útiles

**Cuándo leer**: Para DBA o SQL engineers

---

### 6. **SECURITY.md** 🔐
**Propósito**: Estrategia completa de seguridad  
**Contenido**:
- Autenticación JWT
- Encryption (AES-256, bcrypt)
- Rate limiting
- Input validation
- CORS y headers de seguridad
- Auditoría completa
- Incident response
- Compliance (OWASP, ISO 27001)
- Penetration testing
- Secret management

**Cuándo leer**: Para security reviews o compliance

---

### 7. **POCKETBASE_CONFIG.md** 🔧
**Propósito**: Configuración y hooks de PocketBase  
**Contenido**:
- Scripts de inicialización
- Hooks en JavaScript
- Auto-generación de IDs
- Validaciones customizadas
- Auditoría automática
- KPI calculations
- Backups automáticos
- Reglas de acceso
- Índices de performance

**Cuándo leer**: Para personalizar PocketBase

---

### 8. **DEPLOYMENT.md** 🚀
**Propósito**: Guía completa de despliegue  
**Contenido**:
- Prerrequisitos
- Despliegue local (docker-compose)
- Despliegue en servidores
- Configuración Nginx
- SSL/TLS certificates
- Alternativas: Render, Kubernetes
- Monitoreo y logs
- Backups y disaster recovery
- Scaling horizontal y vertical
- Troubleshooting
- Checklist de despliegue
- Security hardening

**Opciones de Despliegue**:
- Docker Compose (desarrollo)
- VPS + Docker (pequeño/mediano)
- Render PaaS (manejado)
- Kubernetes (grande)

**Cuándo leer**: Para desplegar a producción

---

### 9. **ROADMAP.md** 🗺️
**Propósito**: Planes de desarrollo y próximos pasos  
**Contenido**:
- Resumen ejecutivo
- Archivos generados
- 4 fases de desarrollo (6 semanas)
- Priorización de tareas
- Estimación de recursos
- Costos de infraestructura
- Performance optimization
- Monitoring & observability
- Compliance standards
- Testing strategy
- Recomendaciones finales

**Fases**:
1. Base (Semana 1-2): Repositorios, servicios, handlers
2. Features (Semana 3-4): Analytics, reportes, notificaciones
3. Producción (Semana 5-6): Testing, optimización, deploy
4. Mejoras (Continuo): GraphQL, WebSockets, ML

**Cuándo leer**: Para planning y prioritización

---

### 10. **ADR.md** 📋
**Propósito**: Architecture Decision Records - todas las decisiones importantes  
**Contenido**:
- 18 ADRs documentadas
- Clean Architecture
- Go + Chi Router
- PocketBase
- JWT Authentication
- RBAC
- Structured logging
- Rate limiting
- Error handling
- API versioning
- Deployment strategy
- Configuration management
- Database migrations
- Testing strategy
- Documentation
- Code quality
- Performance targets
- Security standards

**Cuándo leer**: Para entender por qué cada decisión

---

### 11. **ADR.md - Futuro** 🔮
**Propósito**: Consideraciones para fases futuras  
**Tecnologías Consideradas**:
- GraphQL (Phase 3+)
- WebSockets (Phase 2)
- Microservicios (Phase 3+)
- Machine Learning (Phase 4)
- API Gateway (Phase 2)
- Event Streaming (Phase 2)

**Cuándo leer**: Para planning a largo plazo

---

## 🗂️ Estructura de Carpetas del Proyecto

```
PsicoSSTCloud/
├── docs/                           # 📚 Documentación
│   ├── README.md                   # Punto de inicio
│   ├── QUICKSTART.md               # ⚡ 5 minutos
│   ├── ARCHITECTURE.md             # 🏗️ Diseño
│   ├── API_ENDPOINTS.md            # 🔌 APIs
│   ├── COLLECTIONS.md              # 📊 BD
│   ├── DATABASE_DESIGN.md          # 🗄️ SQL
│   ├── SECURITY.md                 # 🔐 Seguridad
│   ├── POCKETBASE_CONFIG.md        # 🔧 PocketBase
│   ├── DEPLOYMENT.md               # 🚀 Deploy
│   ├── ROADMAP.md                  # 🗺️ Fases
│   ├── ADR.md                      # 📋 Decisiones
│   └── INDEX.md                    # 📑 Este archivo
│
├── backend/                        # Go backend
│   ├── cmd/main.go                 # Punto de entrada
│   ├── internal/
│   │   ├── config/
│   │   │   ├── config.go           # Configuración
│   │   │   └── constants.go        # Constantes
│   │   ├── domain/entities/
│   │   │   └── models.go           # Modelos
│   │   ├── handler/                # HTTP handlers
│   │   │   ├── auth.go
│   │   │   ├── user.go
│   │   │   └── ...
│   │   ├── service/                # Lógica
│   │   │   ├── interfaces.go       # Contratos
│   │   │   ├── auth/service.go     # Auth impl
│   │   │   └── ...
│   │   ├── repository/             # Acceso datos
│   │   │   ├── pocketbase.go
│   │   │   └── ...
│   │   ├── middleware/             # Middlewares
│   │   │   ├── auth.go
│   │   │   ├── ratelimit.go
│   │   │   ├── logger.go
│   │   │   └── cors.go
│   │   ├── dto/dtos.go             # DTOs
│   │   ├── pkg/                    # Utilidades
│   │   │   ├── logger/
│   │   │   ├── response/
│   │   │   └── validator/
│   │   └── router/routes.go        # Rutas
│   ├── .env.example                # Variables ejemplo
│   ├── go.mod                      # Dependencias
│   └── Dockerfile                  # Imagen Docker
│
├── pocketbase/                     # Configuración
├── docker-compose.yml              # Orquestación
└── .gitignore
```

---

## 🔍 Mapeo de Responsabilidades

### Por Componente:

| Componente | Archivo | Responsable |
|-----------|---------|------------|
| Autenticación | `SECURITY.md`, `internal/service/auth` | Backend Dev |
| API Design | `API_ENDPOINTS.md`, `internal/handler` | Architect |
| Base de Datos | `DATABASE_DESIGN.md`, `internal/repository` | DBA / Backend |
| Despliegue | `DEPLOYMENT.md`, `docker-compose.yml` | DevOps |
| Seguridad | `SECURITY.md`, `ADR.md` | Architect / Security |
| Performance | `ROADMAP.md`, `ADR.md` | Architect / DevOps |

---

## 📞 Referencias Rápidas

### Comandos Frecuentes

```bash
# Iniciar
docker-compose up -d

# Logs
docker-compose logs -f backend

# Testing
go test -v ./...

# Linting
golangci-lint run ./...

# Build
docker build -t psicosst-backend .
```

### Variables Críticas

- `JWT_SECRET`: Min 32 caracteres
- `POCKETBASE_ADMIN_PASSWORD`: Segura
- `CORS_ORIGINS`: Whitelist de dominios
- `ENVIRONMENT`: development/production

### Error Codes

- `AUTH_001`: Missing token
- `AUTH_004`: Insufficient permissions
- `USER_NOT_FOUND`: Usuario inexistente
- `VALIDATION_ERROR`: Input inválido
- `DB_001`: Database error

### Health Checks

```bash
curl http://localhost:8080/health      # Backend
curl http://localhost:8090/api/health  # PocketBase
```

---

## 🎓 Rutas de Aprendizaje

### Para Principiante en el Proyecto
1. QUICKSTART.md (15 min)
2. ARCHITECTURE.md (30 min)
3. API_ENDPOINTS.md (20 min)
4. Empezar implementación

### Para Nuevo Backend Developer
1. QUICKSTART.md (15 min)
2. ARCHITECTURE.md (30 min)
3. COLLECTIONS.md (20 min)
4. ADR.md secciones 1-6 (20 min)
5. CODE (1-2 horas)

### Para DevOps Engineer
1. DEPLOYMENT.md (1 hora)
2. docker-compose.yml (15 min)
3. SECURITY.md (20 min)
4. Establecer infraestructura

### Para Security Auditor
1. SECURITY.md (1 hora)
2. ADR.md secciones 18 (20 min)
3. API_ENDPOINTS.md (30 min)
4. Code review

---

## ✅ Checklist de Onboarding

- [ ] Clonar repositorio
- [ ] Ejecutar docker-compose up -d
- [ ] Leer QUICKSTART.md
- [ ] Test primeros endpoints
- [ ] Leer ARCHITECTURE.md
- [ ] Revisar código en backend/internal
- [ ] Entender Clean Architecture
- [ ] Leer API_ENDPOINTS.md
- [ ] Revisar tu área de responsabilidad
- [ ] Preguntar dudas

---

## 📊 Estadísticas del Proyecto

| Métrica | Valor |
|---------|-------|
| Líneas de código Go | ~3,000 |
| Líneas de documentación | ~5,000 |
| Colecciones PocketBase | 14 |
| Endpoints API | 35+ |
| Tipos de errores | 10+ |
| Roles de acceso | 5 |
| Permisos granulares | 20+ |
| Pruebas unitarias | ~50 |
| Coverage mínimo | 70% |

---

## 🤝 Contribución

Para contribuir:
1. Branching: feature/XXX, bugfix/XXX
2. Commits: Conventional Commits
3. PRs: Incluir descripción y tests
4. Documentación: Actualizar ADRs si cambia arquitectura

---

## 📞 Contacto y Soporte

**Problemas técnicos**:
- GitHub Issues: https://github.com/apariciocch/PsicoSSTCloud/issues

**Documentación**:
- Todos los archivos en `/docs`
- Código documentado en línea

**Arquitecto**:
- Review arquitectónico
- Pair programming
- Mentoring

---

## 📝 Changelog de Documentación

| Versión | Fecha | Cambios |
|---------|-------|---------|
| 1.0 | 2026-05-18 | Documentación inicial completa |

---

## 🎯 Próximas Actualizaciones

- [ ] OpenAPI/Swagger spec
- [ ] Postman collection
- [ ] Runbooks operacionales
- [ ] Training videos
- [ ] Case studies de implementación

---

**Documento**: INDEX.md  
**Versión**: 1.0  
**Última actualización**: Mayo 18, 2026  
**Mantenido por**: Equipo de Desarrollo
