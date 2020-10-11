package storage

import (
	"books/internal/models"
	"fmt"
)

// Author ...
type Author struct {
	stor *Storage
}

// List ...
func (a *Author) List(limit, offset int) ([]*models.Author, error) {
	var authors []*models.Author

	rows, err := a.stor.db.Query("SELECT * FROM authors LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return authors, err
	}

	defer rows.Close()

	for rows.Next() {
		author := &models.Author{}
		if err := rows.Scan(author.GetFields()...); err != nil {
			return authors, err
		}

		authors = append(authors, author)
	}

	return authors, nil
}

// IsExists ...
func (a *Author) IsExists(id int) (bool, error) {
	isExists := false
	var count int

	err := a.stor.db.QueryRow("SELECT COUNT(*) FROM authors WHERE id = $1", id).Scan(&count)
	if err != nil {
		return isExists, err
	}

	if count > 0 {
		isExists = true
	}

	return isExists, nil
}

// Get ...
func (a *Author) Get(id int) (*models.Author, error) {
	author := new(models.Author)

	row := a.stor.db.QueryRow("SELECT * FROM authors WHERE id = $1", id)
	if err := row.Scan(author.GetFields()...); err != nil {
		return author, err
	}

	return author, nil
}

// Create ...
func (a *Author) Create(author *models.Author) (int, error) {
	var id int

	row := a.stor.db.QueryRow(
		`INSERT INTO authors
        (full_name, image, biography, first_name, last_name, middle_name)
        VALUES
        ($1, $2, $3, $4, $5, $6)
        RETURNING id`,
		author.FullName,
		author.Image,
		author.Biography,
		author.FirstName,
		author.LastName,
		author.MiddleName,
	)

	if err := row.Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}

// Delete ...
func (a *Author) Delete(id int) error {
	res, err := a.stor.db.Exec("DELETE FROM authors WHERE id = $1", id)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

// Update ...
// Удалить возвращаемый обновляемый объект
// рефакторинг
func (a *Author) Update(id int, updates map[string]interface{}) (*models.Author, error) {
	updatedAuthor := new(models.Author)

	query := "UPDATE authors SET "
	var placeholders []interface{}
	i := 1
	for k, v := range updates {
		if k == "id" {
			continue
		}

		query += fmt.Sprintf("%s=$%d, ", k, i)
		placeholders = append(placeholders, v)
		i++
	}

	if len(placeholders) == 0 {
		return updatedAuthor, nil
	}

	query = query[:len(query)-2] + fmt.Sprintf(" WHERE id = $%d", i)
	placeholders = append(placeholders, id)

	_, err := a.stor.db.Exec(query, placeholders...)
	if err != nil {
		return updatedAuthor, err
	}

	updatedAuthor, err = a.Get(id)
	if err != nil {
		return updatedAuthor, err
	}

	return updatedAuthor, nil
}

// QuotesList ...
func (a *Author) QuotesList(id int) ([]*models.Quote, error) {
	db := a.stor.db

	var quotes []*models.Quote

	rows, err := db.Query(`
        SELECT quotes.* FROM quotes
        LEFT JOIN books ON quotes.book_id = books.id
        WHERE books.author_id = $1
    `, id)

	if err != nil {
		return quotes, err
	}

	defer rows.Close()

	for rows.Next() {
		quote := new(models.Quote)

		if err := rows.Scan(quote.GetFields()...); err != nil {
			return quotes, err
		}

		quotes = append(quotes, quote)
	}

	return quotes, nil
}
