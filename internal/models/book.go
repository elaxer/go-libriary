package models

// Book ...
type Book struct {
	ID               int    `json:"id,omitempty"`
	Title            string `json:"title,omitempty"`
	YearOfPublishing int    `json:"year_of_publishing,omitempty"`
	Language         string `json:"language,omitempty"`
	Description      string `json:"description,omitempty"`
	AuthorID         int    `json:"author_id,omitempty"`
	Image            string `json:"image,omitempty"`
}

// GetFields ...
func (b *Book) GetFields() []interface{} {
	return []interface{}{
		&b.ID,
		&b.Title,
		&b.YearOfPublishing,
		&b.Language,
		&b.Description,
		&b.AuthorID,
		&b.Image,
	}
}
