# Configuración y Hooks de PocketBase

## 1. INICIALIZACIÓN DE COLECCIONES

### Script de Inicialización (JavaScript - en PocketBase)

```javascript
// pocketbase/pb_migrations/init_collections.js

async function migrate(db) {
  // ==================== ROLES ====================
  const rolesCollection = new Collection({
    name: 'roles',
    type: 'base',
    system: false,
    createRule: null,
    updateRule: null,
    deleteRule: null,
    listRule: null,
    viewRule: null,
    schema: [
      {
        id: 'id',
        name: 'id',
        type: 'text',
        system: true,
        required: true,
        unique: true,
        options: {
          min: null,
          max: null,
          pattern: ''
        }
      },
      {
        id: 'name',
        name: 'name',
        type: 'text',
        required: true,
        options: {
          min: 3,
          max: 100,
          pattern: ''
        }
      },
      {
        id: 'slug',
        name: 'slug',
        type: 'text',
        required: true,
        unique: true,
        options: {
          min: 3,
          max: 50,
          pattern: ''
        }
      },
      {
        id: 'description',
        name: 'description',
        type: 'text',
        required: true,
        options: {
          min: 10,
          max: 500,
          pattern: ''
        }
      },
      {
        id: 'permissions',
        name: 'permissions',
        type: 'json',
        required: true
      },
      {
        id: 'level',
        name: 'level',
        type: 'number',
        required: true,
        options: {
          min: 1,
          max: 100
        }
      },
      {
        id: 'is_system_role',
        name: 'is_system_role',
        type: 'bool',
        required: true
      }
    ],
    indexes: [
      'CREATE UNIQUE INDEX idx_roles_slug ON roles(slug)'
    ]
  });

  try {
    await db.collections.import([rolesCollection]);
  } catch (err) {
    console.log(`Collection "roles" already exists. Error: ${err.message}`);
  }

  // ==================== Insertar Roles Iniciales ====================
  const adminRole = new Record(rolesCollection, {
    name: 'Administrador',
    slug: 'admin',
    description: 'Administrador con acceso total',
    permissions: [
      'users.view', 'users.create', 'users.edit', 'users.delete',
      'observations.view', 'observations.create', 'observations.edit', 'observations.delete',
      'incidents.view', 'incidents.create', 'incidents.edit',
      'reports.generate', 'reports.export',
      'analytics.view',
      'audit_logs.view'
    ],
    level: 1,
    is_system_role: true
  });

  try {
    await db.collection('roles').create(adminRole);
  } catch (err) {
    console.log(`Role "admin" already exists`);
  }

  // ... (más roles)
}
```

## 2. HOOKS Y VALIDACIONES CUSTOMIZADAS

### Hook: Auto-generar números secuenciales

```javascript
// Hook para generar observation_number automáticamente
$app.dao().registerRecordSubscriber('observations', (e) => {
  if (e.isCreate) {
    // Obtener el próximo número
    const lastObs = $app.dao()
      .findRecordsByFilter('observations', '', '-observation_number', 1)[0];
    
    let nextNumber = 1;
    if (lastObs) {
      const match = lastObs.get('observation_number').match(/OBS-\d+-(\d+)/);
      if (match) {
        nextNumber = parseInt(match[1]) + 1;
      }
    }

    const year = new Date().getFullYear();
    const obsNumber = `OBS-${year}-${String(nextNumber).padStart(5, '0')}`;
    e.record.set('observation_number', obsNumber);
  }
});
```

### Hook: Validación de conductas

```javascript
// Hook para validar que una observación tenga coherencia
$app.dao().registerRecordSubscriber('observations', (e) => {
  if (e.isCreate || e.isUpdate) {
    const hasSafe = e.record.get('has_safe_behaviors');
    const hasUnsafe = e.record.get('has_unsafe_behaviors');
    const status = e.record.get('overall_status');

    // Validar coherencia
    if (hasSafe && !hasUnsafe && status !== 'safe') {
      throw new Error('Overall status must be "safe" when only safe behaviors are present');
    }
    if (hasUnsafe && !hasSafe && status !== 'unsafe') {
      throw new Error('Overall status must be "unsafe" when only unsafe behaviors are present');
    }
  }
});
```

### Hook: Auditoría automática

