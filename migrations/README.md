# 🗄️ Migraciones de Base de Datos

Esta carpeta contiene todo lo relacionado con las migraciones de base de datos del proyecto Yakka Backend.

## 📁 Estructura

```
migrations/
├── cmd/
│   └── migrate/
│       └── main.go          # Comando de migración
├── scripts/
│   ├── migrate.sh           # Script de migración
│   └── init-db.sql         # Inicialización de BD
└── README.md               # Esta documentación
```

## 🚀 Uso

### Migración Rápida
```bash
# Desarrollo
./migrations/scripts/migrate.sh dev

# Producción
./migrations/scripts/migrate.sh prod
```

### Comando Directo
```bash
# Desarrollo
ENVIRONMENT=development go run migrations/cmd/migrate/main.go

# Producción
ENVIRONMENT=production go run migrations/cmd/migrate/main.go
```

### Desde Makefile
```bash
# Desarrollo
make migrate-dev

# Producción
make migrate-prod
```

## 🔧 Configuración

### Variables de Entorno Requeridas
- `DB_HOST` - Host de la base de datos
- `DB_PORT` - Puerto de la base de datos
- `DB_USER` - Usuario de la base de datos
- `DB_PASSWORD` - Contraseña de la base de datos
- `DB_NAME` - Nombre de la base de datos
- `DB_SSLMODE` - Modo SSL

### Archivos de Configuración
- `.env.dev` - Configuración de desarrollo
- `.env.prod` - Configuración de producción

## 📊 Tablas Creadas

Las migraciones crean automáticamente las siguientes tablas:

- `users` - Usuarios del sistema
- `auth_users` - Usuarios de autenticación
- `sessions` - Sesiones de usuario
- `password_resets` - Tokens de restablecimiento
- `email_verifications` - Verificaciones de email

## 🔄 Flujo de Migración

1. **Conexión**: Se conecta a la base de datos
2. **Validación**: Verifica que la conexión sea exitosa
3. **Migración**: Ejecuta las migraciones de GORM
4. **Confirmación**: Confirma que las migraciones se completaron

## 🚨 Solución de Problemas

### Error de Conexión
```bash
# Verificar que PostgreSQL esté corriendo
docker ps | grep postgres

# Verificar credenciales en .env.dev o .env.prod
```

### Error de Permisos
```bash
# Hacer ejecutable el script
chmod +x migrations/scripts/migrate.sh
```

### Error de Base de Datos
```bash
# Recrear base de datos
docker-compose down
docker-compose up -d postgres
./migrations/scripts/migrate.sh dev
```

## 📝 Notas

- Las migraciones son **idempotentes** (se pueden ejecutar múltiples veces)
- GORM maneja automáticamente los cambios en las tablas
- Los datos existentes se preservan durante las migraciones
- Siempre hacer backup antes de migraciones en producción
