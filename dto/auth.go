package dto

type RegisterRequest struct {
	Username string `json:"username" form:"username" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
	// VerificationToken string `json:"verification_token"`
}

type LoginRequest struct {
	Email    string `json:"email" form:"email" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
	Role     string `json:"role" form:"role" validate:"required"`
}

type RegisterResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	// VerificationToken string `json:"verification_token"`
}

type LoginResponse struct {
	// ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
	Role     string `json:"role"`
}

type CheckAuth struct {
	ID       int    `json:"id" form:"id"`
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Role     string `json:"role" gorm:"type: varchar(255)"`
}
