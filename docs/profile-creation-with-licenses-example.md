# Profile Creation with Licenses - Examples

## Labour Profile with Skills and Licenses

```json
{
  "first_name": "Juan",
  "last_name": "Pérez",
  "location": "Ciudad de México, México",
  "bio": "Constructor con 10 años de experiencia en construcción residencial",
  "avatar_url": "https://example.com/avatar.jpg",
  "phone": "+525512345678",
  "skills": [
    {
      "category_id": "550e8400-e29b-41d4-a716-446655440001",
      "subcategory_id": "550e8400-e29b-41d4-a716-446655440002", 
      "experience_level_id": "550e8400-e29b-41d4-a716-446655440003",
      "years_experience": 5.5,
      "is_primary": true
    }
  ],
  "licenses": [
    {
      "license_id": "550e8400-e29b-41d4-a716-446655440010",
      "photo_url": "https://example.com/license1.jpg",
      "issued_at": "2020-01-15T00:00:00Z",
      "expires_at": "2025-01-15T00:00:00Z"
    },
    {
      "license_id": "550e8400-e29b-41d4-a716-446655440011",
      "photo_url": "https://example.com/license2.jpg",
      "issued_at": "2019-06-01T00:00:00Z",
      "expires_at": "2024-06-01T00:00:00Z"
    }
  ]
}
```

## Builder Profile with Licenses

```json
{
  "company_name": "Constructora ABC S.A.",
  "display_name": "Juan Constructor",
  "location": "Ciudad de México, México",
  "bio": "Empresa constructora con 15 años de experiencia",
  "avatar_url": "https://example.com/avatar.jpg",
  "phone": "+525512345678",
  "licenses": [
    {
      "license_id": "550e8400-e29b-41d4-a716-446655440020",
      "photo_url": "https://example.com/builder-license.jpg",
      "issued_at": "2018-03-01T00:00:00Z",
      "expires_at": "2026-03-01T00:00:00Z"
    }
  ]
}
```

## Field Descriptions

### Common Profile Fields
- `first_name`/`company_name`: Name (required)
- `last_name`/`display_name`: Last name or display name (required)
- `location`: Location (required)
- `bio`: Optional biography
- `avatar_url`: Optional profile picture URL
- `phone`: Optional phone number

### Labour Profile Specific
- `skills`: Array of skills (optional)
  - `category_id`: UUID of skill category
  - `subcategory_id`: UUID of skill subcategory
  - `experience_level_id`: UUID of experience level
  - `years_experience`: Years of experience (0-99.9)
  - `is_primary`: Boolean for primary skill

### Licenses (Both Profiles)
- `licenses`: Array of user licenses (optional)
  - `license_id`: UUID of the license type (required)
  - `photo_url`: Optional URL to license image
  - `issued_at`: Optional issue date (ISO 8601 format)
  - `expires_at`: Optional expiration date (ISO 8601 format)

## Validation Rules
- All UUIDs must be valid and exist in their respective master tables
- Dates must be in ISO 8601 format (2006-01-02T15:04:05Z07:00)
- Years of experience must be between 0 and 99.9
- Only one skill can be marked as primary per profile
- No duplicate licenses allowed per user (enforced by unique constraint)
- No duplicate skill subcategories allowed per profile (enforced by unique constraint)

## Endpoints
- **Labour Profile**: `POST /api/v1/profiles/labour`
- **Builder Profile**: `POST /api/v1/profiles/builder`

Both endpoints now support creating profiles with associated licenses and skills (for labour profiles).
