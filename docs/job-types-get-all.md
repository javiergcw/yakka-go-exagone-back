# Get All Job Types

Obtiene todos los tipos de trabajo disponibles en el sistema.

## Endpoint

```
GET /api/v1/job-types
```

## Headers

```
Authorization: Bearer <token>
Content-Type: application/json
```

## Query Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `active_only` | boolean | No | Si es `true`, solo retorna tipos de trabajo activos. Por defecto retorna todos. |

## Response

### Success Response (200 OK)

```json
{
  "types": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "Casual Job",
      "description": "Casual employment with flexible hours",
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": "123e4567-e89b-12d3-a456-426614174001",
      "name": "Part Time",
      "description": "Part-time employment",
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": "123e4567-e89b-12d3-a456-426614174002",
      "name": "Full Time",
      "description": "Full-time employment",
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": "123e4567-e89b-12d3-a456-426614174003",
      "name": "Farms Job",
      "description": "Agricultural and farming work",
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": "123e4567-e89b-12d3-a456-426614174004",
      "name": "Mining Job",
      "description": "Mining industry work",
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": "123e4567-e89b-12d3-a456-426614174005",
      "name": "FIFO",
      "description": "Fly In Fly Out work arrangements",
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": "123e4567-e89b-12d3-a456-426614174006",
      "name": "Seasonal Job",
      "description": "Seasonal employment",
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": "123e4567-e89b-12d3-a456-426614174007",
      "name": "W&H Visa",
      "description": "Working Holiday Visa jobs",
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": "123e4567-e89b-12d3-a456-426614174008",
      "name": "Other",
      "description": "Other types of employment",
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ],
  "message": "Job types retrieved successfully"
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
  "message": "Failed to get job types"
}
```

## Examples

### Get All Job Types

```bash
curl -X GET "http://localhost:8081/api/v1/job-types" \
  -H "Authorization: Bearer your-token-here" \
  -H "Content-Type: application/json"
```

### Get Only Active Job Types

```bash
curl -X GET "http://localhost:8081/api/v1/job-types?active_only=true" \
  -H "Authorization: Bearer your-token-here" \
  -H "Content-Type: application/json"
```

## Notes

- Este endpoint requiere autenticación con middleware de licencia
- Los tipos de trabajo incluyen diferentes modalidades de empleo como casual, part-time, full-time, FIFO, etc.
- El parámetro `active_only=true` es útil para obtener solo los tipos de trabajo que están actualmente disponibles
- Todos los timestamps están en formato ISO 8601 UTC
