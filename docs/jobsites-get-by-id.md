# Obtener Jobsite por ID

## Endpoint
**GET** `/api/v1/jobsites/{id}`

## Descripción
Obtiene un jobsite específico por su ID. Solo funciona si el jobsite pertenece al builder autenticado.

## Headers
```
Authorization: Bearer <jwt_token>
```

## Parámetros de URL
| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `id` | UUID | ID único del jobsite |

## Ejemplo de URL
```
GET /api/v1/jobsites/123e4567-e89b-12d3-a456-426614174000
```

## Response Success (200)
```json
{
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
}
```

## Response Error (400)
```json
{
  "error": "Invalid jobsite ID"
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
