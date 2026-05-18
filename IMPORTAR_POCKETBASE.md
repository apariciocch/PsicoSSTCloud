# Importar Colecciones en PocketBase

## 📥 Pasos para Importar

### 1. Descargar PocketBase
- Ve a: https://pocketbase.io/
- Descarga la versión para tu sistema operativo (Windows, Mac, Linux)
- Descomprime en una carpeta

### 2. Iniciar PocketBase
```bash
# En Windows
pocketbase.exe serve

# En Mac/Linux
./pocketbase serve
```

Debe mostrar:
```
Server started at:
  API: http://127.0.0.1:8090
  Admin UI: http://127.0.0.1:8090/_/
```

### 3. Abrir Admin Panel
- URL: http://localhost:8090/_/
- Crea usuario admin (primera vez)

### 4. Importar Colecciones

**Opción A: Importar el JSON completo (RECOMENDADO)**

1. Click en **Settings** (engranaje)
2. Click en **Import collections**
3. Abre el archivo: `pocketbase/collections.json`
4. O copia el contenido y pégalo
5. Click en **Import**

✅ Se crearán todas las 14 colecciones automáticamente

**Opción B: Crear manualmente**

Si prefieres, puedes ir creando cada colección manualmente en el admin panel.

---

## 📁 Archivos JSON

### Ubicación
```
PsicoSSTCloud/
└── pocketbase/
    └── collections.json          ← Importa esto
```

### Contenido
14 colecciones:
- roles
- users
- work_areas
- employees
- behavior_categories
- observations
- behaviors
- incidents
- corrective_actions
- attachments
- notifications
- audit_logs
- kpis
- reports

---

## ✅ Verificar Importación

En el admin panel, debe aparecer:
```
Collections (14)
├── roles
├── users
├── work_areas
├── employees
├── behavior_categories
├── observations
├── behaviors
├── incidents
├── corrective_actions
├── attachments
├── notifications
├── audit_logs
├── kpis
└── reports
```

---

**Listo. Eso es todo lo que necesitas.**
