package storage

import (
	"books/internal/models"
	"fmt"
)

// Book ...
type Book struct {
	stor *Storage
}

// List ...
func (b *Book) List(limit, offset int) ([]*models.Book, error) {
	var books []*models.Book

	rows, err := b.stor.db.Query("SELECT * FROM books LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return books, err
	}

	defer rows.Close()

	for rows.Next() {
		book := &models.Book{}
		if err := rows.Scan(book.GetFields()...); err != nil {
			return books, err
		}

		books = append(books, book)
	}

	return books, nil
}

// Create ...
func (b *Book) Create(book *models.Book) (int, error) {
	var id int
	row := b.stor.db.QueryRow(
		`INSERT INTO books
        (title, year_of_publishing, language, description, author_id, image)
        VALUES
        ($1, $2, $3, $4, $5, $6)
        RETURNING id`,
		book.Title,
		book.YearOfPublishing,
		book.Language,
		book.Description,
		book.AuthorID,
		book.Image,
	)

	if err := row.Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}

// IsExists ...
func (b *Book) IsExists(id int) (bool, error) {
	var count int
	isExists := false

	err := b.stor.db.QueryRow("SELECT COUNT(*) FROM books WHERE id = $1", id).Scan(&count)
	if err != nil {
		return isExists, err
	}

	if count > 0 {
		isExists = true
	}

	return isExists, nil
}

// Get ...
func (b *Book) Get(id int) (*models.Book, error) {
	book := new(models.Book)

	row := b.stor.db.QueryRow("SELECT * FROM books WHERE id = $1", id)
	if err := row.Scan(book.GetFields()...); err != nil {
		return book, err
	}

	return book, nil
}

// ListByAuthorID ...
func (b *Book) ListByAuthorID(id int) ([]*models.Book, error) {
	var books []*models.Book

	rows, err := b.stor.db.Query("SELECT * FROM books WHERE author_id = $1", id)
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
func (b *Book) Update(id int, updates map[string]interface{}) error {
	query := "UPDATE books SET "
	var placeholders []interface{}
	i := 1

	for k, v := range updates {
		if k == "id" {
			continue
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

	_, err := b.stor.db.Exec(query, placeholders...)
	if err != nil {
		return err
	}

	return nil
}

// Delete ...
func (b *Book) Delete(id int) error {
	_, err := b.stor.db.Exec("DELETE FROM books WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

// AddTag ...
func (b *Book) AddTag(bookID, tagID int) error {
	_, err := b.stor.db.Exec("INSERT INTO books_tags (book_id, tag_id) VALUES ($1, $2)", bookID, tagID)
	if err != nil {
		return err
	}

	return nil
}

// TagsList ...
func (b *Book) TagsList(id int) ([]*models.Tag, error) {
	var tags []*models.Tag

	rows, err := b.stor.db.Query(`
        SELECT t.id, t.name FROM books AS b
        LEFT JOIN books_tags AS bt ON bt.book_id = b.id
        LEFT JOIN book_tags AS t ON t.id = bt.tag_id
        WHERE b.id = $1;
    `, id)

	if err != nil {
		return tags, err
	}

	defer rows.Close()

	for rows.Next() {
		tag := new(models.Tag)

		err := rows.Scan(tag.GetFields()...)
		if err != nil {
			return tags, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

// QuotesList ...
func (b *Book) QuotesList(id int) ([]*models.Quote, error) {
	db := b.stor.db
	var quotes []*models.Quote

	rows, err := db.Query("SELECT * FROM quotes WHERE book_id = $1", id)

	if err != nil {
		return quotes, err
	}

	defer rows.Close()

	for rows.Next() {
		quote := new(models.Quote)

		err := rows.Scan(quote.GetFields()...)
		if err != nil {
			return quotes, err
		}

		quotes = append(quotes, quote)
	}

	return quotes, nil
}

// ReviewsList ...
func (b *Book) ReviewsList(id, limit, offset int) ([]*models.Review, error) {
	db := b.stor.db

	var reviews []*models.Review

	rows, err := db.Query("SELECT * FROM reviews WHERE book_id = $1 LIMIT $2 OFFSET $3", id, limit, offset)
	if err != nil {
		return reviews, err
	}

	for rows.Next() {
		review := new(models.Review)

		if err := rows.Scan(review.GetFields()...); err != nil {
			return reviews, err
		}

		reviews = append(reviews, review)
	}

	return reviews, nil
}

// RemoveTag ...
func (b *Book) RemoveTag(bookID, tagID int) error {
	_, err := b.stor.db.Exec("DELETE FROM books_tags WHERE book_id = $1 AND tag_id = $2", bookID, tagID)
	if err != nil {
		return err
	}

	return nil
}

// HasTag ...
func (b *Book) HasTag(bookID, tagID int) (bool, error) {
	var count int
	hasTag := false

	row := b.stor.db.QueryRow(`
        SELECT COUNT(*)
        FROM books_tags
        WHERE book_id = $1 AND tag_id = $2
    `, bookID, tagID)

	if err := row.Scan(&count); err != nil {
		return hasTag, err
	}

	if count > 0 {
		hasTag = true
	}

	return hasTag, nil
}
