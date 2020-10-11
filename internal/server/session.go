package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func handleAuthorize(srv *Server) http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := new(request)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		user, err := srv.storage.User.GetByEmail(req.Email)
		if err == sql.ErrNoRows || !user.CheckPasswordHash(req.Password) {
			errorResponse(w, r, http.StatusNotFound, errUserNotFound)
			return
		}
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		session, err := srv.sessionStore.Get(r, "SESSID")
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = user.ID
		if err := srv.sessionStore.Save(r, w, session); err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		noContentResponse(w, r, http.StatusNoContent)
	}
}

func handleLogout(srv *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := srv.sessionStore.Get(r, "SESSID")
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		if _, ok := session.Values["user_id"]; !ok {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		delete(session.Values, "user_id")

		if err := srv.sessionStore.Save(r, w, session); err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		noContentResponse(w, r, http.StatusNoContent)
	}
}
