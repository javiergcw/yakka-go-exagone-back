# Licenses API Documentation

Esta documentación describe el endpoint para acceder a las licencias del sistema Yakka.

## Autenticación

El endpoint requiere el header `X-License` con una licencia válida.

**Header requerido:**
```
X-License: YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B
```

## Endpoint

### GET /api/v1/licenses

Obtiene todas las licencias disponibles en el sistema.

**Headers:**
```
X-License: YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B
Content-Type: application/json
```

**Ejemplo de Request:**
```bash
curl -X GET "https://api.yakka.com/api/v1/licenses" \
  -H "X-License: YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B" \
  -H "Content-Type: application/json"
```

**Ejemplo de Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid-here",
      "name": "Driving Licence",
      "description": "Official license to operate motor vehicles legally",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid-here",
      "name": "Work with Children",
      "description": "Certification required to work with children in various settings",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid-here",
      "name": "SIS20122 – Certificate II in Sport and Recreation",
      "description": "Entry-level qualification for sport and recreation industry",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid-here",
      "name": "SIS20321 – Certificate II in Sport Coaching",
      "description": "Foundation qualification for sport coaching roles",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid-here",
      "name": "SIS30122 – Certificate III in Sport, Aquatics and Recreation",
      "description": "Intermediate qualification in sport, aquatics and recreation management",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid-here",
      "name": "SIS30521 – Certificate III in Sport Coaching",
      "description": "Advanced coaching qualification for sport professionals",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid-here",
      "name": "SIS40321 – Certificate IV in Sport Coaching",
      "description": "Senior coaching qualification for leadership roles",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid-here",
      "name": "SIS50321 – Diploma of Sport",
      "description": "Comprehensive diploma program in sport management and coaching",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid-here",
      "name": "SIS50122 – Diploma of Sport, Aquatics and Recreation Management",
      "description": "Specialized diploma in sport, aquatics and recreation management",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid-here",
      "name": "Diploma of Sport & Exercise Science (provider specific)",
      "description": "Provider-specific diploma in sport and exercise science",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid-here",
      "name": "Bachelor Degree in Sport Science",
      "description": "Undergraduate degree in sport science and related fields",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid-here",
      "name": "Bachelor Honours",
      "description": "Honours degree in sport science with research component",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid-here",
      "name": "Master's Degree in Sport",
      "description": "Postgraduate degree in sport science and management",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid-here",
      "name": "Doctoral Degree (PhD) in Sport",
      "description": "Highest academic qualification in sport science and research",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "message": "Licenses retrieved successfully"
}
```

## Códigos de Error

### 401 Unauthorized
```json
{
  "success": false,
  "error": "License header required",
  "message": "X-License header is missing or invalid"
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "error": "Failed to retrieve licenses",
  "message": "Internal server error occurred while fetching licenses"
}
```

## Notas Importantes

1. **Licencia Requerida**: El endpoint requiere el header `X-License` con la licencia válida.
2. **Rate Limiting**: El endpoint está sujeto a rate limiting (100 requests por minuto).
3. **CORS**: El endpoint soporta CORS para requests desde dominios autorizados.
4. **Formato de Fechas**: Todas las fechas están en formato ISO 8601 (RFC3339).
5. **UUIDs**: Todos los IDs son UUIDs v4.

## Tipos de Licencias Disponibles

### Licencias Básicas
- **Driving Licence**: Licencia de conducir oficial
- **Work with Children**: Certificación para trabajar con niños

### Certificados SIS (Sport Industry Skills)
- **SIS20122**: Certificate II in Sport and Recreation
- **SIS20321**: Certificate II in Sport Coaching
- **SIS30122**: Certificate III in Sport, Aquatics and Recreation
- **SIS30521**: Certificate III in Sport Coaching
- **SIS40321**: Certificate IV in Sport Coaching

### Diplomas
- **SIS50321**: Diploma of Sport
- **SIS50122**: Diploma of Sport, Aquatics and Recreation Management
- **Diploma of Sport & Exercise Science**: Específico del proveedor

### Grados Universitarios
- **Bachelor Degree in Sport Science**: Grado en Ciencias del Deporte
- **Bachelor Honours**: Grado con Honores
- **Master's Degree in Sport**: Máster en Deporte
- **Doctoral Degree (PhD) in Sport**: Doctorado en Deporte
