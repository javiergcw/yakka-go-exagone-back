# Skills Complete API Documentation

Esta documentación describe el endpoint para acceder a las categorías de habilidades con sus subcategorías en una sola respuesta del sistema Yakka.

## Autenticación

El endpoint requiere el header `X-License` con una licencia válida.

**Header requerido:**
```
X-License: YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B
```

## Endpoint

### GET /api/v1/skills

Obtiene todas las categorías de habilidades con sus subcategorías en una sola respuesta. Este endpoint es ideal para obtener la estructura completa de skills de una vez.

**Headers:**
```
X-License: YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B
Content-Type: application/json
```

**Ejemplo de Request:**
```bash
curl -X GET "https://api.yakka.com/api/v1/skills" \
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
      "updated_at": "2024-01-01T00:00:00Z",
      "subcategories": [
        {
          "id": "uuid-here",
          "name": "Soccer"
        },
        {
          "id": "uuid-here",
          "name": "Rugby Union"
        },
        {
          "id": "uuid-here",
          "name": "Rugby League"
        },
        {
          "id": "uuid-here",
          "name": "Australian Rules Football (AFL)"
        },
        {
          "id": "uuid-here",
          "name": "Basketball"
        },
        {
          "id": "uuid-here",
          "name": "Netball"
        },
        {
          "id": "uuid-here",
          "name": "Tennis"
        },
        {
          "id": "uuid-here",
          "name": "Swimming"
        },
        {
          "id": "uuid-here",
          "name": "Athletics"
        },
        {
          "id": "uuid-here",
          "name": "Cricket"
        },
        {
          "id": "uuid-here",
          "name": "Hockey"
        },
        {
          "id": "uuid-here",
          "name": "Volleyball"
        },
        {
          "id": "uuid-here",
          "name": "Gymnastics"
        },
        {
          "id": "uuid-here",
          "name": "Baseball"
        },
        {
          "id": "uuid-here",
          "name": "Softball"
        },
        {
          "id": "uuid-here",
          "name": "Surfing"
        },
        {
          "id": "uuid-here",
          "name": "Cycling"
        },
        {
          "id": "uuid-here",
          "name": "Boxing"
        },
        {
          "id": "uuid-here",
          "name": "Martial Arts"
        },
        {
          "id": "uuid-here",
          "name": "Golf"
        },
        {
          "id": "uuid-here",
          "name": "Other"
        }
      ]
    },
    {
      "id": "uuid-here",
      "name": "Personal Trainer (PT)",
      "description": "Personal training and fitness coaching",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "subcategories": [
        {
          "id": "uuid-here",
          "name": "Fitness"
        },
        {
          "id": "uuid-here",
          "name": "Gym"
        },
        {
          "id": "uuid-here",
          "name": "CrossFit"
        },
        {
          "id": "uuid-here",
          "name": "Functional Training"
        },
        {
          "id": "uuid-here",
          "name": "Athletics"
        },
        {
          "id": "uuid-here",
          "name": "Swimming"
        },
        {
          "id": "uuid-here",
          "name": "Cycling"
        },
        {
          "id": "uuid-here",
          "name": "Boxing"
        },
        {
          "id": "uuid-here",
          "name": "Martial Arts"
        },
        {
          "id": "uuid-here",
          "name": "Yoga"
        },
        {
          "id": "uuid-here",
          "name": "Pilates"
        },
        {
          "id": "uuid-here",
          "name": "Other"
        }
      ]
    },
    {
      "id": "uuid-here",
      "name": "Referee",
      "description": "Sports officiating and refereeing roles",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "subcategories": [
        {
          "id": "uuid-here",
          "name": "Soccer"
        },
        {
          "id": "uuid-here",
          "name": "Basketball"
        },
        {
          "id": "uuid-here",
          "name": "Tennis"
        },
        {
          "id": "uuid-here",
          "name": "Swimming"
        },
        {
          "id": "uuid-here",
          "name": "Athletics"
        },
        {
          "id": "uuid-here",
          "name": "Other"
        }
      ]
    }
  ],
  "message": "Skills with categories and subcategories retrieved successfully"
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
  "error": "Failed to retrieve skills",
  "message": "Internal server error occurred while fetching skills"
}
```

## Notas Importantes

1. **Licencia Requerida**: El endpoint requiere el header `X-License` con la licencia válida.
2. **Rate Limiting**: El endpoint está sujeto a rate limiting (100 requests por minuto).
3. **CORS**: El endpoint soporta CORS para requests desde dominios autorizados.
4. **Formato de Fechas**: Todas las fechas están en formato ISO 8601 (RFC3339).
5. **UUIDs**: Todos los IDs son UUIDs v4.
6. **Eficiencia**: Este endpoint es más eficiente que hacer múltiples requests separados.

## Ventajas del Endpoint Completo

### Eficiencia
- **Un solo request**: Obtiene toda la estructura de skills de una vez
- **Menos latencia**: Evita múltiples round-trips al servidor
- **Mejor rendimiento**: Ideal para cargar datos iniciales de la aplicación

### Completitud
- **Datos completos**: Incluye categorías con todas sus subcategorías
- **IDs disponibles**: Cada subcategoría incluye su ID para futuras operaciones
- **Estructura jerárquica**: Mantiene la relación padre-hijo

## Casos de Uso

### Para Aplicaciones Frontend
- **Carga inicial**: Obtener toda la estructura de skills al cargar la aplicación
- **Formularios**: Poblar dropdowns de categorías y subcategorías
- **Filtros**: Implementar filtros jerárquicos de skills

### Para Trabajadores (Labour)
- **Perfil completo**: Seleccionar categorías y subcategorías específicas
- **Especialización**: Mostrar áreas de expertise detalladas
- **Búsqueda**: Filtrar por especializaciones específicas

### Para Empleadores (Builders)
- **Publicación de trabajos**: Especificar categorías y subcategorías requeridas
- **Búsqueda de talento**: Filtrar candidatos por especializaciones
- **Matching**: Conectar trabajadores con empleos apropiados

## Comparación con Otros Endpoints

### Endpoint Completo vs. Separados

**Endpoint Completo (Recomendado):**
```bash
# Un solo request
curl -X GET "https://api.yakka.com/api/v1/skills" \
  -H "X-License: YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B"
```

**Endpoints Separados:**
```bash
# Múltiples requests
curl -X GET "https://api.yakka.com/api/v1/skill-categories" \
  -H "X-License: YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B"

# Para cada categoría
curl -X GET "https://api.yakka.com/api/v1/skill-categories/{categoryId}/subcategories" \
  -H "X-License: YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B"
```

## Integración con Perfiles

Las skills completas se integran con:
- **Labour Profiles**: Para definir habilidades y especializaciones
- **Job Postings**: Para especificar requisitos de skills
- **Matching System**: Para conectar trabajadores con empleos
- **Search System**: Para filtrar por skills específicas

## Mejores Prácticas

1. **Carga inicial**: Usar este endpoint para cargar datos iniciales de la aplicación
2. **Caché**: Implementar caché para evitar requests repetidos
3. **Lazy loading**: Cargar subcategorías solo cuando sea necesario
4. **Filtrado**: Implementar filtros en el frontend para mejor UX
