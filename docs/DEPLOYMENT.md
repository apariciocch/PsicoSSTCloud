# Guía de Despliegue - SBC Backend

## 1. PREREQUISITES

### Sistema Operativo
- Linux (Ubuntu 20.04+, Debian 11+)
- macOS (Big Sur+)
- Windows + WSL2

### Software Requerido
- Docker 20.10+
- Docker Compose 2.0+
- Git
- Go 1.21+ (si compilas manualmente)

### Recursos Recomendados
- **CPU**: 2+ cores
- **RAM**: 4GB mínimo, 8GB recomendado
- **Almacenamiento**: 20GB mínimo

---

## 2. DESPLIEGUE LOCAL (Desarrollo)

### 2.1 Clonar el Repositorio

```bash
git clone https://github.com/apariciocch/PsicoSSTCloud.git
cd PsicoSSTCloud
```

### 2.2 Configurar Variables de Entorno

```bash
# Copiar archivo de ejemplo
cp backend/.env.example backend/.env

# Editar con tus valores
nano backend/.env
```

**Variables críticas:**
```
JWT_SECRET=tu_clave_secreta_de_32_caracteres_minimo_aqui_xxx
POCKETBASE_ADMIN_PASSWORD=tu_contraseña_segura_aqui
ENVIRONMENT=development
```

### 2.3 Iniciar con Docker Compose

```bash
# Iniciar todos los servicios
docker-compose up -d

# Verificar estado
docker-compose ps

# Ver logs
docker-compose logs -f backend
docker-compose logs -f pocketbase
```

### 2.4 Inicializar Base de Datos

```bash
# Acceder a PocketBase Admin
# http://localhost:8090/_/

# Crear colecciones manualmente o usar script
# Ver POCKETBASE_CONFIG.md para más detalles
```

### 2.5 Verificar Servicios

```bash
# Health check del backend
curl http://localhost:8080/health

# Health check de PocketBase
curl http://localhost:8090/api/health

# Intentar login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@example.com", "password": "..."}'
```

---

## 3. DESPLIEGUE EN PRODUCCIÓN

### 3.1 Opción A: Docker Compose en Servidor