```javascript
// Hook para registrar todos los cambios en audit_logs
$app.dao().registerRecordSubscriber('*', (e) => {
  // Ignorar colecciones del sistema
  if (e.collection.name.startsWith('_pb_')) {
    return;
  }

  const auditLog = new Record(
    $app.dao().findCollectionByNameOrId('audit_logs'),
    {
      user_id: $app.auth?.record?.id || 'system',
      action: e.isCreate ? 'create' : e.isUpdate ? 'update' : 'delete',
      entity_type: e.collection.name,
      entity_id: e.record.id,
      old_values: e.isUpdate ? e.oldRecord.export() : null,
      new_values: e.record.export(),
      status: 'success',
      ip_address: $app.request?.header?.get('X-Forwarded-For') || 
                  $app.request?.remoteAddr || 'unknown',
      user_agent: $app.request?.header?.get('User-Agent') || 'unknown',
    }
  );

  try {
    $app.dao().saveRecord(auditLog);
  } catch (err) {
    console.error('Error saving audit log:', err);
  }
});
```

### Hook: Crear notificación para acción correctiva vencida

```javascript
// Hook para verificar acciones vencidas
$app.cron().add('check_overdue_actions', '0 * * * *', () => {
  const now = new Date().toISOString();
  
  const overdueActions = $app.dao()
    .findRecordsByFilter('corrective_actions',
      `due_date <= "${now}" && status != "completed" && status != "overdue"`
    );

  for (const action of overdueActions) {
    // Actualizar estado
    action.set('status', 'overdue');
    $app.dao().saveRecord(action);

    // Crear notificación
    const notification = new Record(
      $app.dao().findCollectionByNameOrId('notifications'),
      {
        recipient_id: action.get('responsible_id'),
        notification_type: 'critical',
        title: 'Acción Correctiva Vencida',
        message: `La acción correctiva "${action.get('description')}" ha vencido.`,
        related_entity_type: 'corrective_actions',
        related_entity_id: action.id,
        is_read: false,
      }
    );

    $app.dao().saveRecord(notification);
  }
});
```

### Hook: KPI Auto-calculado

```javascript
// Hook para recalcular KPIs cuando hay cambios en observaciones
$app.dao().registerRecordSubscriber('observations', (e) => {
  if (!e.isCreate) return;

  const workAreaId = e.record.get('work_area_id');
  const today = new Date().toISOString().split('T')[0];

  // Consultar observaciones de hoy
  const todayRecords = $app.dao()
    .findRecordsByFilter('observations',
      `DATE(observation_date) = DATE("${today}") && work_area_id = "${workAreaId}"`
    );

  let safeBehaviors = 0;
  let unsafeBehaviors = 0;

  for (const obs of todayRecords) {
    const behaviors = $app.dao()
      .findRecordsByFilter('behaviors', `observation_id = "${obs.id}"`);
    
    for (const behavior of behaviors) {
      if (behavior.get('type') === 'safe') safeBehaviors++;
      else unsafeBehaviors++;
    }
  }

  const total = safeBehaviors + unsafeBehaviors;
  const safePercentage = total > 0 ? (safeBehaviors / total) * 100 : 0;

  // Actualizar o crear KPI del día
  let kpi = $app.dao()
    .findFirstRecordByFilter('kpis',
      `period = "daily" && period_start = "${today}" && work_area_id = "${workAreaId}"`
    );

  if (!kpi) {
    kpi = new Record($app.dao().findCollectionByNameOrId('kpis'));
    kpi.set('period', 'daily');
    kpi.set('period_start', today);
    kpi.set('period_end', today);
    kpi.set('work_area_id', workAreaId);
    kpi.set('kpi_code', `KPI-${today}-${workAreaId}`);
    kpi.set('kpi_name', `Daily KPI - ${today}`);
  }

  kpi.set('safe_behaviors_count', safeBehaviors);
  kpi.set('unsafe_behaviors_count', unsafeBehaviors);
  kpi.set('safe_percentage', safePercentage);
  kpi.set('unsafe_percentage', 100 - safePercentage);
  kpi.set('total_observations', todayRecords.length);
  kpi.set('calculated_at', new Date().toISOString());

  $app.dao().saveRecord(kpi);
});
```

## 3. REGLAS DE ACCESO

### Reglas de Lectura Granular

