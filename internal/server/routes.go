package server

func setRoutes(srv *Server) {
	r := srv.router

	authorized := r.PathPrefix("/").Subrouter()
	authorized.Use(srv.authMiddleware)

	// BOOKS
	r.HandleFunc("/books", handleGetBooks(srv)).Methods("GET")
	r.HandleFunc("/books", handleCreateBook(srv)).Methods("POST")
	r.HandleFunc("/books/{id:[0-9]+}", handleGetBook(srv)).Methods("GET")
	r.HandleFunc("/books/{id:[0-9]+}", handleUpdateBook(srv)).Methods("PATCH")
	r.HandleFunc("/books/{id:[0-9]+}", handleDeleteBook(srv)).Methods("DELETE")
	r.HandleFunc("/books/{id:[0-9]+}/author", handleGetBookAuthor(srv)).Methods("GET")
	r.HandleFunc("/books/{id:[0-9]+}/tags", handleGetBookTags(srv)).Methods("GET")
	r.HandleFunc("/books/{id:[0-9]+}/quotes", handleGetBookQuotes(srv)).Methods("GET")
	r.HandleFunc("/books/{id:[0-9]+}/reviews", handleGetBookReviews(srv)).Methods("GET")
	r.HandleFunc("/books/{id:[0-9]+}/tags", handleAddTagToBook(srv)).Methods("POST")
	r.HandleFunc("/books/{id:[0-9]+}/tags/{tag_id:[0-9]+}", handleRemoveBookTag(srv)).Methods("POST")

	// QUOTES
	r.HandleFunc("/quotes", handleGetQuotes(srv)).Methods("GET")
	r.HandleFunc("/quotes", handleCreateQuote(srv)).Methods("POST")
	r.HandleFunc("/quotes/{id:[0-9]+}", handleGetQuote(srv)).Methods("GET")
	r.HandleFunc("/quotes/{id:[0-9]+}", handleUpdateQuote(srv)).Methods("PATCH")
	r.HandleFunc("/quotes/{id:[0-9]+}", handleDeleteQuote(srv)).Methods("DELETE")

	// REVIEWS
	r.HandleFunc("/reviews", handleGetReviews(srv)).Methods("GET")
	r.HandleFunc("/reviews/{id:[0-9]+}", handleGetReview(srv)).Methods("GET")

	// TAGS
	r.HandleFunc("/tags", handleGetTags(srv)).Methods("GET")
	r.HandleFunc("/tags", handleCreateTag(srv)).Methods("POST")
	r.HandleFunc("/tags/{id:[0-9]+}", handleGetTag(srv)).Methods("GET")
	r.HandleFunc("/tags/{id:[0-9]+}", handleUpdateTag(srv)).Methods("PATCH")
	r.HandleFunc("/tags/{id:[0-9]+}", handleDeleteTag(srv)).Methods("DELETE")
	r.HandleFunc("/tags/{id:[0-9]+}/books", handleGetTagBooks(srv)).Methods("GET")

	// AUTHORS
	r.HandleFunc("/authors", handleGetAuthors(srv)).Methods("GET")
	r.HandleFunc("/authors", handleCreateAuthor(srv)).Methods("POST")
	r.HandleFunc("/authors/{id:[0-9]+}", handleGetAuthor(srv)).Methods("GET")
	r.HandleFunc("/authors/{id:[0-9]+}", handleUpdateAuthor(srv)).Methods("PATCH")
	r.HandleFunc("/authors/{id:[0-9]+}", handleDeleteAuthor(srv)).Methods("DELETE")
	r.HandleFunc("/authors/{id:[0-9]+}/books", handleGetAuthorBooks(srv)).Methods("GET")
	r.HandleFunc("/authors/{id:[0-9]+}/quotes", handleGetAuthorQuotes(srv)).Methods("GET")

	// USERS
	r.HandleFunc("/users", handleGetUsers(srv)).Methods("GET")
	r.HandleFunc("/users", handleCreateUser(srv)).Methods("POST")
	r.HandleFunc("/users/{id:[0-9]+}", handleGetUser(srv)).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", handleUpdateUser(srv)).Methods("PATCH")
	r.HandleFunc("/users/{id:[0-9]+}", handleDeleteUser(srv)).Methods("DELETE")
	r.HandleFunc("/users/{id:[0-9]+}/books", handleGetUserBooks(srv)).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}/quotes", handleGetUserQuotes(srv)).Methods("GET")   // implement
	r.HandleFunc("/users/{id:[0-9]+}/reviews", handleGetUserReviews(srv)).Methods("GET") // implement

	authorized.HandleFunc("/me", handleGetMe(srv)).Methods("GET")

	// SESSIONS
	r.HandleFunc("/sessions", handleAuthorize(srv)).Methods("POST")
	authorized.HandleFunc("/sessions", handleLogout(srv)).Methods("DELETE")
}
