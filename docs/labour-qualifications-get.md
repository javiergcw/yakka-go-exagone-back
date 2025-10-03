# Labour Qualifications - Get All

## Endpoint
```
GET /api/v1/labour/qualifications
```

## Description
Retrieves all qualifications associated with the authenticated labour profile. This endpoint returns the qualifications that the labour user has been assigned, including details about the qualification, sport, and personal details like date obtained and expiration.

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

## Response

### Success Response (200 OK)
```json
{
  "qualifications": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
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
      "id": "550e8400-e29b-41d4-a716-446655440002",
      "qualification_id": "550e8400-e29b-41d4-a716-446655440003",
      "title": "Smart Rugby (mandatory)",
      "organization": "Rugby Australia",
      "country": "Australia",
      "sport": "Rugby Union",
      "date_obtained": "2022-06-10T00:00:00Z",
      "expires_at": null,
      "status": "valid"
    }
  ],
  "total": 2,
  "message": "Labour qualifications retrieved successfully"
}
```

### Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `qualifications` | Array | List of qualifications assigned to the labour profile |
| `qualifications[].id` | String (UUID) | Unique identifier for the labour-qualification relationship |
| `qualifications[].qualification_id` | String (UUID) | ID of the qualification from the master table |
| `qualifications[].title` | String | Title/name of the qualification |
| `qualifications[].organization` | String | Organization that issued the qualification |
| `qualifications[].country` | String | Country where the qualification is valid |
| `qualifications[].sport` | String | Sport associated with the qualification |
| `qualifications[].date_obtained` | String (ISO 8601) | Date when the qualification was obtained (optional) |
| `qualifications[].expires_at` | String (ISO 8601) | Date when the qualification expires (optional) |
| `qualifications[].status` | String | Status of the qualification (valid, expired, pending) |
| `total` | Integer | Total number of qualifications |
| `message` | String | Success message |

## Error Responses

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
  "error": "Failed to get labour qualifications",
  "message": "An internal server error occurred"
}
```

## Example Usage

### cURL
```bash
curl -X GET "https://api.yakka.com/api/v1/labour/qualifications" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json"
```

### JavaScript (Fetch)
```javascript
const response = await fetch('/api/v1/labour/qualifications', {
  method: 'GET',
  headers: {
    'Authorization': 'Bearer ' + token,
    'Content-Type': 'application/json'
  }
});

const data = await response.json();
console.log('My qualifications:', data.qualifications);
```

### Python (Requests)
```python
import requests

headers = {
    'Authorization': 'Bearer ' + token,
    'Content-Type': 'application/json'
}

response = requests.get(
    'https://api.yakka.com/api/v1/labour/qualifications',
    headers=headers
)

data = response.json()
print(f"Total qualifications: {data['total']}")
```

## Notes

- The endpoint automatically filters qualifications based on the authenticated user's labour profile
- Only qualifications assigned to the current labour profile are returned
- The `date_obtained` and `expires_at` fields are optional and may be null
- The `status` field defaults to "valid" if not specified
- The endpoint includes full qualification details including sport information
- This endpoint is protected by `LabourMiddleware` and requires a valid JWT token with labour role
