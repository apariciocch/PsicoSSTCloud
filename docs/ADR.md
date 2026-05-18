# Architecture Decision Record (ADR)

## ADR-001: Clean Architecture

**Status**: ACCEPTED  
**Date**: 2026-05-18  
**Context**: 
- Equipo con experiencia en arquitectura limpia
- Necesidad de escalabilidad a largo plazo
- Independencia de frameworks

**Decision**: 
Adoptar Clean Architecture con 4 capas: Handler → UseCase → Service → Repository

**Rationale**:
- Separación clara de responsabilidades
- Facilita testing unitario
- Facilita cambio de tecnologías
- Reduce acoplamiento

**Consequences**:
- Más archivos (pero mejor organizados)
- Curva de aprendizaje para nuevos developers
- Performance comparable a arquitecturas simples

---

## ADR-002: Go + Chi Router

**Status**: ACCEPTED  
**Date**: 2026-05-18  
**Alternatives Considered**:
- Python/FastAPI
- Node.js/Express
- Rust/Actix

**Decision**: 
Go 1.21 + Chi Router

**Rationale**:
- Performance superior (100k+ req/sec posible)
- Compilación a binario único, fácil despliegue
- Goroutines = concurrencia eficiente
- Error handling explícito
- Type-safe, compilado

**Consequences**:
- Team needs Go learning (1-2 semanas)
- Menos librerías de terceros que Python/Node
- Excelente para APIs y microservicios

---

## ADR-003: PocketBase + SQLite

**Status**: ACCEPTED  
**Date**: 2026-05-18  
**Alternatives Considered**:
- PostgreSQL directo
- MongoDB
- Firebase

**Decision**: 
PocketBase (SQLite en dev, PostgreSQL en prod)

**Rationale**:
- PocketBase = admin UI + auth + real-time + API gratis
- SQLite perfecto para desarrollo
- PostgreSQL para producción
- Zero-config deployment
- Migraciones automáticas

**Consequences**:
- SQLite no es multi-user/distribuido (ok para dev)
- PostgreSQL required para prod
- Vendor lock-in moderado (pero fácil migrate)

---

## ADR-004: JWT Authentication

**Status**: ACCEPTED  
**Date**: 2026-05-18  
**Alternatives Considered**:
- OAuth2/OIDC
- mTLS certificates
- API Keys

**Decision**: 
JWT con claims personalizados, refresh tokens, y JTI para revocación

**Rationale**:
- Stateless, escalable horizontalmente
- Claims pueden incluir permisos
- Refresh tokens balance seguridad/UX
- Standard industry

**Consequences**:
- Token revocation requiere fallback a lista negra
- Payload visible en base64 (no secrets en JWT)
- Implementación futura: OAuth2 layer

---

## ADR-005: RBAC (Role-Based Access Control)

**Status**: ACCEPTED  
**Date**: 2026-05-18  
**Alternatives Considered**:
- ABAC (Attribute-Based)
- Directivas granulares por campo

**Decision**: 
RBAC con 5 roles (Admin, Supervisor, Observer, Worker, Auditor) y permisos granulares

**Rationale**:
- Suficientemente flexible para SBC
- Fácil de entender y mantener
- ABAC agrega complejidad innecesaria por ahora
- Escalable a ABAC en futuro

**Consequences**:
- Rol es required field para toda user
- Permisos en JWT claims
- Migrations cuando cambien permisos

---

## ADR-006: Structured Logging (Zap)

**Status**: ACCEPTED  
**Date**: 2026-05-18  
**Alternatives Considered**:
- logrus
- log/slog (Go 1.21)
- custom logger

**Decision**: 
Uber Zap para structured logging JSON

**Rationale**:
- Performance superior (100x+ rápido que logrus)
- JSON output para parsing automatizado
- Field-based logging
- Levels: debug, info, warn, error, fatal

**Consequences**:
- Logs estructurados facilitan debugging
- Parseable por herramientas de monitoreo
- Agregación en Splunk/DataDog/CloudWatch

---

## ADR-007: Rate Limiting

**Status**: ACCEPTED  
**Date**: 2026-05-18  
**Alternatives Considered**:
- Redis-based distributed rate limiting
- Database-backed rate limiting

**Decision**: 
Token bucket algorithm per-IP, in-memory

**Rationale**:
- Suficiente para SBC (no es API pública)
- Sin dependencia adicional
- Performance excelente
- Por IP previene abuso

**Consequences**:
- No funciona con múltiples instancias (cada servidor tiene su estado)
- Para producción con múltiples servidores: implementar con Redis
- Cleanup automático de IPs inactivas

---

## ADR-008: Error Handling

**Status**: ACCEPTED  
**Date**: 2026-05-18  
**Decision**: 
Error codes estandarizados con formato: `{DOMAIN}_{NUMBER}`  
Ej: `AUTH_001`, `USER_NOT_FOUND`, `VALIDATION_ERROR`

**Rationale**:
- Clients pueden parsear y actuar en errores específicos
- Logging más fácil
- API consistente

**Error Codes**:
```
AUTH_001: Missing token
AUTH_002: Invalid format
AUTH_003: Token expired
AUTH_004: Insufficient permissions
USER_001: User not found
USER_002: Email already exists
VALIDATION_001: Invalid input
DB_001: Database error
```

---

## ADR-009: API Versioning

**Status**: ACCEPTED  
**Date**: 2026-05-18  
**Decision**: 
URL-based versioning: `/api/v1`, `/api/v2`

