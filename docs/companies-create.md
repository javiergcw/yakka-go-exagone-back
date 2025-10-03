# Create Company

## Endpoint
```
POST /api/v1/companies
```

## Description
Creates a new company in the system. This endpoint requires a valid license key in the request headers.

## Authentication
- **Required**: License key in headers
- **License Header**: `X-License-Key: YOUR_LICENSE_KEY`

## Request

### Headers
```
Content-Type: application/json
X-License-Key: YOUR_LICENSE_KEY
```

### Body
```json
{
  "name": "Constructora ABC",
  "description": "Empresa líder en construcción residencial y comercial",
  "website": "https://constructora-abc.com"
}
```

### Request Fields

| Field | Type | Required | Description | Validation |
|-------|------|----------|-------------|------------|
| `name` | string | Yes | Company name | min: 2, max: 255 characters |
| `description` | string | No | Company description | Optional |
| `website` | string | No | Company website URL | Must be valid URL format |

## Response

### Success Response (201 Created)
```json
{
  "company": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Constructora ABC",
    "description": "Empresa líder en construcción residencial y comercial",
    "website": "https://constructora-abc.com",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  },
  "message": "Company created successfully"
}
```

### Error Responses

#### 400 Bad Request - Invalid Data
```json
{
  "error": "Invalid request body",
  "message": "Validation failed: Name is required"
}
```

#### 400 Bad Request - Invalid URL
```json
{
  "error": "Invalid website URL",
  "message": "Validation failed: Website must be a valid URL"
}
```

#### 401 Unauthorized - Missing License
```json
{
  "error": "License required",
  "message": "Valid license key is required to access this endpoint"
}
```

#### 409 Conflict - Company Already Exists
```json
{
  "error": "Company with this name already exists",
  "message": "A company with this name already exists in the system"
}
```

#### 500 Internal Server Error
```json
{
  "error": "Failed to create company",
  "message": "An internal server error occurred"
}
```

## Example Usage

### cURL
```bash
curl -X POST "https://api.yakka.com/api/v1/companies" \
  -H "Content-Type: application/json" \
  -H "X-License-Key: YOUR_LICENSE_KEY" \
  -d '{
    "name": "Constructora ABC",
    "description": "Empresa líder en construcción residencial y comercial",
    "website": "https://constructora-abc.com"
  }'
```

### JavaScript (Fetch)
```javascript
const response = await fetch('https://api.yakka.com/api/v1/companies', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'X-License-Key': 'YOUR_LICENSE_KEY'
  },
  body: JSON.stringify({
    name: 'Constructora ABC',
    description: 'Empresa líder en construcción residencial y comercial',
    website: 'https://constructora-abc.com'
  })
});

const data = await response.json();
console.log(data);
```

### Python (requests)
```python
import requests

url = "https://api.yakka.com/api/v1/companies"
headers = {
    "Content-Type": "application/json",
    "X-License-Key": "YOUR_LICENSE_KEY"
}
data = {
    "name": "Constructora ABC",
    "description": "Empresa líder en construcción residencial y comercial",
    "website": "https://constructora-abc.com"
}

response = requests.post(url, headers=headers, json=data)
print(response.json())
```

## Notes

- **Company Names**: Must be unique across the system
- **Website URLs**: Must follow valid URL format (e.g., https://example.com)
- **Description**: Optional field for additional company information
- **License Required**: All company operations require a valid license key
- **Auto-generated Fields**: `id`, `created_at`, and `updated_at` are automatically generated

## Related Endpoints

- `GET /api/v1/companies` - Get all companies
- `POST /api/v1/builder/companies` - Assign company to builder profile
