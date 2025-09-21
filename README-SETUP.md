# 🚀 Yakka Backend - Guía de Configuración

Esta guía te ayudará a configurar y ejecutar el proyecto Yakka Backend por primera vez.

## 📋 Prerrequisitos

- **Go 1.21+** instalado
- **Docker** y **Docker Compose** instalados
- **PostgreSQL** (opcional, si no usas Docker)

## 🛠️ Configuración Inicial

### 1. Crear Archivos de Entorno

Crea los archivos `.env.dev` y `.env.prod` en la raíz del proyecto:

#### `.env.dev` (Desarrollo)
```bash
# Development Environment Configuration
ENVIRONMENT=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=yakka_dev
DB_SSLMODE=disable

# Server Configuration
PORT=8080

# Logging Configuration
LOG_LEVEL=debug

# JWT Configuration
JWT_SECRET=your_jwt_secret_key_here
JWT_EXPIRATION_HOURS=24
```

#### `.env.prod` (Producción)
```bash
# Production Environment Configuration
ENVIRONMENT=production

# Database Configuration
DB_HOST=your_production_db_host
DB_PORT=5432
DB_USER=your_production_user
DB_PASSWORD=your_production_password
DB_NAME=yakka_prod
DB_SSLMODE=require

# Server Configuration
PORT=8080

# Logging Configuration
LOG_LEVEL=info

# JWT Configuration
JWT_SECRET=your_production_jwt_secret_key_here
JWT_EXPIRATION_HOURS=24
```

### 2. Instalar Dependencias

```bash
go mod download
go mod tidy
```

## 🐳 Opción 1: Usando Docker (Recomendado)

### Desarrollo Rápido
```bash
# Iniciar solo PostgreSQL
make docker-up

# Ejecutar migraciones
make migrate-dev

# Iniciar aplicación
make dev
```

### Desarrollo Completo
```bash
# Todo en uno: Docker + Migraciones + Aplicación
make quick-start
```

## 💻 Opción 2: Sin Docker

### 1. Configurar PostgreSQL Manualmente
- Instalar PostgreSQL
- Crear bases de datos: `yakka_dev` y `yakka_prod`
- Configurar usuarios y permisos

### 2. Ejecutar Migraciones
```bash
# Desarrollo
make migrate-dev

# Producción
make migrate-prod
```

### 3. Ejecutar Aplicación
```bash
# Desarrollo
make dev

# Producción
make prod
```

## 📊 Comandos Disponibles

### Desarrollo
```bash
make dev              # Iniciar desarrollo con migraciones
make migrate-dev      # Solo migraciones de desarrollo
make docker-up        # Iniciar PostgreSQL con Docker
make docker-down      # Detener Docker
```

### Producción
```bash
make prod             # Iniciar producción
make migrate-prod     # Solo migraciones de producción
```

### Utilidades
```bash
make build            # Compilar aplicación
make test             # Ejecutar tests
make clean            # Limpiar archivos compilados
make deps             # Instalar dependencias
```

## 🔧 Comandos Manuales

### Migraciones
```bash
# Desarrollo
./migrations/scripts/migrate.sh dev

# Producción
./migrations/scripts/migrate.sh prod

# Solo migraciones (comando directo)
go run migrations/cmd/migrate/main.go
```

### Ejecutar Aplicación
```bash
# Desarrollo
./migrations/scripts/dev.sh

# Producción
./migrations/scripts/prod.sh
```

## 🗄️ Estructura de Base de Datos

El proyecto usa **GORM** para migraciones automáticas. Las tablas se crean automáticamente:

- `users` - Usuarios del sistema
- `auth_users` - Usuarios de autenticación
- `sessions` - Sesiones de usuario
- `password_resets` - Tokens de restablecimiento
- `email_verifications` - Verificaciones de email

## 🚨 Solución de Problemas

### Error de Conexión a Base de Datos
```bash
# Verificar que PostgreSQL esté corriendo
docker ps | grep postgres

# Reiniciar PostgreSQL
make docker-down
make docker-up
```

### Error de Migraciones
```bash
# Limpiar y recrear
make docker-down
make docker-up
make migrate-dev
```

### Puerto en Uso
```bash
# Cambiar puerto en .env.dev
PORT=8081
```

## 📝 Notas Importantes

1. **Primera vez**: Usa `make quick-start` para configuración completa
2. **Desarrollo**: Los datos se persisten en volúmenes de Docker
3. **Producción**: Configura credenciales reales en `.env.prod`
4. **Migraciones**: Se ejecutan automáticamente al iniciar la aplicación
5. **Logs**: Revisa los logs para debugging

## 🔄 Flujo de Desarrollo

1. **Configurar entorno**: Crear archivos `.env.dev` y `.env.prod`
2. **Iniciar base de datos**: `make docker-up`
3. **Ejecutar migraciones**: `make migrate-dev`
4. **Iniciar aplicación**: `make dev`
5. **Desarrollar**: Hacer cambios en el código
6. **Probar**: Los cambios se reflejan automáticamente

## 🎯 Próximos Pasos

- Configurar CI/CD
- Implementar tests automatizados
- Configurar monitoreo
- Optimizar para producción
