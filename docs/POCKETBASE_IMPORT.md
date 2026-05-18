# Guía de Importación de Colecciones en PocketBase

## 🚀 Método 1: Importar JSON Directamente (Más Fácil)

### Paso 1: Iniciar PocketBase
```bash
cd /workspaces/PsicoSSTCloud
docker-compose up -d pocketbase
```

Espera 5-10 segundos a que inicie.

### Paso 2: Acceder al Admin Panel
- URL: http://localhost:8090/_/
- Email: admin@example.com (crear nueva)
- Password: secure_password

### Paso 3: Importar Colecciones

1. En el panel admin, click en **"Settings"** (engranaje)
2. Click en **"Import collections"**
3. Copiar el contenido de `/pocketbase/collections.json`
4. Pegarlo en el campo de texto
5. Click en **"Import"**

✅ **Todas las 14 colecciones se crearán automáticamente**

---

## 🔧 Método 2: Crear Manualmente (Alternativa)

Si prefieres crear las colecciones manualmente:

### Paso 1: Roles
1. Click en "**+" → New Collection**
2. Nombre: `roles`
3. Agregar campos:
   - `name` (text, required, unique)
   - `slug` (text, required, unique)
   - `description` (text, required)
   - `permissions` (json, required)
   - `level` (number, required)
   - `is_system_role` (bool)
4. Click **Save**

### Paso 2: Users
1. Click en "**+" → New Collection**
2. Nombre: `users`
3. Agregar campos:
   - `email` (email, required, unique)
   - `username` (text, required, unique)
   - `password_hash` (text, required)
   - `role_id` (relation → roles, required)
   - `is_active` (bool)
   - `failed_attempts` (number)
   - `locked_until` (date)
4. Click **Save**

...y así para cada colección.

---

## 📋 Colecciones Incluidas

| # | Nombre | Descripción | Registros Iniciales |
|---|--------|-------------|-------------------|
| 1 | `roles` | 5 roles del sistema | Admin, Supervisor, Observer, Worker, Auditor |
| 2 | `users` | Usuarios del sistema | Admin de prueba |
| 3 | `work_areas` | Áreas de trabajo | Según tu empresa |
| 4 | `employees` | Empleados | Vinculados a áreas |
| 5 | `behavior_categories` | Categorías de conductas | Safe / Unsafe |
| 6 | `observations` | Registros de observaciones | Vacía (se llena con API) |
| 7 | `behaviors` | Conductas observadas | Relacionada con observations |
| 8 | `incidents` | Incidentes registrados | Vacía (se llena con API) |
| 9 | `corrective_actions` | Acciones correctivas | Vacía (se llena con API) |
| 10 | `attachments` | Archivos adjuntos | Vacía |
| 11 | `notifications` | Notificaciones | Vacía (se llena automática) |
| 12 | `audit_logs` | Auditoría completa | Vacía (se llena automática) |
| 13 | `kpis` | Indicadores de desempeño | Vacía (se calcula automática) |
| 14 | `reports` | Reportes generados | Vacía |

---

## 🔐 Roles Iniciales

Se crean automáticamente:

```javascript
{
  "admin": {
    "permissions": ["*"],  // Acceso total
    "level": 1
  },
  "supervisor": {
    "permissions": ["users.view", "observations.create", "reports.view"],
    "level": 2
  },
  "observer": {
    "permissions": ["observations.create", "behaviors.view"],
    "level": 3
  },
  "worker": {
    "permissions": ["observations.view"],
    "level": 4
  },
  "auditor": {
    "permissions": ["audit_logs.view", "reports.view"],
    "level": 5
  }
}
```

---

## 👤 Usuario Admin Inicial

```
Email: admin@example.com
Password: secure_password
Role: Admin (acceso total)
```

⚠️ **IMPORTANTE**: Cambiar la contraseña en producción

---

## ✅ Validar Importación

Después de importar, verificar:

1. **En PocketBase Admin Panel:**
   - ✅ Aparecen 14 colecciones en la lista
   - ✅ Cada colección tiene sus campos (click en colección)
   - ✅ Las relaciones están configuradas

