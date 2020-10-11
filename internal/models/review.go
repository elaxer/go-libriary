package models

// Review ...
type Review struct {
	ID          int    `json:"id"`
	BookID      int    `json:"book_id"`
	UserID      int    `json:"user_id"`
	Text        string `json:"string"`
	HasSpoilers bool   `json:"has_spoilers"`
	IsHidden    bool   `json:"is_hidden"`
}

// GetFields ...
func (r *Review) GetFields() []interface{} {
	return []interface{}{
		&r.ID,
		&r.BookID,
		&r.UserID,
		&r.Text,
		&r.HasSpoilers,
		&r.IsHidden,
	}
}
