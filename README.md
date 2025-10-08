# PsicoSSTCloud

## Configuración rápida

1. Instala las dependencias:
   ```bash
   npm install
   ```
2. Copia el archivo de variables de entorno de ejemplo y actualiza los datos si fuera necesario:
   ```bash
   cp server/.env.example server/.env
   ```
3. Inicia el servidor de desarrollo:
   ```bash
   npm run dev
   ```

## Variables de entorno

La API usa `pg` para conectarse a la base de datos. Define la variable `DATABASE_URL` en `server/.env` con el string de conexión que entrega tu proveedor. Por ejemplo:

```env
DATABASE_URL="postgresql://neondb_owner:npg_qrXE48eofWVC@ep-divine-shadow-adp4gjr5-pooler.c-2.us-east-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"
```

El proyecto ya está preparado para cargar automáticamente estas variables mediante `dotenv`.
