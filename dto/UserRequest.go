package dto

import "time"

type UpdateUserRequest struct {
	FirstName        *string `json:"first_name,omitempty" validate:"omitempty,min=1,max=255"`
	LastName         *string `json:"last_name,omitempty" validate:"omitempty,min=1,max=255"`
	Age              *int    `json:"age,omitempty" validate:"omitempty,min=1,max=150"`
	BelieverCategory *int    `json:"believer_category,omitempty" validate:"omitempty,min=1,max=5"`
}

type UserResponse struct {
	ID               int       `json:"id"`
	Email            string    `json:"email"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	Name             string    `json:"name"`
	Age              int       `json:"age"`
	BelieverCategory int       `json:"believer_category"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