#### Requisitos
- Servidor VPS o Cloud (AWS EC2, DigitalOcean, Linode, etc.)
- Domain name con DNS configurado
- SSL Certificate (Let's Encrypt)

#### Pasos

1. **Conectar al servidor**
```bash
ssh usuario@ip_servidor
```

2. **Instalar Docker y Docker Compose**
```bash
# Ubuntu/Debian
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

3. **Clonar repositorio**
```bash
cd /opt
sudo git clone https://github.com/apariciocch/PsicoSSTCloud.git
cd PsicoSSTCloud
sudo chown -R $USER:$USER .
```

4. **Configurar variables de producción**
```bash
# Crear .env con valores seguros
cat > backend/.env << EOF
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
ENVIRONMENT=production

POCKETBASE_URL=http://pocketbase:8090
POCKETBASE_ADMIN_EMAIL=admin@empresa.com
POCKETBASE_ADMIN_PASSWORD=$(openssl rand -base64 32)

JWT_SECRET=$(openssl rand -base64 32)
JWT_EXPIRATION=900
JWT_REFRESH_EXPIRATION=604800

CORS_ORIGINS=https://app.empresa.com

LOG_LEVEL=warn
LOG_FORMAT=json
EOF
```

5. **Configurar Nginx como reverse proxy**
```bash
# Instalar Nginx
sudo apt-get update
sudo apt-get install -y nginx

# Crear configuración
sudo tee /etc/nginx/sites-available/psicosst > /dev/null << EOF
upstream backend {
    server 127.0.0.1:8080;
}

upstream pocketbase {
    server 127.0.0.1:8090;
}

server {
    listen 80;
    server_name api.empresa.com;

    # Redirigir a HTTPS
    return 301 https://\$server_name\$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.empresa.com;

    ssl_certificate /etc/letsencrypt/live/api.empresa.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.empresa.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # GZIP compression
    gzip on;
    gzip_types text/plain text/css application/json application/javascript;

    # Backend API
    location /api {
        proxy_pass http://backend;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # PocketBase Admin (protegido)
    location /admin {
        proxy_pass http://pocketbase/_/;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
    }

    # Health checks
    location /health {
        proxy_pass http://backend;
    }
}
EOF

# Habilitar sitio
sudo ln -s /etc/nginx/sites-available/psicosst /etc/nginx/sites-enabled/

# Validar configuración
sudo nginx -t

# Iniciar Nginx
sudo systemctl start nginx
sudo systemctl enable nginx
```

6. **Instalar SSL Certificate (Let's Encrypt)**
```bash
# Instalar certbot
sudo apt-get install -y certbot python3-certbot-nginx

# Generar certificado
sudo certbot certonly --standalone -d api.empresa.com

# Renovación automática
sudo systemctl enable certbot.timer
```

7. **Iniciar servicios**
```bash
docker-compose -f docker-compose.prod.yml up -d
```

8. **Configurar Firewall**
```bash
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

### 3.2 Opción B: Render (PaaS)

#### Conectar PocketBase a Base de Datos Externa

```yaml
# Render PostgreSQL
Name: psicosst-db
Instance Type: Standard
Plan: $7/month (10GB)
```

#### Crear servicio Go backend

1. Conectar repositorio GitHub
2. Build Command: `go build -o backend ./cmd/main.go`
3. Start Command: `./backend`
4. Environment Variables:
```
DATABASE_URL=postgresql://user:pass@host/db
POCKETBASE_URL=https://pocketbase.render.com
JWT_SECRET=tu_secreto_aqui
```

#### Deploy PocketBase separadamente
```bash
# Push Docker image a registry
docker tag pocketbase:latest your-registry.com/pocketbase:latest
docker push your-registry.com/pocketbase:latest

# Crear Render deployment desde imagen
```

### 3.3 Opción C: Kubernetes

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: psicosst-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: psicosst-backend
  template:
    metadata:
      labels:
        app: psicosst-backend
    spec:
      containers:
      - name: backend
        image: your-registry.com/psicosst-backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: POCKETBASE_URL
          value: "http://pocketbase:8090"
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: psicosst-secrets
              key: jwt-secret
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5

---
apiVersion: v1
kind: Service
metadata:
  name: psicosst-backend
spec:
  type: LoadBalancer
  selector:
    app: psicosst-backend
  ports:
  - port: 80
    targetPort: 8080
```

Deploy:
```bash
kubectl apply -f deployment.yaml
```

---

## 4. MONITOREO Y MANTENIMIENTO

### 4.1 Logs

```bash
# Ver logs en tiempo real
docker-compose logs -f backend

# Ver últimas 100 líneas
docker-compose logs --tail=100 backend

# Guardar logs en archivo
docker-compose logs backend > logs.txt
```

### 4.2 Backup de Base de Datos

```bash
# Backup automático usando cron
# Agregar a crontab: 0 2 * * * /opt/scripts/backup_db.sh

#!/bin/bash
# /opt/scripts/backup_db.sh
BACKUP_DIR="/backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

docker exec pocketbase sqlite3 /pb_data/data.db ".backup /pb_data/backup_$TIMESTAMP.db"
docker cp pocketbase:/pb_data/backup_$TIMESTAMP.db $BACKUP_DIR/

# Eliminar backups más antiguos de 30 días
find $BACKUP_DIR -name "backup_*.db" -mtime +30 -delete
```

### 4.3 Performance Tuning

```bash
# Aumentar conexiones máximas de PostgreSQL (si aplica)
# En postgresql.conf:
# max_connections = 200
# shared_buffers = 256MB
# effective_cache_size = 1GB

# Configurar límites de recursos en Docker
docker update --memory 4g --cpus 2 psicosst-backend
```

### 4.4 Monitoreo de Salud

```bash
# Script de monitoreo
#!/bin/bash
API_URL="https://api.empresa.com/health"

RESPONSE=$(curl -s -w "\n%{http_code}" $API_URL)
HTTP_CODE=$(echo "$RESPONSE" | tail -1)

if [ "$HTTP_CODE" != "200" ]; then
    # Enviar alerta
    mail -s "Alert: API Down" admin@empresa.com
    # O usar webhook:
    curl -X POST https://hooks.slack.com/... \
      -d '{"text":"API is down!"}'
fi
```

---

## 5. TROUBLESHOOTING

### Error: Connection Refused

```bash
# Verificar que PocketBase esté corriendo
docker ps | grep pocketbase

# Reiniciar PocketBase
docker restart pocketbase

# Verificar logs
docker logs pocketbase
```

### Error: Out of Memory

```bash
# Aumentar límite de memoria Docker
docker update --memory 8g psicosst-backend

# O en docker-compose.yml:
# mem_limit: 8g
```

### Permisos en PocketBase

```bash
# Resync colecciones
docker exec pocketbase pocketbase migrate collections

# Reset permisos
docker exec pocketbase pocketbase admin create admin@empresa.com password
```

### SSL Certificate Issues

```bash
# Renovar certificado Let's Encrypt
sudo certbot renew --dry-run

# Renovar forzado
sudo certbot renew --force-renewal
```

---

## 6. SCALING

### Horizontal Scaling (múltiples instancias)

```yaml
# docker-compose.prod.yml con múltiples backends
services:
  backend_1:
    image: psicosst-backend
    environment:
      INSTANCE_ID: "1"
    
  backend_2:
    image: psicosst-backend
    environment:
      INSTANCE_ID: "2"
      
  backend_3:
    image: psicosst-backend
    environment:
      INSTANCE_ID: "3"
      
  pocketbase:
    image: pocketbase
    volumes:
      - pocketbase_shared:/pb_data  # Volumen compartido
```

### Vertical Scaling (más recursos)

```bash
# Actualizar recursos en servidor
docker update --memory 16g --cpus 8 psicosst-backend
```

### Load Balancing

Nginx automáticamente balanceará carga entre múltiples instancias:

```nginx
upstream backend {
    least_conn;
    server backend_1:8080;
    server backend_2:8080;
    server backend_3:8080;
}
```

---

## 7. DISASTER RECOVERY

### 1. Plan de Recuperación
- **RTO (Recovery Time Objective)**: 1 hora
- **RPO (Recovery Point Objective)**: 1 hora

### 2. Backup Strategy
```bash
# Backup diario
0 2 * * * docker exec pocketbase sqlite3 /pb_data/data.db ".backup /pb_data/backup-$(date +%Y%m%d).db"

# Backup a almacenamiento externo (S3, Google Cloud)
0 3 * * * aws s3 sync /backups s3://empresa-backups/
```

### 3. Recuperación
```bash
# Restaurar desde backup
docker exec pocketbase sqlite3 /pb_data/data.db ".restore /pb_data/backup-20260518.db"

# Verificar integridad
docker exec pocketbase sqlite3 /pb_data/data.db "PRAGMA integrity_check;"
```

---

## 8. SECURITY HARDENING

### En Producción

```bash
# Cambiar contraseña de admin
docker exec -it pocketbase pocketbase admin update

# Deshabilitar admin UI en producción (opcional)
# Agregar flag: --disable-admin-ui

# Usar secrets en lugar de .env
# Docker secrets, Vault, Secrets Manager
```

### Firewall Rules

```bash
# Solo permitir desde Nginx
sudo ufw allow from 127.0.0.1 to 127.0.0.1 port 8080
sudo ufw allow from 127.0.0.1 to 127.0.0.1 port 8090

# O desde rango de red
sudo ufw allow from 10.0.0.0/8 to 10.0.0.0/8 port 8080
```

---

## 9. CHECKLIST DE DESPLIEGUE

- [ ] Variables de entorno configuradas
- [ ] SSL/TLS habilitado
- [ ] Firewall configurado
- [ ] Backups programados
- [ ] Monitoreo activado
- [ ] Logs centralizados
- [ ] Health checks funcionando
- [ ] Rate limiting configurado
- [ ] CORS correctamente restringido
- [ ] Base de datos optimizada
- [ ] Nginx caché configurado
- [ ] CDN implementado (opcional)
- [ ] Documentación actualizada
- [ ] Plan de disaster recovery creado
- [ ] Tests de carga realizados

---

## 10. CONTACTO Y SOPORTE

Para problemas o consultas:
- GitHub Issues: https://github.com/apariciocch/PsicoSSTCloud/issues
- Email: soporte@empresa.com
- Documentación: /docs
