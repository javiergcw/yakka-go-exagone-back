package payload

// CreatePaymentConstantRequest represents the request to create a payment constant
type CreatePaymentConstantRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Value       int     `json:"value" validate:"required,min=0"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=255"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

// UpdatePaymentConstantRequest represents the request to update a payment constant
type UpdatePaymentConstantRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Value       *int    `json:"value,omitempty" validate:"omitempty,min=0"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=255"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

// UpdatePaymentConstantValueRequest represents the request to update only the value of a payment constant
type UpdatePaymentConstantValueRequest struct {
	Value int `json:"value" validate:"required,min=0"`
}
