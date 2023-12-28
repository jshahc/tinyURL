package linkGenerator

import "tinyURL/databaseConnectors"

// LinkGenerator is an interface for generating and managing short links.
type LinkGenerator interface {
	// GenerateLink creates or updates a short link in the specified database.
	GenerateLink(input string, db databaseConnectors.DatabaseConnector) (string, error)
}
