# Eliminar Jobsite

## Endpoint
**DELETE** `/api/v1/jobsites/{id}`

## Descripción
Elimina un jobsite existente. Solo funciona si el jobsite pertenece al builder autenticado.

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
DELETE /api/v1/jobsites/123e4567-e89b-12d3-a456-426614174000
```

## Body
No requiere body

## Response Success (200)
```json
{
  "message": "Jobsite deleted successfully"
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

## Response Error (500)
```json
{
  "error": "Failed to delete jobsite: database error"
}
```
