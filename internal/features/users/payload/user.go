package payload

import (
	"time"

	"github.com/yakka-backend/internal/features/users/models"
)

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
}

// UserResponse represents the response payload for user data
type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToUserResponse converts a User model to UserResponse
func ToUserResponse(user *models.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToUserResponseList converts a slice of User models to UserResponse slice
func ToUserResponseList(users []models.User) []*UserResponse {
	responses := make([]*UserResponse, len(users))
	for i, user := range users {
		responses[i] = ToUserResponse(&user)
	}
	return responses
}
