package payload

import "time"

// PaymentConstantResponse represents a payment constant in responses
type PaymentConstantResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Value       int       `json:"value"`
	Description *string   `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreatePaymentConstantResponse represents the response when creating a payment constant
type CreatePaymentConstantResponse struct {
	Constant PaymentConstantResponse `json:"constant"`
	Message  string                  `json:"message"`
}

// UpdatePaymentConstantResponse represents the response when updating a payment constant
type UpdatePaymentConstantResponse struct {
	Constant PaymentConstantResponse `json:"constant"`
	Message  string                  `json:"message"`
}

// GetPaymentConstantsResponse represents the response when getting all payment constants
type GetPaymentConstantsResponse struct {
	Constants []PaymentConstantResponse `json:"constants"`
	Message   string                    `json:"message"`
}

// GetPaymentConstantResponse represents the response when getting a single payment constant
type GetPaymentConstantResponse struct {
	Constant PaymentConstantResponse `json:"constant"`
	Message  string                  `json:"message"`
}

// UpdatePaymentConstantValueResponse represents the response when updating only the value
type UpdatePaymentConstantValueResponse struct {
	Constant PaymentConstantResponse `json:"constant"`
	Message  string                  `json:"message"`
}
