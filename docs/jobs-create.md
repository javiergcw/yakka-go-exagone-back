# Create Job

Crea un nuevo trabajo en el sistema.

## Endpoint

```
POST /api/v1/jobs
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
    "123e4567-e89b-12d3-a456-426614174005"
  ],
  "skill_subcategory_ids": [
    "123e4567-e89b-12d3-a456-426614174006"
  ]
}
```

## Response

### Success Response (201 Created)

```json
{
  "job": {
    "id": "123e4567-e89b-12d3-a456-426614174007",
    "builder_profile_id": "123e4567-e89b-12d3-a456-426614174000",
    "jobsite_id": "123e4567-e89b-12d3-a456-426614174001",
    "job_type_id": "123e4567-e89b-12d3-a456-426614174002",
    "many_labours": 5,
    "ongoing_work": false,
    "wage_site_allowance": 50.00,
    "wage_leading_hand_allowance": 25.00,
    "wage_productivity_allowance": 15.00,
    "extras_overtime_rate": 1.5,
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
    "updated_at": "2024-01-15T10:30:00Z"
  },
  "message": "Job created successfully"
}
```

### Error Response (400 Bad Request)

```json
{
  "error": "Bad Request",
  "message": "Invalid request body"
}
```

### Error Response (401 Unauthorized)

```json
{
  "error": "Unauthorized",
  "message": "Invalid or missing authentication token"
}
```

## Field Descriptions

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `jobsite_id` | UUID | Yes | ID of the jobsite where work will be performed |
| `job_type_id` | UUID | Yes | ID of the job type (from job_types master table) |
| `many_labours` | Integer | Yes | Number of labourers needed (minimum 1) |
| `ongoing_work` | Boolean | No | Whether this is ongoing work (default: false) |
| `wage_site_allowance` | Float | No | Site allowance amount |
| `wage_leading_hand_allowance` | Float | No | Leading hand allowance amount |
| `wage_productivity_allowance` | Float | No | Productivity allowance amount |
| `extras_overtime_rate` | Float | No | Overtime rate multiplier |
| `start_date_work` | DateTime | No | Start date of work |
| `end_date_work` | DateTime | No | End date of work |
| `work_saturday` | Boolean | No | Whether work is required on Saturday (default: false) |
| `work_sunday` | Boolean | No | Whether work is required on Sunday (default: false) |
| `start_time` | String | No | Start time in HH:MM:SS format |
| `end_time` | String | No | End time in HH:MM:SS format |
| `description` | String | No | Job description |
| `payment_day` | Integer | No | Day of month for payment (for FIXED_DAY payment type) |
| `requires_supervisor_signature` | Boolean | No | Whether supervisor signature is required (default: false) |
| `supervisor_name` | String | No | Name of the supervisor |
| `visibility` | Enum | No | Job visibility (PUBLIC, PRIVATE, BANNED, ARCHIVED, DRAFT) |
| `payment_type` | Enum | No | Payment frequency (FIXED_DAY, WEEKLY, FORTNIGHTLY) |
| `license_ids` | Array | No | Array of license IDs required for this job |
| `skill_category_ids` | Array | No | Array of skill category IDs required |
| `skill_subcategory_ids` | Array | No | Array of skill subcategory IDs required |

## Examples

### Create a Basic Job

```bash
curl -X POST "http://localhost:8081/api/v1/jobs" \
  -H "Authorization: Bearer your-token-here" \
  -H "Content-Type: application/json" \
  -d '{
    "jobsite_id": "123e4567-e89b-12d3-a456-426614174001",
    "job_type_id": "123e4567-e89b-12d3-a456-426614174002",
    "many_labours": 3,
    "description": "Basic construction work",
    "visibility": "PUBLIC",
    "payment_type": "WEEKLY"
  }'
```

### Create a Complete Job with All Fields

```bash
curl -X POST "http://localhost:8081/api/v1/jobs" \
  -H "Authorization: Bearer your-token-here" \
  -H "Content-Type: application/json" \
  -d '{
    "jobsite_id": "123e4567-e89b-12d3-a456-426614174001",
    "job_type_id": "123e4567-e89b-12d3-a456-426614174002",
    "many_labours": 5,
    "ongoing_work": false,
    "wage_site_allowance": 50.00,
    "wage_leading_hand_allowance": 25.00,
    "wage_productivity_allowance": 15.00,
    "extras_overtime_rate": 1.5,
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
      "123e4567-e89b-12d3-a456-426614174005"
    ],
    "skill_subcategory_ids": [
      "123e4567-e89b-12d3-a456-426614174006"
    ]
  }'
```

## Notes

- The job creator must be authenticated and have a valid builder profile
- The `builder_profile_id` is automatically extracted from the JWT token - do not include it in the request
- All referenced IDs (jobsite_id, job_type_id) must exist in the database
- **Security**: The `jobsite_id` must belong to the authenticated builder - you can only create jobs for your own jobsites
- License and skill relationships are optional but recommended for better job matching
- Visibility defaults to DRAFT if not specified
- Payment type defaults to WEEKLY if not specified
