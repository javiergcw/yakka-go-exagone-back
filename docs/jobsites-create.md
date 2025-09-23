# Crear Jobsite

## Endpoint
**POST** `/api/v1/jobsites`

## Descripción
Crea un nuevo jobsite para el builder autenticado. El `builder_id` se toma automáticamente del JWT del usuario.

## Headers
```
Content-Type: application/json
Authorization: Bearer <jwt_token>
```

## Body (JSON)
```json
{
  "address": "123 Main Street, Sydney NSW 2000",
  "city": "Sydney",
  "suburb": "CBD",
  "description": "Nueva construcción de edificio residencial",
  "latitude": -33.8688,
  "longitude": 151.2093,
  "phone": "+61 2 1234 5678"
}
```

## Parámetros del Body
| Campo | Tipo | Requerido | Descripción |
|-------|------|-----------|-------------|
| `address` | string | ✅ | Dirección completa del jobsite (10-500 caracteres) |
| `city` | string | ✅ | Ciudad (2-120 caracteres) |
| `suburb` | string | ❌ | Suburbio o distrito (2-120 caracteres) |
| `description` | string | ❌ | Descripción adicional (máx 1000 caracteres) |
| `latitude` | number | ✅ | Latitud (-90 a 90) |
| `longitude` | number | ✅ | Longitud (-180 a 180) |
| `phone` | string | ❌ | Teléfono de contacto (10-32 caracteres) |

## Response Success (201)
```json
{
  "jobsite": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "builder_id": "987fcdeb-51a2-43d1-9f12-345678901234",
    "address": "123 Main Street, Sydney NSW 2000",
    "city": "Sydney",
    "suburb": "CBD",
    "description": "Nueva construcción de edificio residencial",
    "latitude": -33.8688,
    "longitude": 151.2093,
    "phone": "+61 2 1234 5678",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  },
  "message": "Jobsite created successfully"
}
```

## Response Error (400)
```json
{
  "error": "Invalid request body"
}
```

## Response Error (401)
```json
{
  "error": "User not authenticated"
}
```

## Response Error (500)
```json
{
  "error": "Failed to create jobsite: database error"
}
```
