package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"tinyURL/databaseConnectors"
	"tinyURL/linkGenerator"
)

// HTTPServer represents the HTTP server for handling URL shortening operations.
type HTTPServer struct {
	router             *mux.Router
	URLStore           databaseConnectors.DatabaseConnector
	ShortLinkGenerator linkGenerator.LinkGenerator
}

// Init initializes the HTTP server and sets up the routing.
func (s *HTTPServer) Init() {
	s.router = mux.NewRouter()
	s.router.HandleFunc("/{param}", s.getHandler).Methods("GET")
	s.router.HandleFunc("/", s.postHandler).Methods("POST")
}

// Run starts the HTTP server on port 8080.
func (s *HTTPServer) Run() error {
	err := http.ListenAndServe(":8080", s.router)
	if err != nil {
		return fmt.Errorf("unable to start server: %s", err)
	}
	return nil
}
