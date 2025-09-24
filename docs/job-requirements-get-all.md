# Get All Job Requirements

Obtiene todos los requisitos de trabajo disponibles en el sistema.

## Endpoint

```
GET /api/v1/job-requirements
```

## Headers

```
Authorization: Bearer <token>
Content-Type: application/json
```

## Query Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `active_only` | boolean | No | Si es `true`, solo retorna requisitos activos. Por defecto retorna todos. |

## Response

### Success Response (200 OK)

```json
{
  "requirements": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "White Card",
      "description": "Construction industry safety card",
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": "123e4567-e89b-12d3-a456-426614174001",
      "name": "First Aid",
      "description": "First aid certification required",
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": "123e4567-e89b-12d3-a456-426614174002",
      "name": "Driver License",
      "description": "Valid driver's license required",
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": "123e4567-e89b-12d3-a456-426614174003",
      "name": "Own Tools",
      "description": "Must provide own tools",
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": "123e4567-e89b-12d3-a456-426614174004",
      "name": "Safety Boots",
      "description": "Safety boots required",
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": "123e4567-e89b-12d3-a456-426614174005",
      "name": "Hard Hat",
      "description": "Hard hat required",
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": "123e4567-e89b-12d3-a456-426614174006",
      "name": "High Vis Vest",
      "description": "High visibility vest required",
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": "123e4567-e89b-12d3-a456-426614174007",
      "name": "Experience Required",
      "description": "Previous experience in the field required",
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ],
  "message": "Job requirements retrieved successfully"
}
```

### Error Response (401 Unauthorized)

```json
{
  "error": "Unauthorized",
  "message": "Invalid or missing license key"
}
```

### Error Response (500 Internal Server Error)

```json
{
  "error": "Internal Server Error",
  "message": "Failed to get job requirements"
}
```

## Examples

### Get All Job Requirements

```bash
curl -X GET "http://localhost:8081/api/v1/job-requirements" \
  -H "Authorization: Bearer your-token-here" \
  -H "Content-Type: application/json"
```

### Get Only Active Job Requirements

```bash
curl -X GET "http://localhost:8081/api/v1/job-requirements?active_only=true" \
  -H "Authorization: Bearer your-token-here" \
  -H "Content-Type: application/json"
```

## Notes

- Este endpoint requiere autenticación con middleware de licencia
- Los requisitos de trabajo incluyen elementos como certificaciones, licencias, equipos de seguridad, etc.
- El parámetro `active_only=true` es útil para obtener solo los requisitos que están actualmente disponibles
- Todos los timestamps están en formato ISO 8601 UTC
