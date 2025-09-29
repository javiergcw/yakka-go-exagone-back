# Create Job

Crea un nuevo trabajo en el sistema.

## Endpoint

```
POST /api/v1/builder/jobs
```

## Headers

```
Authorization: Bearer <token>
Content-Type: application/json
```

## Authentication

- **Required**: Bearer token authentication
- **Role**: User must have builder role
- **Builder Profile**: The system automatically extracts the `builder_profile_id` from the authenticated user's JWT token
- **No Manual Input**: Do not include `builder_profile_id` in the request body - it will be automatically set

## Request Body

```json
{
  "jobsite_id": "123e4567-e89b-12d3-a456-426614174001",
  "job_type_id": "123e4567-e89b-12d3-a456-426614174002",
  "many_labours": 5,
  "ongoing_work": false,
  "wage_site_allowance": 50.00,
  "wage_leading_hand_allowance": 25.00,
  "wage_productivity_allowance": 15.00,
  "extras_overtime_rate": 1.5,
  "wage_hourly_rate": 25.50,
  "travel_allowance": 15.00,
  "gst": 2.50,
  "start_date_work": "2024-02-01T00:00:00Z",
  "end_date_work": "2024-02-28T00:00:00Z",
  "work_saturday": true,
  "work_sunday": false,
  "start_time": "08:00:00",
  "end_time": "17:00:00",
  "description": "Construction work for new building project",
  "payment_day": 15,
  "requires_supervisor_signature": true,
  "supervisor_name": "John Smith",
  "visibility": "PUBLIC",
  "payment_type": "WEEKLY",
  "license_ids": [
    "123e4567-e89b-12d3-a456-426614174003",
    "123e4567-e89b-12d3-a456-426614174004"
  ],
  "skill_category_ids": [
    "123e4567-e89b-12d3-a456-426614174005",
    "123e4567-e89b-12d3-a456-426614174006"
  ]
}
```

## Field Descriptions

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `jobsite_id` | UUID | ID of the jobsite where the work will be performed |
| `job_type_id` | UUID | ID of the job type (e.g., construction, maintenance) |
| `many_labours` | Integer | Number of labourers needed for this job |
| `visibility` | String | Job visibility: `DRAFT`, `PUBLIC`, `PRIVATE` |
| `payment_type` | String | Payment frequency: `WEEKLY`, `DAILY`, `FIXED_DAY` |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `ongoing_work` | Boolean | Whether this is ongoing work (default: false) |
| `wage_site_allowance` | Float | Site allowance amount |
| `wage_leading_hand_allowance` | Float | Leading hand allowance amount |
| `wage_productivity_allowance` | Float | Productivity allowance amount |
| `extras_overtime_rate` | Float | Overtime rate multiplier |
| `wage_hourly_rate` | Float | Hourly wage rate |
| `travel_allowance` | Float | Travel allowance amount |
| `gst` | Float | Goods and Services Tax amount |
| `start_date_work` | DateTime | Start date of work |
| `end_date_work` | DateTime | End date of work |
| `work_saturday` | Boolean | Whether work is required on Saturday (default: false) |
| `work_sunday` | Boolean | Whether work is required on Sunday (default: false) |
| `start_time` | String | Start time in HH:MM:SS format |
| `end_time` | String | End time in HH:MM:SS format |
| `description` | String | Job description |
| `payment_day` | Integer | Day of month for FIXED_DAY payment type |
| `requires_supervisor_signature` | Boolean | Whether supervisor signature is required (default: false) |
| `supervisor_name` | String | Name of the supervisor |
| `license_ids` | Array[UUID] | Array of required license IDs |
| `skill_category_ids` | Array[UUID] | Array of required skill category IDs |

## Response

### Success Response (201 Created)

