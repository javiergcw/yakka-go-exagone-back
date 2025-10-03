# Builder Assign Company

## Endpoint
```
POST /api/v1/builder/companies
```

## Descripción
Permite a un builder asignarse a una compañía existente. Si el builder ya tiene una compañía asignada, se actualiza a la nueva compañía. Si no tiene ninguna compañía asignada, se asigna normalmente.

## Autenticación
- **Tipo**: Bearer Token
- **Middleware**: `BuilderMiddleware`
- **Requerido**: Sí

## Headers
```
Authorization: Bearer <token>
Content-Type: application/json
```

## Request Body

### Estructura
```json
{
  "company_id": "string (UUID)"
}
```

### Campos
| Campo | Tipo | Requerido | Validación | Descripción |
|-------|------|-----------|------------|-------------|
| `company_id` | string | ✅ | UUID válido | ID de la compañía a asignar |

### Ejemplo de Request
```json
{
  "company_id": "123e4567-e89b-12d3-a456-426614174000"
}
```

## Response

### Estructura de Respuesta Exitosa (200 OK)
```json
{
  "builder_profile": {
    "id": "string",
    "user_id": "string",
    "company_id": "string",
    "company": {
      "id": "string",
      "name": "string",
      "description": "string",
      "website": "string",
      "created_at": "string (ISO 8601)",
      "updated_at": "string (ISO 8601)"
    },
    "display_name": "string",
    "location": "string",
    "bio": "string",
    "avatar_url": "string",
    "created_at": "string (ISO 8601)",
    "updated_at": "string (ISO 8601)"
  },
  "message": "string"
}
```

### Ejemplo de Respuesta Exitosa
```json
{
  "builder_profile": {
    "id": "456e7890-e89b-12d3-a456-426614174001",
    "user_id": "789e0123-e89b-12d3-a456-426614174002",
    "company_id": "123e4567-e89b-12d3-a456-426614174000",
    "company": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "TechCorp Solutions",
      "description": "Empresa líder en tecnología",
      "website": "https://techcorp.com",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    "display_name": "Juan Pérez",
    "location": "Madrid, España",
    "bio": "Desarrollador con 5 años de experiencia",
    "avatar_url": "https://example.com/avatar.jpg",
    "created_at": "2024-01-10T09:00:00Z",
    "updated_at": "2024-01-20T14:30:00Z"
  },
  "message": "Builder successfully assigned to company"
}
```

## Códigos de Error

### 400 Bad Request
```json
{
  "error": "Invalid request data",
  "message": "Validation failed",
  "details": [
    {
      "field": "company_id",
      "message": "company_id is required"
    }
  ]
}
```

### 401 Unauthorized
```json
{
  "error": "Unauthorized",
  "message": "Invalid or missing authentication token"
}
```

### 403 Forbidden
```json
{
  "error": "Forbidden",
  "message": "User is not a builder"
}
```

### 404 Not Found
```json
{
  "error": "Not Found",
  "message": "Company not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal Server Error",
  "message": "An unexpected error occurred"
}
```

## Ejemplos de Uso

### cURL
```bash
curl -X POST "http://localhost:8080/api/v1/builder/companies" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "company_id": "123e4567-e89b-12d3-a456-426614174000"
  }'
```

### JavaScript (Fetch)
```javascript
const response = await fetch('http://localhost:8080/api/v1/builder/companies', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    company_id: '123e4567-e89b-12d3-a456-426614174000'
  })
});

const data = await response.json();
console.log(data);
```

### Python (Requests)
```python
import requests

url = "http://localhost:8080/api/v1/builder/companies"
headers = {
    "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "Content-Type": "application/json"
}
data = {
    "company_id": "123e4567-e89b-12d3-a456-426614174000"
}

response = requests.post(url, headers=headers, json=data)
print(response.json())
```

## Comportamiento Especial

### Asignación de Nueva Compañía
- Si el builder no tiene ninguna compañía asignada, se asigna la nueva compañía
- Se crea o actualiza el perfil de builder con el `company_id`

### Actualización de Compañía Existente
- Si el builder ya tiene una compañía asignada, se actualiza a la nueva compañía
- El `company_id` anterior se reemplaza por el nuevo

### Validaciones
- El `company_id` debe ser un UUID válido
- La compañía debe existir en la base de datos
- El usuario debe ser un builder autenticado

## Notas Importantes

1. **Middleware Requerido**: Este endpoint requiere `BuilderMiddleware`, asegurando que solo builders autenticados puedan usarlo
2. **Actualización Automática**: Si el builder ya tiene una compañía, se actualiza automáticamente
3. **Relación Opcional**: El `company_id` es opcional en el perfil de builder
4. **Datos Completos**: La respuesta incluye tanto el `company_id` como los datos completos de la compañía

## Endpoints Relacionados

- `GET /api/v1/companies` - Obtener todas las compañías disponibles
- `POST /api/v1/companies` - Crear nueva compañía
- `GET /api/v1/profiles/builder` - Obtener perfil de builder actual
- `POST /api/v1/profiles/builder` - Crear/actualizar perfil de builder
