package storage

import (
	"books/internal/models"
	"fmt"
)

// User ...
type User struct {
	store *Storage
}

// List ...
func (u *User) List(limit, offset int) ([]*models.User, error) {
	db := u.store.db

	var users []*models.User

	rows, err := db.Query("SELECT * FROM users LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return users, err
	}

	defer rows.Close()

	for rows.Next() {
		user := new(models.User)

		if err := rows.Scan(user.GetFields()...); err != nil {
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}

// Get ...
func (u *User) Get(id int) (*models.User, error) {
	db := u.store.db

	user := new(models.User)
	if err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(user.GetFields()...); err != nil {
		return user, err
	}

	return user, nil
}

// GetByEmail ...
func (u *User) GetByEmail(email string) (*models.User, error) {
	db := u.store.db

	user := new(models.User)
	if err := db.QueryRow("SELECT * FROM users WHERE email = $1", email).Scan(user.GetFields()...); err != nil {
		return user, err
	}

	return user, nil
}

// Create ...
func (u *User) Create(user *models.User) error {
	user.BeforeCreate()

	q := "INSERT INTO users (email, password_hash, name) VALUES ($1, $2, $3) RETURNING id, created_at"

	row := u.store.db.QueryRow(q, user.Email, user.PasswordHash, user.Name)

	return row.Scan(&user.ID, &user.CreatedAt)
}

// Update ...
func (u *User) Update(id int, updates map[string]interface{}) error {
	db := u.store.db

	query := "UPDATE users SET "
	var placeholders []interface{}
	i := 1

	for field, value := range updates {
		isAllowed := false
		for _, f := range models.UserUpdatableFields {
			if field == f {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			return ErrUnallowedField
		}

		query += fmt.Sprintf("%s = $%d, ", field, i)
		placeholders = append(placeholders, value)
		i++
	}

	if len(placeholders) == 0 {
		return nil
	}

	query = query[:len(query)-2] + fmt.Sprintf(" WHERE id = $%d", i)
	placeholders = append(placeholders, id)

	if _, err := db.Exec(query, placeholders...); err != nil {
		return err
	}

	return nil
}

// Delete ...
func (u *User) Delete(id int) error {
	db := u.store.db

	_, err := db.Exec("DELETE FROM users WHERE id = $1", id)

	return err
}

// BooksList ...
func (u *User) BooksList(id int) ([]*models.UserBook, error) {
	db := u.store.db

	var books []*models.UserBook

	rows, err := db.Query(`
        SELECT ub.book_id, ub.user_id, b.title,
            b.year_of_publishing, b.language, b.description,
            b.author_id, b.image, ub.type, ub.rating
        FROM users AS u
        JOIN users_books AS ub ON u.id = ub.user_id
        JOIN books AS b ON ub.book_id = b.id
        WHERE u.id = $1
    `, id)

	if err != nil {
		return books, err
	}

	defer rows.Close()

	for rows.Next() {
		book := new(models.UserBook)

		if err := rows.Scan(book.GetFields()...); err != nil {
			return books, err
		}

		books = append(books, book)
	}

	return books, nil
}

// IsExists ...
func (u *User) IsExists(id int) (bool, error) {
	db := u.store.db

	isExists := false
	var count int

	if err := db.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", id).Scan(&count); err != nil {
		return isExists, err
	}

	if count > 0 {
		isExists = true
	}

	return isExists, nil
}

// IsExistsByEmail ...
func (u *User) IsExistsByEmail(email string) (bool, error) {
	db := u.store.db

	isExists := false
	var count int

	if err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count); err != nil {
		return isExists, err
	}

	if count > 0 {
		isExists = true
	}

	return isExists, nil
}
