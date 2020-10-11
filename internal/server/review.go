package server

import (
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func handleGetReviews(srv *Server) http.HandlerFunc {
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

		reviews, err := srv.storage.Review.List(limit, offset)
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		jsonResponse(w, r, http.StatusOK, reviews)
	}
}

func handleGetReview(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		review, err := srv.storage.Review.Get(id)
		switch {
		case err == sql.ErrNoRows:
			errorResponse(w, r, http.StatusNotFound, err)
			return
		case err != nil:
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		jsonResponse(w, r, http.StatusOK, review)
	}
}
