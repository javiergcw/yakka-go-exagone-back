# Master Tables API Documentation

Esta documentación describe los endpoints disponibles para acceder a las tablas maestras del sistema Yakka. Todos los endpoints requieren autenticación mediante licencia.

## Autenticación

Todos los endpoints de tablas maestras requieren el header `X-License` con una licencia válida.

**Header requerido:**
```
X-License: YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B
```

## Endpoints Disponibles

### 1. Licenses

#### GET /api/v1/licenses

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
      "name": "Licencia de Conducir",
      "description": "Permiso para conducir vehículos automotores",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid-here",
      "name": "Licencia de Construcción",
      "description": "Permiso para realizar trabajos de construcción",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "message": "Licenses retrieved successfully"
}
```

### 2. Experience Levels

#### GET /api/v1/experience-levels

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
    }
  ],
  "message": "Experience levels retrieved successfully"
}
```

### 3. Skill Categories

#### GET /api/v1/skill-categories

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
    }
  ],
  "message": "Skill categories retrieved successfully"
}
```

### 4. Skill Subcategories

#### GET /api/v1/skill-subcategories

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
    }
  ],
  "message": "Skill subcategories retrieved successfully"
}
```

#### GET /api/v1/skill-categories/{categoryId}/subcategories

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
    }
  ],
  "message": "Skill subcategories retrieved successfully"
}
```

### 5. Skills Complete

#### GET /api/v1/skills

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
  "error": "Failed to retrieve [resource]",
  "message": "Internal server error occurred while fetching data"
}
```

## Notas Importantes

1. **Licencia Requerida**: Todos los endpoints requieren el header `X-License` con la licencia válida.
2. **Rate Limiting**: Los endpoints están sujetos a rate limiting (100 requests por minuto).
3. **CORS**: Los endpoints soportan CORS para requests desde dominios autorizados.
4. **Formato de Fechas**: Todas las fechas están en formato ISO 8601 (RFC3339).
5. **UUIDs**: Todos los IDs son UUIDs v4.

## Ejemplos de Uso Común

### Obtener todas las categorías y sus subcategorías (Recomendado)
```bash
# Opción 1: Endpoint completo (más eficiente)
curl -X GET "https://api.yakka.com/api/v1/skills" \
  -H "X-License: YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B"
```

### Obtener categorías y subcategorías por separado
```bash
# 1. Obtener categorías
curl -X GET "https://api.yakka.com/api/v1/skill-categories" \
  -H "X-License: YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B"

# 2. Para cada categoría, obtener sus subcategorías
curl -X GET "https://api.yakka.com/api/v1/skill-categories/{categoryId}/subcategories" \
  -H "X-License: YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B"
```

### Obtener datos maestros completos
```bash
# Obtener todos los datos maestros
curl -X GET "https://api.yakka.com/api/v1/licenses" \
  -H "X-License: YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B"

curl -X GET "https://api.yakka.com/api/v1/experience-levels" \
  -H "X-License: YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B"

curl -X GET "https://api.yakka.com/api/v1/skill-categories" \
  -H "X-License: YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B"

curl -X GET "https://api.yakka.com/api/v1/skill-subcategories" \
  -H "X-License: YAKKA-PROD-2024-8F9E2A1B-3C4D5E6F-7A8B9C0D-1E2F3A4B"
```
