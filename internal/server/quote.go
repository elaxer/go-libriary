package server

import (
	"books/internal/models"
	"books/internal/storage"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"net/http"
	"strconv"
)

func handleGetQuotes(srv *Server) http.HandlerFunc {
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

		quotes, err := srv.storage.Quote.List(limit, offset)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		jsonResponse(w, r, http.StatusOK, quotes)
	}
}

func handleCreateQuote(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		quote := new(models.Quote)

		if err := json.NewDecoder(r.Body).Decode(quote); err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		if err := quote.ValidateCreation(); err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		err := srv.storage.Quote.Create(quote)
		if err, ok := err.(*pq.Error); ok {
			if err.Code == storage.ForeignKeyViolationCode {
				errorResponse(w, r, http.StatusBadRequest, err)
				return
			}

			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		jsonResponse(w, r, http.StatusCreated, quote)
	}
}

func handleGetQuote(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		quote, err := srv.storage.Quote.Get(id)
		if err == sql.ErrNoRows {
			errorResponse(w, r, http.StatusNotFound, err)
			return
		} else if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		jsonResponse(w, r, http.StatusOK, quote)
	}
}

func handleUpdateQuote(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func handleDeleteQuote(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])

		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		err = srv.storage.Quote.Delete(id)
		if err == storage.ErrNoRowsAffected {
			errorResponse(w, r, http.StatusNotFound, err)
			return
		} else if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		noContentResponse(w, r, http.StatusNoContent)
	}
}
