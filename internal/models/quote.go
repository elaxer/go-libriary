package models

// Quote ...
type Quote struct {
	ID       int    `json:"id"`
	BookID   int    `json:"book_id"`
	UserID   int    `json:"user_id"`
	Text     string `json:"text"`
	IsHidden bool   `json:"is_hidden"`
}

// GetFields ...
func (q *Quote) GetFields() []interface{} {
	return []interface{}{
		&q.ID,
		&q.BookID,
		&q.UserID,
		&q.Text,
		&q.IsHidden,
	}
}

// ValidateCreation ...
func (q *Quote) ValidateCreation() error {
	return validation.ValidateStruct(
		q,
		validation.Field(&q.BookID, validation.Required),
		validation.Field(&q.UserID, validation.Required),
		validation.Field(&q.Text, validation.Required),
	)
}
