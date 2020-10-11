package storage

import (
	"books/internal/models"
	"fmt"
)

// Tag ...
type Tag struct {
	stor *Storage
}

// List ...
func (t *Tag) List(limit, offset int) ([]*models.Tag, error) {
	var tags []*models.Tag

	rows, err := t.stor.db.Query("SELECT * FROM book_tags LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return tags, err
	}

	defer rows.Close()

	for rows.Next() {
		tag := new(models.Tag)

		if err := rows.Scan(tag.GetFields()...); err != nil {
			return tags, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

// Get ...
func (t *Tag) Get(id int) (*models.Tag, error) {
	tag := new(models.Tag)

	err := t.stor.db.QueryRow("SELECT * FROM book_tags WHERE id = $1", id).Scan(tag.GetFields()...)
	if err != nil {
		return tag, err
	}

	return tag, nil
}

// IsExists ...
func (t *Tag) IsExists(id int) (bool, error) {
	var count int
	isExists := false

	err := t.stor.db.QueryRow("SELECT COUNT(*) FROM book_tags WHERE id = $1", id).Scan(&count)
	if err != nil {
		return isExists, err
	}

	if count > 0 {
		isExists = true
	}

	return isExists, nil
}

// Create ...
func (t *Tag) Create(tag *models.Tag) (int, error) {
	var id int

	row := t.stor.db.QueryRow("INSERT INTO book_tags (name) VALUES ($1) RETURNING id", tag.Name)
	if err := row.Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}

// BooksList ...
func (t *Tag) BooksList(id int) ([]*models.Book, error) {
	var books []*models.Book

	rows, err := t.stor.db.Query(`
        SELECT b.* FROM book_tags AS t
        LEFT JOIN books_tags AS bt ON bt.tag_id = t.id
        LEFT JOIN books AS b ON b.id = bt.book_id
        WHERE t.id = $1
    `, id)
	if err != nil {
		return books, err
	}

	defer rows.Close()

	for rows.Next() {
		book := new(models.Book)

		if err := rows.Scan(book.GetFields()...); err != nil {
			return books, err
		}

		books = append(books, book)
	}

	return books, nil
}

// Update ...
func (t *Tag) Update(id int, updates map[string]interface{}) error {
	query := "UPDATE book_tags SET "
	var placeholders []interface{}
	i := 1

	for k, v := range updates {
		isAvailable := false
		for _, field := range models.TagUpdatableFields {
			if k == field {
				isAvailable = true
				break
			}
		}

		if !isAvailable {
			return ErrUnallowedField
		}

		query += fmt.Sprintf("%s = $%d, ", k, i)
		placeholders = append(placeholders, v)
		i++
	}

	if len(placeholders) == 0 {
		return nil
	}

	query = query[:len(query)-2] + fmt.Sprintf(" WHERE id = $%d", i)
	placeholders = append(placeholders, id)

	if _, err := t.stor.db.Exec(query, placeholders...); err != nil {
		return err
	}

	return nil
}

// Delete ...
func (t *Tag) Delete(id int) error {
	db := t.stor.db

	if _, err := db.Exec("DELETE FROM book_tags WHERE id = $1", id); err != nil {
		return err
	}

	return nil
}
