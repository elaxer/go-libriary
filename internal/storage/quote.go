package storage

import (
	"books/internal/models"
)

// Quote ...
type Quote struct {
	store *Storage
}

// List ...
func (q *Quote) List(limit, offset int) ([]*models.Quote, error) {
	db := q.store.db

	var quotes []*models.Quote

	rows, err := db.Query("SELECT * FROM quotes LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return quotes, err
	}

	for rows.Next() {
		quote := new(models.Quote)

		if err := rows.Scan(quote.GetFields()...); err != nil {
			return quotes, err
		}

		quotes = append(quotes, quote)
	}

	return quotes, nil
}

// Create ...
func (q *Quote) Create(quote *models.Quote) error {
	return q.store.db.QueryRow(
		"INSERT INTO quotes (book_id, user_id, text) VALUES ($1, $2, $3) RETURNING id, is_hidden",
		quote.BookID,
		quote.UserID,
		quote.Text,
	).Scan(&quote.ID, &quote.IsHidden)
}

// Get ...
func (q *Quote) Get(id int) (*models.Quote, error) {
	db := q.store.db

	quote := new(models.Quote)
	err := db.QueryRow("SELECT * FROM quotes WHERE id = $1", id).Scan(quote.GetFields()...)

	return quote, err
}

// Delete ...
func (q *Quote) Delete(id int) error {
	db := q.store.db

	res, err := db.Exec("DELETE FROM quotes WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNoRowsAffected
	}

	return nil
}
