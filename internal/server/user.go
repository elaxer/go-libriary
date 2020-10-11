package server

import (
	"books/internal/models"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func handleGetUsers(srv *Server) http.HandlerFunc {
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

		users, err := srv.storage.User.List(limit, offset)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		jsonResponse(w, r, http.StatusOK, users)
	}
}

func handleCreateUser(srv *Server) http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := new(request)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		user := &models.User{
			Email:    req.Email,
			Password: req.Password,
			Name:     req.Name,
		}

		isExists, err := srv.storage.User.IsExistsByEmail(user.Email)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		if isExists {
			errorResponse(w, r, http.StatusBadRequest, errors.New("user with email already exists"))
			return
		}

		if err := user.ValidateCreation(); err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}
		if err := srv.storage.User.Create(user); err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		jsonResponse(w, r, http.StatusCreated, user)
	}

}

func handleGetUser(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		user, err := srv.storage.User.Get(id)
		if err == sql.ErrNoRows {
			errorResponse(w, r, http.StatusNotFound, err)
			return
		} else if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		jsonResponse(w, r, http.StatusOK, user)
	}
}

func handleUpdateUser(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		isExists, err := srv.storage.User.IsExists(id)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		if !isExists {
			errorResponse(w, r, http.StatusNotFound, err)
			return
		}

		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		if err := srv.storage.User.Update(id, req); err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		user, err := srv.storage.User.Get(id)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		jsonResponse(w, r, http.StatusOK, user)
	}
}

func handleDeleteUser(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		if err := srv.storage.User.Delete(id); err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		noContentResponse(w, r, http.StatusNoContent)
	}
}

func handleGetUserBooks(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		isExists, err := srv.storage.User.IsExists(id)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		if !isExists {
			errorResponse(w, r, http.StatusNotFound, err)
			return
		}

		books, err := srv.storage.User.BooksList(id)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		if len(books) == 0 {
			noContentResponse(w, r, http.StatusOK)
			return
		}

		jsonResponse(w, r, http.StatusOK, books)
	}
}

func handleGetUserQuotes(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func handleGetUserReviews(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func handleGetMe(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*models.User))
	}
}
