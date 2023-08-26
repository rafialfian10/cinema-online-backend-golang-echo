package dto

type CreateMovieRequest struct {
	Title       string `json:"title" form:"title" gorm:"type: varchar(255)" validate:"required"`
	ReleaseDate string `json:"release_date" form:"release_date"`
	CategoryID  []int  `json:"category_id" form:"category_id" validate:"required"`
	Price       int    `json:"price" form:"price" gorm:"type: int" validate:"required"`
	Link        string `json:"link" form:"link" gorm:"type: varchar(255)" validate:"required"`
	Description string `json:"description" form:"description" gorm:"type: text" validate:"required"`
	Thumbnail   string `json:"thumbnail" form:"thumbnail" gorm:"type: varchar(255)" validate:"required"`
	Trailer     string `json:"trailer" form:"trailer" gorm:"type: varchar(255)" validate:"required"`
}

type UpdateMovieRequest struct {
	Title       string `json:"title" form:"title"`
	ReleaseDate string `json:"release_date" form:"release_date"`
	CategoryID  []int  `json:"category_id" form:"category_id"`
	Price       int    `json:"price" form:"price"`
	Link        string `json:"link" form:"link"`
	Description string `json:"description" form:"description"`
	Thumbnail   string `json:"thumbnail" form:"thumbnail"`
	Trailer     string `json:"trailer" form:"trailer"`
}
