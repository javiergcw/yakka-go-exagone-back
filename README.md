# Yakka Go Exagone Backend

Un backend moderno en Go con API REST para gestión de usuarios.

## 🚀 Características

- **API REST** completa con endpoints CRUD
- **CORS** habilitado para desarrollo frontend
- **Middleware** de logging y manejo de errores
- **Estructura modular** y fácil de extender
- **Validación** de datos de entrada
- **Respuestas JSON** estandarizadas

## 📋 Endpoints Disponibles

### Health Check
- `GET /health` - Verificar estado del servidor

### Usuarios
- `GET /api/v1/users` - Obtener todos los usuarios
- `POST /api/v1/users` - Crear nuevo usuario
- `GET /api/v1/users/{id}` - Obtener usuario por ID
- `PUT /api/v1/users/{id}` - Actualizar usuario
- `DELETE /api/v1/users/{id}` - Eliminar usuario

## 🛠️ Instalación y Uso

### Prerrequisitos
- Go 1.19 o superior
- Git

### Instalación

1. **Clonar el repositorio:**
```bash
git clone <tu-repositorio>
cd yakka-go-exagone-back
```

2. **Instalar dependencias:**
```bash
go mod tidy
```

3. **Ejecutar el servidor:**
```bash
go run main.go
```

El servidor se ejecutará en `http://localhost:8080`

### Variables de Entorno

Puedes configurar el puerto usando la variable de entorno `PORT`:
```bash
export PORT=3000
go run main.go
```

## 📝 Ejemplos de Uso

### Crear un usuario
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Juan Pérez",
    "email": "juan@example.com"
  }'
```

### Obtener todos los usuarios
```bash
curl http://localhost:8080/api/v1/users
```

### Obtener un usuario específico
```bash
curl http://localhost:8080/api/v1/users/1
```

### Actualizar un usuario
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Juan Carlos Pérez",
    "email": "juancarlos@example.com"
  }'
```

### Eliminar un usuario
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

## 🏗️ Estructura del Proyecto

```
yakka-go-exagone-back/
├── main.go                    # Archivo principal con servidor HTTP
├── go.mod                     # Dependencias del módulo Go
├── go.sum                     # Checksums de dependencias
├── docker-compose.yml         # Configuración de Docker Compose
├── Dockerfile                 # Configuración de Docker
├── env.example                # Variables de entorno de ejemplo
├── .gitignore                 # Archivos ignorados por Git
├── internal/
│   ├── config/                # Configuración de la aplicación
│   │   └── config.go
│   ├── database/              # Configuración de base de datos
│   │   └── database.go
│   ├── http/                  # Router HTTP
│   │   └── router.go
│   └── users/                 # Feature: Gestión de usuarios
│       ├── interface.go       # Interfaces del dominio
│       ├── models/            # Modelos de datos
│       │   └── user.go
│       ├── payload/           # DTOs de request/response
│       │   └── user.go
│       ├── entity/database/   # Implementación de repositorio
│       │   └── user_repo.go
│       ├── usecase/           # Lógica de negocio
│       │   └── user_operations.go
│       └── delivery/rest/     # Handlers HTTP
│           ├── endpoint.go
│           └── methods.go
└── README.md                  # Documentación del proyecto
```

### 🏛️ Arquitectura Hexagonal por Feature

El proyecto sigue una **arquitectura hexagonal** organizada por **features**:

- **`models/`**: Entidades del dominio
- **`payload/`**: DTOs para requests y responses
- **`entity/database/`**: Implementación de repositorios (infraestructura)
- **`usecase/`**: Casos de uso y lógica de negocio
- **`delivery/rest/`**: Handlers HTTP (adaptadores)
- **`interface.go`**: Definición de interfaces del dominio

## 🔧 Desarrollo

### Agregar nuevas rutas
1. Define el handler en `main.go`
2. Registra la ruta en la función `main()`
3. Asegúrate de manejar errores apropiadamente

### Agregar middleware
Puedes agregar middleware personalizado usando la librería Gorilla Mux:

```go
r.Use(loggingMiddleware)
r.Use(authMiddleware)
```

## 🚀 Despliegue

### Docker (Opcional)
```dockerfile
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

### Compilar para producción
```bash
go build -o yakka-backend main.go
./yakka-backend
```

## 📊 Respuesta de la API

Todas las respuestas siguen el formato estándar:

```json
{
  "success": true,
  "message": "Operación exitosa",
  "data": {
    // Datos específicos de la respuesta
  }
}
```

## 🤝 Contribuir

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo `LICENSE` para más detalles.

## 🆘 Soporte

Si tienes alguna pregunta o problema, por favor abre un issue en el repositorio.
