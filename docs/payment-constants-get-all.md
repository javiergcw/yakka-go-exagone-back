# Obtener Todas las Constantes de Pago

## Descripción
Este endpoint permite obtener todas las constantes de pago del sistema, con la opción de filtrar solo las constantes activas.

## Endpoint
```
GET /api/v1/payment-constants
```

## Headers
```
Authorization: Bearer <token>
```

## Query Parameters
- `active_only` (boolean, opcional): Si se especifica como `true`, solo retorna las constantes activas

## Response

### Éxito (200 OK)
```json
{
  "constants": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "GST",
      "value": 10,
      "description": "Goods and Services Tax percentage",
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": "456e7890-e89b-12d3-a456-426614174001",
      "name": "WAGE_HOURLY",
      "value": 25,
      "description": "Default hourly wage rate",
      "is_active": true,
      "created_at": "2024-01-15T10:35:00Z",
      "updated_at": "2024-01-15T10:35:00Z"
    }
  ],
  "message": "Payment constants retrieved successfully"
}
```

### Error (500 Internal Server Error)
```json
{
  "error": "Failed to get payment constants"
}
```

## Ejemplos de Uso

### Obtener todas las constantes
```bash
curl -X GET http://localhost:8080/api/v1/payment-constants \
  -H "Authorization: Bearer <token>"
```

### Obtener solo constantes activas
```bash
curl -X GET "http://localhost:8080/api/v1/payment-constants?active_only=true" \
  -H "Authorization: Bearer <token>"
```

## Notas
- Si no hay constantes en el sistema, se retorna un array vacío
- El parámetro `active_only` es opcional y por defecto retorna todas las constantes
- Las constantes se ordenan por fecha de creación
