package dto

type RegisterRequest struct {
	FirstName    string `json:"first_name" validate:"required,min=1,max=255"`
	LastName     string `json:"last_name" validate:"required,min=1,max=255"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required,min=6"`
	Age          int    `json:"age" validate:"required,min=1,max=150"`
	BelifRating  int    `json:"belif_rating" validate:"required,min=1,max=5"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Access string `json:"access"`
}

type RegisterResponse struct {
	Access string      `json:"access"`
	User   UserResponse `json:"user"`
}

