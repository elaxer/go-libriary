package models

import (
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	UserCreatableFields = []string{"email", "password", "name"}
	UserUpdatableFields = []string{"email", "password", "name"}
)

// User ...
type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Password     string    `json:"-"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
}

// GetFields ...
func (u *User) GetFields() []interface{} {
	return []interface{}{
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.Name,
		&u.CreatedAt,
	}
}

// ValidateCreation ...
func (u *User) ValidateCreation() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 50)),
		validation.Field(&u.Name, validation.Required, validation.Length(2, 35)),
	)
}

// BeforeCreate ...
func (u *User) BeforeCreate() error {
	hash, err := u.hashPassword(u.Password)

	if err != nil {
		return err
	}

	u.PasswordHash = hash

	return nil
}

func (u *User) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}

// CheckPasswordHash ...
func (u *User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))

	return err == nil
}
