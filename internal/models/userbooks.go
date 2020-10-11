package models

// UserBook ...
type UserBook struct {
	UserID           int    `json:"user_id,omitempty"`
	BookID           int    `json:"book_id,omitempty"`
	Title            string `json:"title,omitempty"`
	YearOfPublishing int    `json:"year_of_publishing,omitempty"`
	Language         string `json:"language,omitempty"`
	Description      string `json:"description,omitempty"`
	AuthorID         int    `json:"author_id,omitempty"`
	Image            string `json:"image,omitempty"`
	Type             int    `json:"type,omitempty"`
	Rating           int    `json:"rating,omitempty"`
}

// GetFields ...
func (ub *UserBook) GetFields() []interface{} {
	return []interface{}{
		&ub.UserID,
		&ub.BookID,
		&ub.Title,
		&ub.YearOfPublishing,
		&ub.Language,
		&ub.Description,
		&ub.AuthorID,
		&ub.Image,
		&ub.Type,
		&ub.Rating,
	}
}
