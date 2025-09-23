# Labour Profile Creation with Skills - Example

## Request Body Example

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
    },
    {
      "category_id": "550e8400-e29b-41d4-a716-446655440004",
      "subcategory_id": "550e8400-e29b-41d4-a716-446655440005",
      "experience_level_id": "550e8400-e29b-41d4-a716-446655440006", 
      "years_experience": 3.0,
      "is_primary": false
    }
  ]
}
```

## Field Descriptions

### Basic Profile Fields
- `first_name`: User's first name (required)
- `last_name`: User's last name (required) 
- `location`: User's location (required)
- `bio`: Optional biography
- `avatar_url`: Optional profile picture URL
- `phone`: Optional phone number

### Skills Array
Each skill object contains:
- `category_id`: UUID of the skill category (from skill_categories table)
- `subcategory_id`: UUID of the skill subcategory (from skill_subcategories table)
- `experience_level_id`: UUID of the experience level (from experience_levels table)
- `years_experience`: Number of years of experience (0-99.9)
- `is_primary`: Boolean indicating if this is the user's primary skill

## Validation Rules
- Skills are optional (can create profile without skills)
- Each skill must have valid UUIDs for category, subcategory, and experience level
- Years of experience must be between 0 and 99.9
- Only one skill can be marked as primary per profile
- No duplicate subcategories allowed per profile (enforced by unique constraint)

## Response
The endpoint will return the created labour profile with all associated skills.
