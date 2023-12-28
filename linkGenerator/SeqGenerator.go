package linkGenerator

import (
	"sync"
	"tinyURL/databaseConnectors"
)

// Elements of chosen base 64
var charArray = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_")

// SeqGenerator is a LinkGenerator implementation using sequential numbers and base64 encoding
type SeqGenerator struct {
	BaseSize int        // Base size for base64 conversion
	lock     sync.Mutex // Mutex for thread safety
	Counter  int        // Counter for generating unique sequence numbers
}

// GenerateLink creates a new short link for the given link and updates the database
func (g *SeqGenerator) GenerateLink(link string, db databaseConnectors.DatabaseConnector) (string, error) {
	g.lock.Lock()
	defer g.lock.Unlock()

	// Check if link already exists in database, if yes return corresponding existing short link
	existingShortLink, err := db.GetShortLink(link)
	if err == nil && existingShortLink != "" {
		return existingShortLink, nil
	}

	// Link doesn't exist, create a new short link for it
	seqNumber := g.getSeqNumber()
	shortLink := g.intToBase64(seqNumber)

	// Insert the new link into the database
	err = db.InsertLink(shortLink, link)
	return shortLink, err
}

// intToBase64 converts an integer to base64 using the specified character array
func (g *SeqGenerator) intToBase64(number int) string {
	var base64Bytes []byte

	for number > 0 {
		// Value to add is remainder with base size
		remainder := number % g.BaseSize
		// Add to remainder to start of array
		base64Bytes = append([]byte{charArray[remainder]}, base64Bytes...)
		// Continue loop by dividing by base size
		number /= g.BaseSize
	}

	return string(base64Bytes)
}

// getSeqNumber increments and returns the current sequence number
func (g *SeqGenerator) getSeqNumber() int {
	g.Counter += 1
	return g.Counter
}
