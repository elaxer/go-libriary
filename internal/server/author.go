package server

import (
	"books/internal/models"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func handleGetAuthors(srv *Server) http.HandlerFunc {
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

		authors, err := srv.storage.Author.List(limit, offset)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		jsonResponse(w, r, http.StatusOK, authors)
	}
}

func handleGetAuthor(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		author, err := srv.storage.Author.Get(id)
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

func handleGetAuthorBooks(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		isExists, err := srv.storage.Author.IsExists(id)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		if !isExists {
			errorResponse(w, r, http.StatusNotFound, sql.ErrNoRows)
			return
		}

		// Сделать Author.BooksList
		books, err := srv.storage.Book.ListByAuthorID(id)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		if len(books) == 0 {
			noContentResponse(w, r, http.StatusNoContent)
			return
		}

		jsonResponse(w, r, http.StatusOK, books)
	}
}

// REfactor
func handleCreateAuthor(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		author := new(models.Author)

		if err := json.NewDecoder(r.Body).Decode(author); err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		id, err := srv.storage.Author.Create(author)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		author.ID = id

		jsonResponse(w, r, http.StatusOK, author)
	}
}

func handleDeleteAuthor(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		isExists, err := srv.storage.Author.IsExists(id)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		if !isExists {
			errorResponse(w, r, http.StatusNotFound, sql.ErrNoRows)
			return
		}

		if err := srv.storage.Author.Delete(id); err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		noContentResponse(w, r, http.StatusNoContent)
	}
}

// Refactor
func handleUpdateAuthor(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		isExists, err := srv.storage.Author.IsExists(id)
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

		updatedAuthor, err := srv.storage.Author.Update(id, req)
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		jsonResponse(w, r, http.StatusOK, updatedAuthor)
	}
}

func handleGetAuthorQuotes(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		quotes, err := srv.storage.Author.QuotesList(id)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		jsonResponse(w, r, http.StatusOK, quotes)
	}
}
