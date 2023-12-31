package main

import (
	"log"
	"tinyURL/databaseConnector"
	"tinyURL/linkGenerator"
	"tinyURL/server"
)

// TinyURL represents the main application structure.
type TinyURL struct {
	URLStore           databaseConnector.DatabaseConnector
	ShortLinkGenerator linkGenerator.LinkGenerator
	Server             server.HTTPServer
	config             Config
}

// Run starts the TinyURL application by running the HTTP server.
func (s *TinyURL) Run() error {
	err := s.Server.Run()
	if err != nil {
		return err
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
		log.Fatalf("Error running URL shortener: %v", err)
	}
}