```
// users collection - View rule
user.id = @requestUser.id
||
@requestUser.is_super_admin = true
||
@requestUser.role_id in (select id from roles where slug = 'admin')
||
(
  @requestUser.role_id in (select id from roles where slug = 'supervisor')
  &&
  user.role_id in (select id from roles where slug != 'admin')
)
```

### Reglas de Creación

```
// observations - Create rule
@requestUser.role_id in (
  select id from roles 
  where slug in ('admin', 'supervisor', 'observer')
)
```

## 4. ÍNDICES DE PERFORMANCE

```sql
-- Índices críticos
CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE UNIQUE INDEX idx_users_username ON users(username);
CREATE INDEX idx_observations_employee_id ON observations(employee_id);
CREATE INDEX idx_observations_observer_id ON observations(observer_id);
CREATE INDEX idx_observations_work_area_id ON observations(work_area_id);
CREATE INDEX idx_observations_observation_date ON observations(observation_date DESC);
CREATE INDEX idx_behaviors_observation_id ON behaviors(observation_id);
CREATE INDEX idx_incidents_employee_id ON incidents(employee_id);
CREATE INDEX idx_incidents_work_area_id ON incidents(work_area_id);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);
```

## 5. BACKUPS AUTOMÁTICOS

```javascript
// Configurar backup automático cada día a las 2:00 AM
$app.cron().add('daily_backup', '0 2 * * *', () => {
  const backupDir = '/pb_backups';
  const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
  const backupPath = `${backupDir}/backup-${timestamp}.zip`;

  // Usar comando del sistema para backup
  // Implementación específica según SO
  console.log(`Backup created at ${backupPath}`);
});
```

## 6. LIMPIEZA DE DATOS EXPIRADOS

```javascript
// Limpiar logs de auditoría más antiguos de 90 días
$app.cron().add('cleanup_old_logs', '0 3 * * 0', () => {
  const ninetyDaysAgo = new Date();
  ninetyDaysAgo.setDate(ninetyDaysAgo.getDate() - 90);
  const cutoffDate = ninetyDaysAgo.toISOString();

  const oldLogs = $app.dao()
    .findRecordsByFilter('audit_logs', `created_at < "${cutoffDate}"`);

  for (const log of oldLogs) {
    $app.dao().deleteRecord(log);
  }

  console.log(`Deleted ${oldLogs.length} old audit logs`);
});
```

## 7. VALIDACIONES CUSTOMIZADAS

```javascript
// Validación de formato de documento
$app.validator.addRule('document_number', (value) => {
  // Validar que el documento tenga formato válido
  if (!value || typeof value !== 'string') {
    return false;
  }
  
  // Aceptar números y guiones
  return /^[0-9\-]+$/.test(value) && value.length >= 5;
});

// Validación de teléfono internacional
$app.validator.addRule('phone_international', (value) => {
  if (!value) return true; // opcional
  return /^\+?[0-9\s\-()]{10,}$/.test(value);
});
```

## 8. TESTING DE HOOKS

```javascript
// Script de test para verificar hooks
describe('PocketBase Hooks', () => {
  it('should auto-generate observation numbers', async () => {
    // Crear observación de prueba
    const obs = new Record($pb.collection('observations'));
    obs.set('employee_id', 'emp_test');
    obs.set('observer_id', 'user_test');
    
    const saved = await $pb.collection('observations').create(obs);
    expect(saved.observation_number).toMatch(/^OBS-\d+-\d+$/);
  });

  it('should validate observation status consistency', async () => {
    const obs = new Record($pb.collection('observations'));
    obs.set('has_safe_behaviors', true);
    obs.set('has_unsafe_behaviors', false);
    obs.set('overall_status', 'unsafe'); // Incoherente

    try {
      await $pb.collection('observations').create(obs);
      fail('Should have thrown error');
    } catch (e) {
      expect(e.message).toContain('Overall status must be');
    }
  });
});
```

---

## Deployment en Producción

### Docker Compose con PocketBase

El archivo `docker-compose.yml` ya contiene la configuración completa de PocketBase.

**Iniciar:**
```bash
docker-compose up -d pocketbase
```

**Acceder a dashboard:**
```
http://localhost:8090/_/
```

**Usar CLI de PocketBase:**
```bash
docker exec pocketbase pocketbase --help
```
