package server

import (
	"books/internal/config"
	"books/internal/storage"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"strconv"
)

type ctxKey int8

const (
	ctxKeyUser ctxKey = iota
)

var (
	errUnauthorized = errors.New("authorize to execute this operation")
	errUserNotFound = errors.New("session: user not found")
)

// Server ...
type Server struct {
	config       *config.Config
	router       *mux.Router
	storage      *storage.Storage
	sessionStore sessions.Store
}

// New ...
func New(config *config.Config, storage *storage.Storage, sessionStore sessions.Store) *Server {
	return &Server{
		config:       config,
		router:       mux.NewRouter(),
		storage:      storage,
		sessionStore: sessionStore,
	}
}

// Run ...
func (s *Server) Run() error {
	s.configureRouter()

	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)

	log.Printf("Server is starting at address %s...\n", addr)

	return http.ListenAndServe(addr, s.router)
}

func (s *Server) configureRouter() {
	s.router.Use(s.loggingMiddleware)

	setRoutes(s)
}

func getIntQuery(r *http.Request, name string, _default int) (int, error) {
	if r.URL.Query().Get(name) == "" {
		return _default, nil
	}

	param, err := strconv.Atoi(r.URL.Query().Get(name))
	if err != nil {
		return param, err
	}

	return param, nil
}

func response(w http.ResponseWriter, r *http.Request, code int, body string) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if _, err := fmt.Fprint(w, body); err != nil {

	}
}

func jsonResponse(w http.ResponseWriter, r *http.Request, code int, item interface{}) {
	marshaled, err := json.Marshal(item)
	if err != nil {
		errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	response(w, r, code, string(marshaled))
}

func noContentResponse(w http.ResponseWriter, r *http.Request, code int) {
	response(w, r, code, "")
}

func errorResponse(w http.ResponseWriter, r *http.Request, code int, err error) {
	noContentResponse(w, r, code)
	log.Println(err)
}
