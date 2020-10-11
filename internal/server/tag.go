package server

import (
	"books/internal/models"
	"books/internal/storage"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func handleGetTags(srv *Server) http.HandlerFunc {
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

		tags, err := srv.storage.Tag.List(limit, offset)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		jsonResponse(w, r, http.StatusOK, tags)
	}
}

func handleCreateTag(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(models.Tag)

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		id, err := srv.storage.Tag.Create(req)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		tag, err := srv.storage.Tag.Get(id)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		jsonResponse(w, r, http.StatusCreated, tag)
	}
}

func handleGetTag(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		tag, err := srv.storage.Tag.Get(id)
		if err == sql.ErrNoRows {
			errorResponse(w, r, http.StatusNotFound, err)
			return
		} else if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
		}

		jsonResponse(w, r, http.StatusOK, tag)
	}
}

func handleGetTagBooks(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		isExists, err := srv.storage.Tag.IsExists(id)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		if !isExists {
			errorResponse(w, r, http.StatusNotFound, sql.ErrNoRows)
			return
		}

		books, err := srv.storage.Tag.BooksList(id)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		jsonResponse(w, r, http.StatusOK, books)
	}
}

func handleUpdateTag(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req map[string]interface{}

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		isExists, err := srv.storage.Tag.IsExists(id)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		if !isExists {
			errorResponse(w, r, http.StatusNotFound, err)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		if err := srv.storage.Tag.Update(id, req); err == storage.ErrUnallowedField {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		} else if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		tag, err := srv.storage.Tag.Get(id)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		jsonResponse(w, r, http.StatusOK, tag)
	}
}

func handleDeleteTag(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		isExists, err := srv.storage.Tag.IsExists(id)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		if !isExists {
			errorResponse(w, r, http.StatusNotFound, err)
			return
		}

		if err := srv.storage.Tag.Delete(id); err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		noContentResponse(w, r, http.StatusNoContent)
	}
}
