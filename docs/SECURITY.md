# Estrategia de Seguridad - SBC Backend

## 1. AUTENTICACIÓN Y AUTORIZACIÓN

### 1.1 JWT (JSON Web Tokens)

**Estructura del Token:**
```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT"
  },
  "payload": {
    "sub": "user_id_123",
    "email": "usuario@empresa.com",
    "role_id": "role_observer",
    "permissions": ["observations.view", "observations.create"],
    "is_super_admin": false,
    "iat": 1716033000,
    "exp": 1716033900,
    "jti": "unique_token_id"
  }
}
```

**Tiempos de Expiración:**
- Access Token: 15 minutos
- Refresh Token: 7 días
- Remember Me Token: 30 días

**Seguridad JWT:**
- Secret mínimo 32 caracteres
- Rotar secret cada 90 días
- Token revocation list para logout inmediato
- Verificar JTI para evitar replay attacks

### 1.2 Autenticación Multifactor (MFA)

**Implementación Optional:**
- TOTP (Time-based One-Time Password) con Google Authenticator
- SMS OTP para casos críticos
- Backup codes para recuperación

### 1.3 Control de Acceso Basado en Roles (RBAC)

**Estructura Permisos:**
```
resource.action
ejemplo: observations.create, users.delete
```

**Verificación en Middleware:**
```
GET /api/observations
├─ Verificar token válido
├─ Extraer user_id y role_id
├─ Validar permiso "observations.view"
├─ Validar scope (datos propios, área, global)
└─ Ejecutar handler
```

**Granularidad:**
- A nivel de endpoint
- A nivel de atributo (campo del recurso)
- A nivel de registro (scope)

---

## 2. ENCRIPTACIÓN DE DATOS

### 2.1 Contraseñas

```go
// Hashing con bcrypt
password := "UserPassword123!"
hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// Cost: 12 (recomendado para producción)
```

**Requisitos Contraseña:**
- Mínimo 8 caracteres
- Debe incluir: mayúscula, minúscula, número, símbolo
- No usar información personal
- Historia: últimas 5 contraseñas
- Expiración: 90 días
- Reset fuerza: password temporal + email

### 2.2 Datos Sensibles en Tránsito

**HTTPS/TLS 1.3**
- Certificado SSL válido
- HSTS header habilitado
- Cipher suites modernos
- Certificate pinning para apps móviles

### 2.3 Datos Sensibles en Reposo

**AES-256-GCM para campos críticos:**
- Phone numbers
- Document numbers
- Personal emails
- MFA secrets

```go
// Encriptación AES-256-GCM
cipher, err := aes.NewCipher(key)
gcm, err := cipher.NewGCM()
nonce := make([]byte, gcm.NonceSize())
ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
```

---

## 3. VALIDACIÓN Y SANITIZACIÓN

### 3.1 Validación de Entrada

**DTO con tags de validación:**
```go
type LoginRequest struct {
  Email    string `json:"email" validate:"required,email,max=255"`
  Password string `json:"password" validate:"required,min=8,max=128"`
}
```

**Validadores Custom:**
- Documento: formato según país
- Teléfono: formato internacional
- URL: validar esquema y dominio
- JSON: estructura esperada

### 3.2 Sanitización

**Contra SQL Injection:**
- Usar prepared statements en todas las queries
- Nunca concatenar strings en SQL
- Validar tipos de datos esperados

**Contra XSS:**
- Escapar caracteres especiales en respuestas JSON
- Content-Type: application/json; charset=utf-8
- X-Content-Type-Options: nosniff

**Contra CSRF:**
- CSRF tokens para operaciones POST/PUT/DELETE
- SameSite=Strict en cookies
- Validar Origin header

### 3.3 Rate Limiting

**Por Endpoint:**
```
POST /api/auth/login
├─ Límite: 5 intentos / 15 minutos
├─ Por: IP address
└─ Bloqueo: 30 minutos después de 5 intentos

GET /api/observations
├─ Límite: 100 requests / minuto
├─ Por: User ID
└─ Respuesta: 429 Too Many Requests
```

**Implementación:**
- Token bucket algorithm
- Redis para distribuida
- In-memory para desarrollo

---

## 4. AUDITORÍA Y LOGGING

### 4.1 Audit Log

**Eventos Auditados:**
- Login / Logout
- Create / Update / Delete en cualquier tabla
- Access denegado (permiso)
- Cambio de roles
- Cambio de contraseña
- Descarga de reportes

**Información Capturada:**
```json
{
  "timestamp": "2026-05-18T14:35:00Z",
  "user_id": "user_123",
  "user_email": "juan@empresa.com",
  "action": "update",
  "entity_type": "observations",
  "entity_id": "obs_001",
  "old_values": {"status": "mixed"},
  "new_values": {"status": "safe"},
  "ip_address": "192.168.1.100",
  "user_agent": "Mozilla/5.0...",
  "status": "success",
  "duration_ms": 245
}
```

### 4.2 Structured Logging

**Formato:**
```json
{
  "level": "INFO",
  "timestamp": "2026-05-18T14:35:00.123Z",
  "request_id": "req_abc123",
  "user_id": "user_123",
  "method": "POST",
  "path": "/api/observations",
  "status": 201,
  "duration_ms": 145,
  "message": "Observation created successfully"
}
```

**Niveles:**
- DEBUG: Desarrollo solo
- INFO: Eventos importantes
- WARN: Situaciones anómalas
- ERROR: Errores recuperables
- CRITICAL: Errores no recuperables