**Rationale**:
- Explícito en URL
- Fácil para clientes
- Permite deprecación gradual
- Industry standard

**Consequences**:
- Cambios breaking → nueva versión
- Mantener múltiples versiones temporalmente
- Documentación por versión

---

## ADR-010: Deployment Strategy

**Status**: ACCEPTED  
**Date**: 2026-05-18  
**Alternatives Considered**:
- Blue-Green
- Canary
- Rolling

**Decision**: 
Blue-Green para inicial, Canary en futuro

**Rationale**:
- Blue-Green: zero downtime, fácil rollback
- Canary: para testing en prod (fase 2)

**Process**:
1. Deploy v2 en verde
2. Test v2
3. Switch router a v2
4. Mantener v1 como fallback

---

## ADR-011: Configuration Management

**Status**: ACCEPTED  
**Date**: 2026-05-18  
**Decision**: 
Environment variables + .env file, structured en config.go

**Rationale**:
- 12-factor app compliance
- Diferentes configs por entorno
- Secrets no en código

**Sensitive Variables**:
- JWT_SECRET (min 32 chars)
- DB passwords
- API keys externos

---

## ADR-012: Database Migrations

**Status**: ACCEPTED  
**Date**: 2026-05-18  
**Decision**: 
PocketBase handles migrations automáticamente  
Manual DDL en SQL cuando sea necesario

**Rationale**:
- PocketBase tiene UI para schema changes
- Versioned migrations via hooks
- SQL scripts para índices complejos

---

## ADR-013: Testing Strategy

**Status**: ACCEPTED  
**Date**: 2026-05-18  
**Levels**:
- Unit tests (>70% cobertura) - servicios
- Integration tests - repositorios
- E2E tests - flows completos

**Tools**:
- Testing: Go testing + testify
- Integration: docker-compose test
- E2E: Thunder Client / Postman

---

## ADR-014: Documentation

**Status**: ACCEPTED  
**Date**: 2026-05-18  
**Types**:
1. Architecture (ARCHITECTURE.md)
2. API (API_ENDPOINTS.md, OpenAPI)
3. Database (COLLECTIONS.md, DATABASE_DESIGN.md)
4. Security (SECURITY.md)
5. Deployment (DEPLOYMENT.md)
6. Code (inline comments + doc strings)

---

## ADR-015: Future Considerations

**Status**: PROPOSED  
**Topics for Phase 2+**:

### GraphQL
- Status: Bajo interés actual
- Timeline: Phase 3+
- Rationale: REST es suficiente, GraphQL agrega complejidad

### WebSockets
- Status: Requerido para notificaciones en tiempo real
- Timeline: Phase 2
- Rationale: Observaciones/notificaciones real-time

### Microservicios
- Status: No inmediato
- Timeline: Phase 3+ cuando crecimiento justifique
- Split: Analytics → servicio separado

### Machine Learning
- Status: Fase exploratoria
- Timeline: Phase 4
- Use case: Predicción de riesgos

### API Gateway
- Status: Requerido para múltiples backends
- Timeline: Phase 2 si escalamos
- Options: Kong, Traefik, AWS API Gateway

### Event Streaming
- Status: Considerar para auditoría
- Timeline: Phase 2
- Options: Kafka, RabbitMQ, Redis Streams

---

## ADR-016: Code Quality Standards

**Status**: ACCEPTED  
**Date**: 2026-05-18  
**Linting**: golangci-lint
**Formatting**: gofmt
**Testing**: >70% coverage
**Security**: gosec
**Documentation**: godoc comments

**Pre-commit Hooks**:
```bash
#!/bin/bash
go fmt ./...
golangci-lint run ./...
go test -v ./...
```

---

## ADR-017: Performance Targets

**Status**: ACCEPTED  
**Targets**:
- API response: <100ms (p95)
- Database query: <50ms (p95)
- Throughput: 1000 req/sec
- Availability: 99.9%

**Monitoring**:
- Prometheus metrics
- Response time histograms
- Error rates
- Database query times

---

## ADR-018: Security Standards

**Status**: ACCEPTED  
**Date**: 2026-05-18  

**Encryption**:
- In transit: TLS 1.2+ (HTTPS)
- At rest: AES-256 para datos sensibles
- Passwords: bcrypt cost=12

**Authentication**:
- JWT HS256
- Token rotation
- Session invalidation

**Authorization**:
- RBAC granular
- Audit logging
- Principle of least privilege

**Network**:
- Firewall rules
- Rate limiting
- DDoS protection (Cloudflare)

---

## Review Process for ADRs

New ADRs require:
1. ✅ Senior architect review
2. ✅ Team discussion (async or sync)
3. ✅ Clear alternatives considered
4. ✅ Rationale documentado
5. ✅ Consequences identificadas
6. ✅ Implementation plan

**ADR Template**:
```
# ADR-NNN: [Title]

**Status**: PROPOSED | ACCEPTED | SUPERSEDED | DEPRECATED
**Date**: YYYY-MM-DD
**Author**: Name

## Context
[Situación actual y problema]

## Decision
[Decisión tomada]

## Rationale
[Por qué esta decisión]

## Alternatives Considered
- Option A: [pros/cons]
- Option B: [pros/cons]

## Consequences
- Positivas: [...]
- Negativas: [...]

## References
- [Links relevantes]
```

---

**Document Status**: Complete  
**Version**: 1.0  
**Last Updated**: May 18, 2026  
**Maintained By**: Architecture Team
