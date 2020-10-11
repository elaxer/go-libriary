package storage

import (
	"books/internal/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // ...
)

// Storage ...
type Storage struct {
	config *config.Config
	db     *sql.DB
	*Book
	*Tag
	*Author
	*User
	*Quote
	*Review
}

// New ...
func New(config *config.Config) *Storage {
	s := &Storage{
		config: config,
	}

	s.Book = &Book{stor: s}
	s.Tag = &Tag{stor: s}
	s.Author = &Author{stor: s}
	s.User = &User{store: s}
	s.Quote = &Quote{store: s}
	s.Review = &Review{store: s}

	return s
}

// Open ...
func (s *Storage) Open() error {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		s.config.DB.Host,
		s.config.DB.Port,
		s.config.DB.User,
		s.config.DB.Password,
		s.config.DB.DBName,
		s.config.DB.SSLMode,
	)

	db, err := sql.Open(s.config.DB.Engine, connStr)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	s.db = db

	return nil
}

// MustOpen ...
func (s *Storage) MustOpen() {
	if err := s.Open(); err != nil {
		panic(err)
	}
}

// Close ...
func (s *Storage) Close() {
	s.db.Close()
}
