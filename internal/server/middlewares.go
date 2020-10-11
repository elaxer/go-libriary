package server

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{ResponseWriter: w}

		start := time.Now()
		next.ServeHTTP(rw, r)
		completed := time.Now().Sub(start)

		log.Printf(
			"%s %s: %s(%d) - %s\n",
			r.Method,
			r.URL.String(),
			http.StatusText(rw.status),
			rw.status,
			completed,
		)
	})
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, "SESSID")
		if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		userID, ok := session.Values["user_id"]
		if !ok {
			errorResponse(w, r, http.StatusUnauthorized, errUnauthorized)
			return
		}

		user, err := s.storage.User.Get(userID.(int))
		if err == sql.ErrNoRows {
			errorResponse(w, r, http.StatusNotFound, err)
			return
		} else if err != nil {
			errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		ctx := context.WithValue(r.Context(), ctxKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