```json
{
  "success": true,
  "job": {
    "id": "123e4567-e89b-12d3-a456-426614174007",
    "builder_profile_id": "123e4567-e89b-12d3-a456-426614174008",
    "jobsite_id": "123e4567-e89b-12d3-a456-426614174001",
    "job_type_id": "123e4567-e89b-12d3-a456-426614174002",
    "many_labours": 5,
    "ongoing_work": false,
    "wage_site_allowance": 50.00,
    "wage_leading_hand_allowance": 25.00,
    "wage_productivity_allowance": 15.00,
    "extras_overtime_rate": 1.5,
    "wage_hourly_rate": 25.50,
    "travel_allowance": 15.00,
    "gst": 2.50,
    "start_date_work": "2024-02-01T00:00:00Z",
    "end_date_work": "2024-02-28T00:00:00Z",
    "work_saturday": true,
    "work_sunday": false,
    "start_time": "08:00:00",
    "end_time": "17:00:00",
    "description": "Construction work for new building project",
    "payment_day": 15,
    "requires_supervisor_signature": true,
    "supervisor_name": "John Smith",
    "visibility": "PUBLIC",
    "payment_type": "WEEKLY",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "job_licenses": [
      {
        "id": "123e4567-e89b-12d3-a456-426614174009",
        "job_id": "123e4567-e89b-12d3-a456-426614174007",
        "license_id": "123e4567-e89b-12d3-a456-426614174003",
        "created_at": "2024-01-15T10:30:00Z"
      }
    ],
    "job_skills": [
      {
        "id": "123e4567-e89b-12d3-a456-426614174010",
        "job_id": "123e4567-e89b-12d3-a456-426614174007",
        "skill_category_id": "123e4567-e89b-12d3-a456-426614174005",
        "skill_subcategory_id": null,
        "created_at": "2024-01-15T10:30:00Z"
      }
    ]
  },
  "message": "Job created successfully"
}
```

### Error Responses

#### 400 Bad Request - Invalid Data
```json
{
  "success": false,
  "message": "Invalid request data",
  "error": "Validation failed: jobsite_id is required"
}
```

#### 401 Unauthorized - Invalid Token
```json
{
  "success": false,
  "message": "Invalid token",
  "error": "Invalid token"
}
```

#### 403 Forbidden - Builder Role Required
```json
{
  "success": false,
  "message": "Access denied: Builder role required",
  "error": "Access denied: Builder role required"
}
```

#### 500 Internal Server Error
```json
{
  "success": false,
  "message": "Failed to create job",
  "error": "Database connection failed"
}
```

## Examples

### Minimal Job Creation
```json
{
  "jobsite_id": "123e4567-e89b-12d3-a456-426614174001",
  "job_type_id": "123e4567-e89b-12d3-a456-426614174002",
  "many_labours": 3,
  "visibility": "PUBLIC",
  "payment_type": "WEEKLY"
}
```

### Complete Job Creation
```json
{
  "jobsite_id": "123e4567-e89b-12d3-a456-426614174001",
  "job_type_id": "123e4567-e89b-12d3-a456-426614174002",
  "many_labours": 5,
  "ongoing_work": false,
  "wage_site_allowance": 50.00,
  "wage_leading_hand_allowance": 25.00,
  "wage_productivity_allowance": 15.00,
  "extras_overtime_rate": 1.5,
  "wage_hourly_rate": 25.50,
  "travel_allowance": 15.00,
  "gst": 2.50,
  "start_date_work": "2024-02-01T00:00:00Z",
  "end_date_work": "2024-02-28T00:00:00Z",
  "work_saturday": true,
  "work_sunday": false,
  "start_time": "08:00:00",
  "end_time": "17:00:00",
  "description": "Construction work for new building project",
  "payment_day": 15,
  "requires_supervisor_signature": true,
  "supervisor_name": "John Smith",
  "visibility": "PUBLIC",
  "payment_type": "WEEKLY",
  "license_ids": [
    "123e4567-e89b-12d3-a456-426614174003"
  ],
  "skill_category_ids": [
    "123e4567-e89b-12d3-a456-426614174005"
  ]
}
```

## Notes

- The `builder_profile_id` is automatically extracted from the JWT token and should not be included in the request body
- All monetary fields (wages, allowances, GST) are stored as decimal values with 2 decimal places
- Time fields should be in HH:MM:SS format
- Date fields should be in ISO 8601 format
- The job will be created with `DRAFT` visibility by default unless specified otherwise
- License and skill requirements are optional but recommended for better job matching
