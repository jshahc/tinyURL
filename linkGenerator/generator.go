package linkGenerator

import (
	"math/rand"
	"tinyURL/databaseConnector"
)

// CharArray Elements of chosen base 64
var CharArray = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_")

// LinkGenerator is an interface for generating and managing short links.
type LinkGenerator interface {
	// GenerateLink creates or updates a short link in the specified database.
	GenerateLink(input string, db databaseConnector.DatabaseConnector) (string, error)
}

// Config represents the configuration for the link generator.
type Config struct {
	GeneratorType  string
	Base           int
	StartingNumber int64
	ShortLinkSize  int
}

// GenerateRandomInt generates a random integer in the specified range.
func GenerateRandomInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}
