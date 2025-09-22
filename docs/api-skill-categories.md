# Skill Categories API Documentation

Esta documentación describe el endpoint para acceder a las categorías de habilidades del sistema Yakka.

## Autenticación

El endpoint requiere el header `X-License` con una licencia válida.

**Header requerido:**
```
X-License: YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B
```

## Endpoint

### GET /api/v1/skill-categories

Obtiene todas las categorías de habilidades disponibles.

**Headers:**
```
X-License: YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B
Content-Type: application/json
```

**Ejemplo de Request:**
```bash
curl -X GET "https://api.yakka.com/api/v1/skill-categories" \
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
      "name": "Coach",
      "description": "Professional coaching roles in various sports",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid-here",
      "name": "Personal Trainer (PT)",
      "description": "Personal training and fitness coaching",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid-here",
      "name": "Referee",
      "description": "Sports officiating and refereeing roles",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid-here",
      "name": "Sports Administrator",
      "description": "Administrative roles in sports organizations",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid-here",
      "name": "Sports Medicine",
      "description": "Medical and health support roles in sports",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid-here",
      "name": "Sports Science",
      "description": "Scientific and research roles in sports",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid-here",
      "name": "Event Management",
      "description": "Event planning and management roles",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid-here",
      "name": "Sports Media",
      "description": "Media and communication roles in sports",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "message": "Skill categories retrieved successfully"
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
  "error": "Failed to retrieve skill categories",
  "message": "Internal server error occurred while fetching skill categories"
}
```

## Notas Importantes

1. **Licencia Requerida**: El endpoint requiere el header `X-License` con la licencia válida.
2. **Rate Limiting**: El endpoint está sujeto a rate limiting (100 requests por minuto).
3. **CORS**: El endpoint soporta CORS para requests desde dominios autorizados.
4. **Formato de Fechas**: Todas las fechas están en formato ISO 8601 (RFC3339).
5. **UUIDs**: Todos los IDs son UUIDs v4.

## Categorías de Habilidades Disponibles

### Coaching y Entrenamiento
- **Coach**: Roles profesionales de entrenamiento en varios deportes
- **Personal Trainer (PT)**: Entrenamiento personal y coaching de fitness

### Oficiales y Administración
- **Referee**: Roles de arbitraje y oficialización deportiva
- **Sports Administrator**: Roles administrativos en organizaciones deportivas

### Medicina y Ciencia
- **Sports Medicine**: Roles médicos y de apoyo en salud deportiva
- **Sports Science**: Roles científicos y de investigación en deportes

### Gestión y Medios
- **Event Management**: Roles de planificación y gestión de eventos
- **Sports Media**: Roles de medios y comunicación en deportes

## Casos de Uso

### Para Trabajadores (Labour)
Las categorías de habilidades ayudan a los trabajadores a:
- Definir su área de especialización
- Mostrar sus habilidades a empleadores
- Filtrar oportunidades por categoría

### Para Empleadores (Builders)
Las categorías de habilidades permiten a los empleadores:
- Publicar trabajos con categorías específicas
- Filtrar candidatos por área de especialización
- Buscar talento en áreas específicas

## Relación con Subcategorías

Cada categoría tiene subcategorías específicas que se pueden obtener usando:
- **GET /api/v1/skill-subcategories?category_id={id}**: Filtrar por categoría
- **GET /api/v1/skill-categories/{categoryId}/subcategories**: Obtener subcategorías de una categoría específica
- **GET /api/v1/skills**: Obtener categorías con sus subcategorías en una sola respuesta

## Integración con Perfiles

Las categorías de habilidades se integran con:
- **Labour Profiles**: Para definir las habilidades del trabajador
- **Job Postings**: Para especificar categorías requeridas
- **Matching System**: Para conectar trabajadores con empleos apropiados
- **Skill Subcategories**: Para especificar habilidades más detalladas
