# Experience Levels API Documentation

Esta documentación describe el endpoint para acceder a los niveles de experiencia del sistema Yakka.

## Autenticación

El endpoint requiere el header `X-License` con una licencia válida.

**Header requerido:**
```
X-License: YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B
```

## Endpoint

### GET /api/v1/experience-levels

Obtiene todos los niveles de experiencia disponibles.

**Headers:**
```
X-License: YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B
Content-Type: application/json
```

**Ejemplo de Request:**
```bash
curl -X GET "https://api.yakka.com/api/v1/experience-levels" \
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
      "name": "Less than 6 months",
      "description": "Experience level for workers with less than 6 months of experience",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid-here",
      "name": "6-12 months",
      "description": "Experience level for workers with 6 to 12 months of experience",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid-here",
      "name": "1-2 years",
      "description": "Experience level for workers with 1 to 2 years of experience",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid-here",
      "name": "2-5 years",
      "description": "Experience level for workers with 2 to 5 years of experience",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid-here",
      "name": "More than 5 years",
      "description": "Experience level for workers with more than 5 years of experience",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "message": "Experience levels retrieved successfully"
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
  "error": "Failed to retrieve experience levels",
  "message": "Internal server error occurred while fetching experience levels"
}
```

## Notas Importantes

1. **Licencia Requerida**: El endpoint requiere el header `X-License` con la licencia válida.
2. **Rate Limiting**: El endpoint está sujeto a rate limiting (100 requests por minuto).
3. **CORS**: El endpoint soporta CORS para requests desde dominios autorizados.
4. **Formato de Fechas**: Todas las fechas están en formato ISO 8601 (RFC3339).
5. **UUIDs**: Todos los IDs son UUIDs v4.

## Niveles de Experiencia Disponibles

### Niveles Básicos
- **Less than 6 months**: Para trabajadores con menos de 6 meses de experiencia
- **6-12 months**: Para trabajadores con 6 a 12 meses de experiencia

### Niveles Intermedios
- **1-2 years**: Para trabajadores con 1 a 2 años de experiencia
- **2-5 years**: Para trabajadores con 2 a 5 años de experiencia

### Niveles Avanzados
- **More than 5 years**: Para trabajadores con más de 5 años de experiencia

## Casos de Uso

### Para Trabajadores (Labour)
Los niveles de experiencia ayudan a los trabajadores a:
- Definir su nivel de experiencia en el perfil
- Mostrar su experiencia a empleadores potenciales
- Filtrar oportunidades según su nivel

### Para Empleadores (Builders)
Los niveles de experiencia permiten a los empleadores:
- Filtrar candidatos por nivel de experiencia
- Publicar trabajos con requisitos de experiencia específicos
- Evaluar la idoneidad de los candidatos

## Integración con Perfiles

Los niveles de experiencia se integran con:
- **Labour Profiles**: Para definir la experiencia del trabajador
- **Job Postings**: Para especificar requisitos de experiencia
- **Matching System**: Para conectar trabajadores con empleos apropiados
