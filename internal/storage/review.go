package storage

import (
	"books/internal/models"
)

// Review ...
type Review struct {
	store *Storage
}

// List ...
func (r *Review) List(limit, offset int) ([]*models.Review, error) {
	db := r.store.db

	var reviews []*models.Review

	rows, err := db.Query("SELECT * FROM reviews LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return reviews, err
	}

	defer rows.Close()

	for rows.Next() {
		review := new(models.Review)

		if err := rows.Scan(review.GetFields()...); err != nil {
			return reviews, err
		}

		reviews = append(reviews, review)
	}

	return reviews, nil
}

// Get ...
func (r *Review) Get(id int) (*models.Review, error) {
	db := r.store.db

	review := new(models.Review)

	row := db.QueryRow("SELECT * FROM reviews WHERE id = $1", id)
	if err := row.Scan(review.GetFields()...); err != nil {
		return review, err
	}

	return review, nil
}