**Rotación:**
- Diaria
- Máximo 100MB por archivo
- Retención: 90 días
- Compresión: gzip

---

## 5. PROTECCIÓN CONTRA ATAQUES

### 5.1 CORS (Cross-Origin Resource Sharing)

```go
// Configuración restrictiva
AllowedOrigins: []string{
  "https://app.empresa.com",
  "https://admin.empresa.com",
}
AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
AllowedHeaders: []string{"Authorization", "Content-Type"}
ExposedHeaders: []string{"X-Total-Count", "X-Page"}
MaxAge: 3600  // 1 hora
```

### 5.2 HTTPS y Headers de Seguridad

```
Strict-Transport-Security: max-age=31536000; includeSubDomains
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Content-Security-Policy: default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'
Referrer-Policy: strict-origin-when-cross-origin
Permissions-Policy: geolocation=(), microphone=(), camera=()
```

### 5.3 Protección contra Fuerza Bruta

**Login:**
- Máximo 5 intentos fallidos
- Bloqueo progresivo:
  - Intento 1-3: Sin espera
  - Intento 4: Espera 5 segundos
  - Intento 5: Espera 30 minutos
- CAPTCHA después de 3 intentos
- Notificación email al usuario

### 5.4 Dependency Injection para Testing

```go
type UserRepository interface {
  GetByEmail(ctx context.Context, email string) (*User, error)
  Create(ctx context.Context, user *User) error
}

// Mock para tests
type MockUserRepository struct{}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
  // Implementación mock
}
```

---

## 6. MANEJO DE ERRORES SEGUROS

### 6.1 Respuestas de Error

**NO EXPONER:**
```json
// ❌ Malo
{
  "error": "UNIQUE constraint failed: users.email"
}
```

**EXPONER:**
```json
// ✅ Correcto
{
  "code": "VALIDATION_ERROR",
  "message": "El email ya está registrado",
  "field": "email"
}
```

### 6.2 Códigos de Error

```
AUTH_001: Invalid credentials
AUTH_002: Token expired
AUTH_003: Invalid token
AUTH_004: Permission denied
AUTH_005: Token revoked

VALIDATION_001: Email already registered
VALIDATION_002: Invalid email format
VALIDATION_003: Weak password

DB_001: Database connection error
DB_002: Record not found
DB_003: Duplicate entry

SYSTEM_001: Internal server error
```

---

## 7. SEGURIDAD EN DESARROLLO

### 7.1 Gestión de Secretos

**NO hacer:**
```go
const JWTSecret = "my-super-secret-key"  // ❌ Hardcoded
```

**Hacer:**
```go
jwtSecret := os.Getenv("JWT_SECRET")     // ✅ Variables entorno
// O usar: AWS Secrets Manager, HashiCorp Vault
```

### 7.2 Control de Versiones

```
.env                    # Local, NO commitar
.env.example            # Template, SÍ commitar
.env.production         # Usar sistema de secrets
secrets/                # Carpeta ignorada
.gitignore              # Excluir archivos sensibles
```

### 7.3 Análisis Estático

```bash
# Go vet
go vet ./...

# Go fmt
go fmt ./...

# Linter
golangci-lint run

# SAST (Static Application Security Testing)
gosec ./...
```

---

## 8. COMPLIANCE Y REGULACIONES

### 8.1 GDPR Compliance

- Derecho al olvido (right to be forgotten)
- Portabilidad de datos
- Consentimiento registrado
- Política de privacidad clara
- DPA (Data Processing Agreement)

### 8.2 Auditoría Interna

- Logs no pueden ser modificados
- Cadena de custodia de datos
- Segregación de ambientes
- Acceso a logs solo para admin

### 8.3 Pruebas de Seguridad

- OWASP Top 10 compliance
- Penetration testing trimestral
- Vulnerability scanning mensual
- Security audit anual

---

## 9. RESPUESTA A INCIDENTES

### 9.1 Plan de Incidencia

**Niveles:**
- Crítico: Breach confirmado
- Alto: Intentos de ataque
- Medio: Anomalía detectada
- Bajo: Advertencia

**Procedimiento:**
1. Detección automática + alertas
2. Aislamiento de sistema afectado
3. Análisis de impacto
4. Notificación a stakeholders
5. Remediación
6. Comunicación externa
7. Post-mortem

### 9.2 Alertas Automáticas

```
- 10+ login fallos en 5 min → Alertar admin
- Múltiples DELETE en audit_logs → Alertar CISO
- Token revocation en masa → Investigar
- Acceso a datos fuera de permiso → Bloquear + alertar
```

---

## 10. CHECKLIST DE SEGURIDAD

- [ ] Todas las contraseñas hasheadas con bcrypt (cost 12+)
- [ ] JWT validado en cada request
- [ ] HTTPS habilitado en producción
- [ ] CORS configurado restrictivamente
- [ ] Rate limiting activo
- [ ] Audit logs para todas las operaciones sensibles
- [ ] Validación de entrada en todos los endpoints
- [ ] SQL Injection testing pasado
- [ ] XSS testing pasado
- [ ] CSRF tokens en forms
- [ ] Secrets no en código fuente
- [ ] Headers de seguridad presentes
- [ ] Error messages seguros (no exponen detalles internos)
- [ ] Dependencies auditadas (go mod audit)
- [ ] SAST scan limpio (gosec)
- [ ] Penetration testing realizado
- [ ] Disaster recovery plan documentado
- [ ] Backup strategy implementada
