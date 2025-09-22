# Skill Subcategories API Documentation

Esta documentación describe los endpoints para acceder a las subcategorías de habilidades del sistema Yakka.

## Autenticación

Los endpoints requieren el header `X-License` con una licencia válida.

**Header requerido:**
```
X-License: YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B
```

## Endpoints

### GET /api/v1/skill-subcategories

Obtiene todas las subcategorías de habilidades o las filtra por categoría.

**Headers:**
```
X-License: YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B
Content-Type: application/json
```

**Parámetros de Query (opcionales):**
- `category_id` (string): ID de la categoría para filtrar subcategorías

**Ejemplo de Request (todas las subcategorías):**
```bash
curl -X GET "https://api.yakka.com/api/v1/skill-subcategories" \
  -H "X-License: YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B" \
  -H "Content-Type: application/json"
```

**Ejemplo de Request (filtrar por categoría):**
```bash
curl -X GET "https://api.yakka.com/api/v1/skill-subcategories?category_id=uuid-here" \
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
      "category_id": "uuid-here",
      "name": "Soccer",
      "description": "Sport specialization for Coach",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "category": {
        "id": "uuid-here",
        "name": "Coach",
        "description": "Professional coaching roles in various sports"
      }
    },
    {
      "id": "uuid-here",
      "category_id": "uuid-here",
      "name": "Basketball",
      "description": "Sport specialization for Coach",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "category": {
        "id": "uuid-here",
        "name": "Coach",
        "description": "Professional coaching roles in various sports"
      }
    },
    {
      "id": "uuid-here",
      "category_id": "uuid-here",
      "name": "Fitness",
      "description": "Sport specialization for Personal Trainer (PT)",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "category": {
        "id": "uuid-here",
        "name": "Personal Trainer (PT)",
        "description": "Personal training and fitness coaching"
      }
    }
  ],
  "message": "Skill subcategories retrieved successfully"
}
```

### GET /api/v1/skill-categories/{categoryId}/subcategories

Obtiene las subcategorías para una categoría específica usando el ID de la categoría en la URL.

**Headers:**
```
X-License: YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B
Content-Type: application/json
```

**Parámetros de URL:**
- `categoryId` (string): ID de la categoría

**Ejemplo de Request:**
```bash
curl -X GET "https://api.yakka.com/api/v1/skill-categories/uuid-here/subcategories" \
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
      "category_id": "uuid-here",
      "name": "Soccer",
      "description": "Sport specialization for Coach",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "category": {
        "id": "uuid-here",
        "name": "Coach",
        "description": "Professional coaching roles in various sports"
      }
    },
    {
      "id": "uuid-here",
      "category_id": "uuid-here",
      "name": "Rugby Union",
      "description": "Sport specialization for Coach",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "category": {
        "id": "uuid-here",
        "name": "Coach",
        "description": "Professional coaching roles in various sports"
      }
    }
  ],
  "message": "Skill subcategories retrieved successfully"
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

### 400 Bad Request
```json
{
  "success": false,
  "error": "Invalid category ID",
  "message": "The provided category ID is not a valid UUID"
}
```

### 404 Not Found
```json
{
  "success": false,
  "error": "Category not found",
  "message": "The specified category does not exist"
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "error": "Failed to retrieve skill subcategories",
  "message": "Internal server error occurred while fetching skill subcategories"
}
```

## Notas Importantes

1. **Licencia Requerida**: Los endpoints requieren el header `X-License` con la licencia válida.
2. **Rate Limiting**: Los endpoints están sujetos a rate limiting (100 requests por minuto).
3. **CORS**: Los endpoints soportan CORS para requests desde dominios autorizados.
4. **Formato de Fechas**: Todas las fechas están en formato ISO 8601 (RFC3339).
5. **UUIDs**: Todos los IDs son UUIDs v4.
6. **Filtrado**: Se puede filtrar por categoría usando query parameter o URL parameter.

## Subcategorías Disponibles por Categoría

### Coach
- **Soccer**: Fútbol
- **Rugby Union**: Rugby Union
- **Rugby League**: Rugby League
- **Australian Rules Football (AFL)**: Fútbol Australiano
- **Basketball**: Baloncesto
- **Netball**: Netball
- **Tennis**: Tenis
- **Swimming**: Natación
- **Athletics**: Atletismo
- **Cricket**: Cricket
- **Hockey**: Hockey
- **Volleyball**: Voleibol
- **Gymnastics**: Gimnasia
- **Baseball**: Béisbol
- **Softball**: Softball
- **Surfing**: Surf
- **Cycling**: Ciclismo
- **Boxing**: Boxeo
- **Martial Arts**: Artes Marciales
- **Golf**: Golf
- **Other**: Otros

### Personal Trainer (PT)
- **Fitness**: Fitness
- **Gym**: Gimnasio
- **CrossFit**: CrossFit
- **Functional Training**: Entrenamiento Funcional
- **Athletics**: Atletismo
- **Swimming**: Natación
- **Cycling**: Ciclismo
- **Boxing**: Boxeo
- **Martial Arts**: Artes Marciales
- **Yoga**: Yoga
- **Pilates**: Pilates
- **Other**: Otros

## Casos de Uso

### Para Trabajadores (Labour)
Las subcategorías de habilidades ayudan a los trabajadores a:
- Especificar deportes o áreas específicas de expertise
- Mostrar especializaciones detalladas
- Filtrar oportunidades por especialización específica

### Para Empleadores (Builders)
Las subcategorías de habilidades permiten a los empleadores:
- Publicar trabajos con especializaciones específicas
- Filtrar candidatos por deporte o área específica
- Buscar talento especializado en áreas particulares

## Integración con Otros Endpoints

Las subcategorías se integran con:
- **Skill Categories**: Para obtener la categoría padre
- **Skills Complete**: Para obtener categorías con sus subcategorías
- **Labour Profiles**: Para definir especializaciones del trabajador
- **Job Postings**: Para especificar especializaciones requeridas
