package dto

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required"`
}

type RegisterResponse struct {
	Email    string `json:"email" gorm:"type: varchar(255)"`
	Password string `json:"password" gorm:"type: varchar(255)"`
	Role     string `json:"role" gorm:"type: varchar(255)"`
}

type LoginResponse struct {
	Username string `json:"username" gorm:"type: varchar(255)"`
	Email    string `json:"email" gorm:"type: varchar(255)"`
	Password string `json:"password" gorm:"type: varchar(255)"`
	Token    string `json:"token" gorm:"type: varchar(255)"`
}

type CheckAuth struct {
	ID       int    `json:"id" form:"id"`
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Role     string `json:"role" gorm:"type: varchar(255)"`
}
