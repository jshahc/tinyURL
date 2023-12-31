package linkGenerator

import "tinyURL/databaseConnector"

// LinkGenerator is an interface for generating and managing short links.
type LinkGenerator interface {
	// GenerateLink creates or updates a short link in the specified database.
	GenerateLink(input string, db databaseConnector.DatabaseConnector) (string, error)
}

// Config represents the configuration for the link generator.
type Config struct {
	GeneratorType  string
	Base           int
	StartingNumber int
	ShortLinkSize  int
}
