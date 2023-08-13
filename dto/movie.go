package dto

type CreateMovieRequest struct {
	Title       string `json:"title" form:"title" gorm:"type: varchar(255)" validate:"required"`
	Category    string `json:"category" form:"category" gorm:"type: varchar(255)" validate:"required"`
	Price       int    `json:"price" form:"price" gorm:"type: int" validate:"required"`
	Link        string `json:"link" form:"link" gorm:"type: varchar(255)" validate:"required"`
	Description string `json:"description" form:"description" gorm:"type: text" validate:"required"`
	Thumbnail   string `json:"thumbnail" form:"thumbnail" gorm:"type: varchar(255)" validate:"required"`
}

type UpdateMovieRequest struct {
	Title       string `json:"title" form:"title"`
	Category    string `json:"category" form:"category"`
	Price       int    `json:"price" form:"price"`
	Link        string `json:"link" form:"link"`
	Description string `json:"description" form:"description"`
	Thumbnail   string `json:"thumbnail" form:"thumbnail"`
}

type MovieResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title" validate:"required"`
	Category    string `json:"category" validate:"required"`
	Price       int    `json:"price" validate:"required"`
	Link        string `json:"link" validate:"required"`
	Description string `json:"description" validate:"required"`
	Thumbnail   string `json:"thumbnail" validate:"required"`
}
