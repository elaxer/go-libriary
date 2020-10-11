package server

import (
	"books/internal/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func handleGetBooks(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, err := getIntQuery(r, "limit", 100)
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		offset, err := getIntQuery(r, "offset", 0)
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		books, err := srv.storage.Book.List(limit, offset)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		jsonResponse(w, r, http.StatusOK, books)
	}
}

func handleGetBook(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		book, err := srv.storage.Book.Get(id)
		if err == sql.ErrNoRows {
			errorResponse(w, r, http.StatusNotFound, err)
			return
		} else if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		jsonResponse(w, r, http.StatusOK, book)
	}
}

// Refactor
func handleCreateBook(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(models.Book)

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		isExists, err := srv.storage.Author.IsExists(req.AuthorID)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		if !isExists {
			errorResponse(w, r, http.StatusNotFound, sql.ErrNoRows)
			return
		}

		id, err := srv.storage.Book.Create(req)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		book, err := srv.storage.Book.Get(id)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		jsonResponse(w, r, http.StatusOK, book)
	}
}

func handleUpdateBook(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		isExists, err := srv.storage.Book.IsExists(id)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		if !isExists {
			errorResponse(w, r, http.StatusNotFound, sql.ErrNoRows)
			return
		}

		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		err = srv.storage.Book.Update(id, req)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		book, err := srv.storage.Book.Get(id)
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		jsonResponse(w, r, http.StatusOK, book)
	}
}

func handleDeleteBook(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		isExists, err := srv.storage.Book.IsExists(id)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		if !isExists {
			errorResponse(w, r, http.StatusNotFound, sql.ErrNoRows)
			return
		}

		if err := srv.storage.Book.Delete(id); err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		noContentResponse(w, r, http.StatusNoContent)
	}
}

func handleGetBookAuthor(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		book, err := srv.storage.Book.Get(id)
		if err == sql.ErrNoRows {
			errorResponse(w, r, http.StatusNotFound, sql.ErrNoRows)
			return
		} else if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		author, err := srv.storage.Author.Get(book.AuthorID)
		if err == sql.ErrNoRows {
			errorResponse(w, r, http.StatusNotFound, err)
			return
		} else if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		jsonResponse(w, r, http.StatusOK, author)
	}
}

// Рефактор
func handleAddTagToBook(srv *Server) http.HandlerFunc {
	type request struct {
		TagID int `json:"tag_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		isExists, err := srv.storage.Book.IsExists(id)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		if !isExists {
			errorResponse(w, r, http.StatusNotFound, sql.ErrNoRows)
			return
		}

		req := new(request)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		isExists, err = srv.storage.Tag.IsExists(req.TagID)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		if !isExists {
			errorResponse(w, r, http.StatusNotFound, sql.ErrNoRows)
			return
		}

		isExists, err = srv.storage.Book.HasTag(id, req.TagID)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		if isExists {
			errorResponse(
				w,
				r,
				http.StatusConflict,
				fmt.Errorf("book with id %d already have tag with id %d", id, req.TagID),
			)
			return
		}

		if err := srv.storage.Book.AddTag(id, req.TagID); err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		noContentResponse(w, r, http.StatusCreated)
	}
}

func handleGetBookTags(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		isExists, err := srv.storage.Book.IsExists(id)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		if !isExists {
			errorResponse(w, r, http.StatusNotFound, sql.ErrNoRows)
			return
		}

		tags, err := srv.storage.Book.TagsList(id)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		jsonResponse(w, r, http.StatusOK, tags)
	}
}

func handleGetBookQuotes(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		isExists, err := srv.storage.Book.IsExists(id)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		if !isExists {
			errorResponse(w, r, http.StatusNotFound, sql.ErrNoRows)
			return
		}

		quotes, err := srv.storage.Book.QuotesList(id)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		jsonResponse(w, r, http.StatusOK, quotes)
	}
}

func handleGetBookReviews(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		isExists, err := srv.storage.Book.IsExists(id)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		if !isExists {
			errorResponse(w, r, http.StatusNotFound, sql.ErrNoRows)
			return
		}

		limit, err := getIntQuery(r, "limit", 100)
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		offset, err := getIntQuery(r, "offset", 0)
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		reviews, err := srv.storage.Book.ReviewsList(id, limit, offset)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		jsonResponse(w, r, http.StatusOK, reviews)
	}
}

func handleRemoveBookTag(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		bookID, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		tagID, err := strconv.Atoi(vars["tag_id"])
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		hasTag, err := srv.storage.Book.HasTag(bookID, tagID)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		if !hasTag {
			errorResponse(w, r, http.StatusNotFound, sql.ErrNoRows)
			return
		}

		if err := srv.storage.Book.RemoveTag(bookID, tagID); err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		noContentResponse(w, r, http.StatusNoContent)
	}
}
