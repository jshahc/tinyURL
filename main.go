package main

import (
	"fmt"
	"log"
	"tinyURL/databaseConnectors"
	"tinyURL/linkGenerator"
	"tinyURL/server"
)

// TinyURL represents the main application structure.
type TinyURL struct {
	URLStore           databaseConnectors.DatabaseConnector
	ShortLinkGenerator linkGenerator.LinkGenerator
	Server             server.HTTPServer
}

// Setup initializes the components of the TinyURL application.
func (s *TinyURL) Setup() error {
	// s.ShortLinkGenerator = &linkGenerator.SeqGenerator{ShortLinkSize: 10, BaseSize: 64}

	// Initialize the ShortLinkGenerator with Sequential Generator and starting value of counter 10000000
	s.ShortLinkGenerator = &linkGenerator.SeqGenerator{BaseSize: 64, Counter: 10000000}

	// Initialize the database connector (using a PostgreSQL connector).
	s.URLStore = databaseConnectors.NewPSQLConnector()
	if err := s.URLStore.Init(); err != nil {
		return fmt.Errorf("unable to start database: %s", err)
	}

	// Initialize the HTTP server with the configured components.
	s.Server = server.HTTPServer{URLStore: s.URLStore, ShortLinkGenerator: s.ShortLinkGenerator}
	s.Server.Init()
	return nil
}

// Run starts the TinyURL application by running the HTTP server.
func (s *TinyURL) Run() error {
	err := s.Server.Run()
	if err != nil {
		return fmt.Errorf("error running URL shortener: %s", err)
	}
	return nil
}

func main() {
	// Create an instance of TinyURL.
	service := TinyURL{}

	// Set up the TinyURL application.
	err := service.Setup()
	if err != nil {
		log.Fatalf("Error setting up URL shortener: %s", err)
	}

	// Run the TinyURL application.
	err = service.Run()
	if err != nil {
		log.Fatalf("Error running URL shortener: %s", err)
	}
}
