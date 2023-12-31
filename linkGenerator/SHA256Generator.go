package linkGenerator

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"tinyURL/databaseConnector"
)

var (
	maxRetries = 10
)

// SHA256Generator is a LinkGenerator implementation using SHA256 hashing.
type SHA256Generator struct {
	ShortLinkSize int
	BaseSize      int
}

// generateRandomLink generates a short link based on the SHA256 hash of the input.
func (g *SHA256Generator) generateRandomLink(input string) string {
	hash := sha256.Sum256([]byte(input))
	hashBase64 := base64.URLEncoding.EncodeToString(hash[:])
	return hashBase64[:g.ShortLinkSize]
}

// GenerateLink creates or updates a short link in the specified database,
// handling collisions by updating a random character in the short link.
func (g *SHA256Generator) GenerateLink(link string, db databaseConnector.DatabaseConnector) (string, error) {
	shortLink := g.generateRandomLink(link)
	log.Default().Printf("Inserting %v for %v", shortLink, link)
	i := 0
	for i = 0; i < maxRetries; i++ {
		err := db.InsertLink(shortLink, link)
		if err != nil {
			log.Default().Printf("Collision Detected: Error inserting %v: %v for %v", shortLink, err, link)
			// Collision detected, update a random character in the short link
			positionToUpdate := GenerateRandomInt(0, g.ShortLinkSize-1)
			characterToUse := GenerateRandomInt(0, g.BaseSize-1)

			// Update the character in the short link
			shortLink = shortLink[:positionToUpdate] + string(CharArray[characterToUse]) + shortLink[positionToUpdate+1:]
		} else {
			log.Default().Printf("Successfully inserted %v for %v", shortLink, link)
			return shortLink, nil
		}
	}
	return "", fmt.Errorf("unable to create short link")
}
