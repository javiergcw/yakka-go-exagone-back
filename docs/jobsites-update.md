# Actualizar Jobsite

## Endpoint
**PUT** `/api/v1/jobsites/{id}`

## Descripción
Actualiza un jobsite existente. Solo funciona si el jobsite pertenece al builder autenticado.

## Headers
```
Content-Type: application/json
Authorization: Bearer <jwt_token>
```

## Parámetros de URL
| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `id` | UUID | ID único del jobsite |

## Ejemplo de URL
```
PUT /api/v1/jobsites/123e4567-e89b-12d3-a456-426614174000
```

## Body (JSON)
```json
{
  "address": "123 Main Street, Sydney NSW 2000 - Updated",
  "city": "Sydney",
  "suburb": "CBD",
  "description": "Nueva construcción de edificio residencial - Fase 2",
  "latitude": -33.8688,
  "longitude": 151.2093,
  "phone": "+61 2 1234 5678"
}
```

## Parámetros del Body (Todos opcionales)
| Campo | Tipo | Descripción |
|-------|------|-------------|
| `address` | string | Dirección actualizada (10-500 caracteres) |
| `city` | string | Ciudad actualizada (2-120 caracteres) |
| `suburb` | string | Suburbio actualizado (2-120 caracteres) |
| `description` | string | Descripción actualizada (máx 1000 caracteres) |
| `latitude` | number | Latitud actualizada (-90 a 90) |
| `longitude` | number | Longitud actualizada (-180 a 180) |
| `phone` | string | Teléfono actualizado (10-32 caracteres) |

## Response Success (200)
```json
{
  "jobsite": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "builder_id": "987fcdeb-51a2-43d1-9f12-345678901234",
    "address": "123 Main Street, Sydney NSW 2000 - Updated",
    "city": "Sydney",
    "suburb": "CBD",
    "description": "Nueva construcción de edificio residencial - Fase 2",
    "latitude": -33.8688,
    "longitude": 151.2093,
    "phone": "+61 2 1234 5678",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T15:45:00Z"
  },
  "message": "Jobsite updated successfully"
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

## Response Error (403)
```json
{
  "error": "Access denied: This jobsite does not belong to you"
}
```

## Response Error (404)
```json
{
  "error": "jobsite not found"
}
```
