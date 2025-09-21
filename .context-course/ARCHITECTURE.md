# Arquitectura del Proyecto

## Estructura General
```
yakka-go-exagone-back/
├── config/                    # Configuración global
├── database/                  # Configuración de base de datos
├── http/                     # HTTP handlers (legacy)
├── internal/                 # Código interno de la aplicación
│   ├── features/             # Módulos de funcionalidades
│   │   └── users/            # Módulo de usuarios
│   │       ├── delivery/     # Capa de entrega (REST)
│   │       ├── entity/       # Entidades y repositorios
│   │       ├── models/       # Modelos de datos
│   │       ├── payload/      # DTOs y payloads
│   │       └── usecase/      # Lógica de negocio
│   ├── infrastructure/       # Infraestructura
│   │   ├── config/          # Configuración interna
│   │   ├── database/        # Base de datos
│   │   └── http/            # HTTP server y middleware
│   └── shared/              # Código compartido
│       ├── constants/       # Constantes
│       ├── errors/          # Manejo de errores
│       ├── response/         # Respuestas HTTP
│       └── validation/       # Validaciones
└── main.go                  # Punto de entrada
```

## Patrón de Arquitectura
- **Clean Architecture** con separación de capas
- **Feature-based** organization
- **Dependency Injection** para handlers

## Imports Importantes
```go
// Handlers de usuarios
"github.com/yakka-backend/internal/features/users/delivery/rest"

// Infraestructura
"github.com/yakka-backend/internal/infrastructure/http/middleware"
"github.com/yakka-backend/internal/infrastructure/config"

// Shared
"github.com/yakka-backend/internal/shared/response"
"github.com/yakka-backend/internal/shared/errors"
```

## Convenciones
- **Widgets globales**: `@widgets/` (múltiples módulos)
- **Widgets específicos**: `{module}/widgets/` (por módulo)
- **No usar**: `components/` (solo `widgets/`)

## Stack Tecnológico
- **Go 1.24.2**
- **Gorilla Mux** (router)
- **GORM** (ORM)
- **PostgreSQL** (base de datos)
- **Docker** (contenedores)
