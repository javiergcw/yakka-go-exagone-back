# Obtener Jobsites del Builder

## Endpoint
**GET** `/api/v1/jobsites`

## Descripción
Obtiene todos los jobsites del builder autenticado.

## Headers
```
Authorization: Bearer <jwt_token>
```

## Parámetros de Query
Ninguno

## Response Success (200)
```json
{
  "jobsites": [
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
    },
    {
      "id": "456e7890-e89b-12d3-a456-426614174001",
      "builder_id": "987fcdeb-51a2-43d1-9f12-345678901234",
      "address": "456 Queen Street, Melbourne VIC 3000",
      "city": "Melbourne",
      "suburb": "Central",
      "description": "Proyecto de renovación comercial",
      "latitude": -37.8136,
      "longitude": 144.9631,
      "phone": "+61 3 9876 5432",
      "created_at": "2024-01-16T14:20:00Z",
      "updated_at": "2024-01-16T14:20:00Z"
    }
  ],
  "total": 2
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
  "error": "Failed to retrieve jobsites: database error"
}
```