2. **Desde Terminal:**
```bash
# Listar colecciones
curl -s http://localhost:8090/api/collections | jq '.items | length'
# Debe mostrar: 14

# Contar registros en roles
curl -s http://localhost:8090/api/collections/roles/records | jq '.items | length'
# Debe mostrar: 5 (los 5 roles iniciales)
```

---

## 🔗 Relaciones Entre Colecciones

```
roles ←─────────┐
        ↓       │
      users     │
        │       │
        └─ role_id
        
employees ─ work_area_id ─→ work_areas

observations ─┬─ employee_id ─→ employees
              ├─ observer_id ─→ users
              └─ work_area_id ─→ work_areas
              
behaviors ─┬─ observation_id ─→ observations
           └─ behavior_category_id ─→ behavior_categories

incidents ─┬─ employee_id ─→ employees
           ├─ work_area_id ─→ work_areas
           └─ reported_by_id ─→ users

corrective_actions ─┬─ incident_id ─→ incidents
                    └─ responsible_id ─→ users

notifications ─ recipient_id ─→ users

audit_logs ─ user_id ─→ users

kpis ─ work_area_id ─→ work_areas

reports ─┬─ work_area_id ─→ work_areas
         └─ generated_by_id ─→ users

attachments ─ uploaded_by_id ─→ users
```

---

## 🧪 Test de Conexión Backend ↔ PocketBase

Una vez importadas las colecciones:

```bash
# 1. Iniciar el backend
docker-compose up -d

# 2. Esperar 10 segundos
sleep 10

# 3. Verificar conexión
curl -v http://localhost:8080/health
# Debe responder: 200 OK

# 4. Intentar login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "secure_password"
  }'
  
# Respuesta esperada:
# {
#   "access_token": "eyJhbGci...",
#   "refresh_token": "eyJhbGci...",
#   "user": {...}
# }
```

---

## 🚀 Agregar Datos Iniciales

### Crear Roles (si no se importaron)
```bash
curl -X POST http://localhost:8090/api/collections/roles/records \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Administrador",
    "slug": "admin",
    "description": "Acceso total al sistema",
    "permissions": ["*"],
    "level": 1,
    "is_system_role": true
  }'
```

### Crear Áreas de Trabajo
```bash
curl -X POST http://localhost:8090/api/collections/work_areas/records \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Producción",
    "description": "Planta de producción principal",
    "location": "Piso 1",
    "risk_level": "high",
    "is_active": true
  }'
```

### Crear Categorías de Conductas
```bash
curl -X POST http://localhost:8090/api/collections/behavior_categories/records \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Uso de EPP",
    "description": "Equipos de Protección Personal",
    "category_type": "safe",
    "color_code": "#00AA00"
  }'
```

---

## 📚 Archivos Relacionados

- `pocketbase/collections.json` - Definición de colecciones
- `backend/.env.example` - Variables de entorno
- `docs/COLLECTIONS.md` - Descripción detallada de colecciones
- `docs/API_ENDPOINTS.md` - Endpoints para acceder a datos

---

## ❓ Troubleshooting

### Error: "Cannot connect to PocketBase"
```bash
# Verificar que está corriendo
docker ps | grep pocketbase

# Reiniciar
docker restart pocketbase

# Esperar 10 segundos y reintentar
sleep 10
```

### Error: "Collection already exists"
```bash
# Si importas dos veces, las colecciones ya existen
# Opción 1: Borrar todas y reimportar
# Opción 2: Solo crear las que faltan manualmente
```

### Error: "Invalid JSON"
```bash
# Validar JSON en: https://jsonlint.com/
# Asegurarse de copiar TODO el contenido de collections.json
```

### PocketBase no inicia
```bash
# Revisar logs
docker-compose logs pocketbase

# Reconstruir imagen
docker-compose build --no-cache pocketbase

# Reiniciar
docker-compose up -d pocketbase
```

---

## 🎯 Próximos Pasos

Después de importar las colecciones:

1. ✅ Iniciar backend: `docker-compose up -d backend`
2. ✅ Crear usuario admin: En PocketBase admin panel
3. ✅ Crear áreas de trabajo: Usando curl o PocketBase admin
4. ✅ Crear empleados: Usando API `/api/v1/employees`
5. ✅ Comenzar a registrar observaciones: API `/api/v1/observations`

---

**Versión**: 1.0  
**Fecha**: Mayo 18, 2026  
**Estado**: Listo para importar
