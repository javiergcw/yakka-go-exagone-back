# Labour Qualifications - Update

## Endpoint
```
POST /api/v1/labour/qualifications
```

## Description
Updates the qualifications for the authenticated labour profile. This endpoint performs a complete replacement of the user's qualifications - it deletes all existing qualifications and creates new ones based on the provided data. This ensures the labour profile has exactly the qualifications specified in the request.

## Authentication
- **Required**: Yes
- **Type**: Bearer Token (JWT)
- **Middleware**: `LabourMiddleware`
- **Role**: `labour` only

## Headers
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

## Request Body

### Schema
```json
{
  "qualifications": [
    {
      "qualification_id": "string (UUID)",
      "date_obtained": "string (ISO 8601, optional)",
      "expires_at": "string (ISO 8601, optional)",
      "status": "string (optional)"
    }
  ]
}
```

### Request Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `qualifications` | Array | Yes | List of qualifications to assign to the labour profile |
| `qualifications[].qualification_id` | String (UUID) | Yes | ID of the qualification from the master qualifications table |
| `qualifications[].date_obtained` | String (ISO 8601) | No | Date when the qualification was obtained |
| `qualifications[].expires_at` | String (ISO 8601) | No | Date when the qualification expires |
| `qualifications[].status` | String | No | Status of the qualification (defaults to "valid") |

### Example Request
```json
{
  "qualifications": [
    {
      "qualification_id": "550e8400-e29b-41d4-a716-446655440001",
      "date_obtained": "2023-01-15T00:00:00Z",
      "expires_at": "2025-01-15T00:00:00Z",
      "status": "valid"
    },
    {
      "qualification_id": "550e8400-e29b-41d4-a716-446655440002",
      "date_obtained": "2022-06-10T00:00:00Z",
      "expires_at": null,
      "status": "valid"
    },
    {
      "qualification_id": "550e8400-e29b-41d4-a716-446655440003",
      "date_obtained": "2023-03-20T00:00:00Z",
      "expires_at": "2024-03-20T00:00:00Z",
      "status": "expired"
    }
  ]
}
```

## Response

### Success Response (200 OK)
```json
{
  "qualifications": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440010",
      "qualification_id": "550e8400-e29b-41d4-a716-446655440001",
      "title": "UEFA A Diploma",
      "organization": "UEFA",
      "country": "International",
      "sport": "Football (Soccer)",
      "date_obtained": "2023-01-15T00:00:00Z",
      "expires_at": "2025-01-15T00:00:00Z",
      "status": "valid"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440011",
      "qualification_id": "550e8400-e29b-41d4-a716-446655440002",
      "title": "Smart Rugby (mandatory)",
      "organization": "Rugby Australia",
      "country": "Australia",
      "sport": "Rugby Union",
      "date_obtained": "2022-06-10T00:00:00Z",
      "expires_at": null,
      "status": "valid"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440012",
      "qualification_id": "550e8400-e29b-41d4-a716-446655440003",
      "title": "Foundation Coach",
      "organization": "Rugby Australia",
      "country": "Australia",
      "sport": "Rugby Union",
      "date_obtained": "2023-03-20T00:00:00Z",
      "expires_at": "2024-03-20T00:00:00Z",
      "status": "expired"
    }
  ],
  "message": "Labour qualifications updated successfully"
}
```

### Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `qualifications` | Array | List of newly created qualifications for the labour profile |
| `qualifications[].id` | String (UUID) | Unique identifier for the labour-qualification relationship |
| `qualifications[].qualification_id` | String (UUID) | ID of the qualification from the master table |
| `qualifications[].title` | String | Title/name of the qualification |
| `qualifications[].organization` | String | Organization that issued the qualification |
| `qualifications[].country` | String | Country where the qualification is valid |
| `qualifications[].sport` | String | Sport associated with the qualification |
| `qualifications[].date_obtained` | String (ISO 8601) | Date when the qualification was obtained |
| `qualifications[].expires_at` | String (ISO 8601) | Date when the qualification expires |
| `qualifications[].status` | String | Status of the qualification |
| `message` | String | Success message |

## Error Responses

### 400 Bad Request
```json
{
  "error": "Invalid request body",
  "message": "Request body must be valid JSON"
}
```

```json
{
  "error": "At least one qualification is required",
  "message": "The qualifications array cannot be empty"
}
```

```json
{
  "error": "Invalid qualification ID: 550e8400-e29b-41d4-a716-446655440999",
  "message": "The qualification ID must be a valid UUID"
}
```

```json
{
  "error": "Qualification not found: 550e8400-e29b-41d4-a716-446655440999",
  "message": "The specified qualification does not exist in the system"
}
```

### 401 Unauthorized
```json
{
  "error": "User ID not found in context",
  "message": "Authentication required"
}
```

### 500 Internal Server Error
```json
{
  "error": "Failed to delete existing qualifications",
  "message": "An internal server error occurred while removing old qualifications"
}
```

```json
{
  "error": "Failed to create qualification",
  "message": "An internal server error occurred while creating new qualifications"
}
```

## Example Usage

### cURL
```bash
curl -X POST "https://api.yakka.com/api/v1/labour/qualifications" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "qualifications": [
      {
        "qualification_id": "550e8400-e29b-41d4-a716-446655440001",
        "date_obtained": "2023-01-15T00:00:00Z",
        "expires_at": "2025-01-15T00:00:00Z",
        "status": "valid"
      }
    ]
  }'
```

### JavaScript (Fetch)
```javascript
const qualifications = [
  {
    qualification_id: "550e8400-e29b-41d4-a716-446655440001",
    date_obtained: "2023-01-15T00:00:00Z",
    expires_at: "2025-01-15T00:00:00Z",
    status: "valid"
  }
];

const response = await fetch('/api/v1/labour/qualifications', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer ' + token,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({ qualifications })
});

const data = await response.json();
console.log('Updated qualifications:', data.qualifications);
```

### Python (Requests)
```python
import requests
from datetime import datetime

qualifications = [
    {
        "qualification_id": "550e8400-e29b-41d4-a716-446655440001",
        "date_obtained": "2023-01-15T00:00:00Z",
        "expires_at": "2025-01-15T00:00:00Z",
        "status": "valid"
    }
]

headers = {
    'Authorization': 'Bearer ' + token,
    'Content-Type': 'application/json'
}

response = requests.post(
    'https://api.yakka.com/api/v1/labour/qualifications',
    headers=headers,
    json={'qualifications': qualifications}
)

data = response.json()
print(f"Updated {len(data['qualifications'])} qualifications")
```

## Important Notes

### Complete Replacement
- This endpoint performs a **complete replacement** of qualifications
- All existing qualifications for the labour profile are **deleted** first
- Only the qualifications specified in the request will remain
- If you want to keep existing qualifications, include them in the request

### Validation
- All `qualification_id` values must exist in the master qualifications table
- Invalid qualification IDs will result in a 400 error
- At least one qualification must be provided
- Empty qualifications array is not allowed

### Status Field
- If `status` is not provided, it defaults to "valid"
- Common status values: "valid", "expired", "pending"

### Date Format
- All dates must be in ISO 8601 format (e.g., "2023-01-15T00:00:00Z")
- Both `date_obtained` and `expires_at` are optional
- `expires_at` can be null for qualifications that don't expire

### Performance
- This operation is atomic - either all qualifications are updated or none
- If any qualification fails to create, the entire operation fails
- The endpoint returns the complete list of newly created qualifications
