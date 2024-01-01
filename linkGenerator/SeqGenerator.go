package linkGenerator

import (
	"fmt"
	"sync"
	"tinyURL/databaseConnector"
)

// SeqGenerator is a LinkGenerator implementation using sequential numbers and base64 encoding
type SeqGenerator struct {
	BaseSize int        // Base size for base64 conversion
	lock     sync.Mutex // Mutex for thread safety
	Counter  int64      // Counter for generating unique sequence numbers, used only when counter is not partitioned
}

// GenerateLink creates a new short link for the given link and updates the database
func (g *SeqGenerator) GenerateLink(link string, db databaseConnector.DatabaseConnector) (string, error) {
	g.lock.Lock()
	defer g.lock.Unlock()

	// Check if link already exists in database, if yes return corresponding existing short link
	existingShortLink, err := db.GetShortLink(link)
	if err == nil && existingShortLink != "" {
		return existingShortLink, nil
	}

	// Link doesn't exist, create a new short link for it
	seqNumber, err := g.getSeqNumber(db)
	if err != nil {
		return "", fmt.Errorf("error getting next sequence number: %s", err)
	}
	shortLink := g.intToBase64(seqNumber)

	// Insert the new link into the database
	err = db.InsertLink(shortLink, link)
	return shortLink, err
}

// intToBase64 converts an integer to base64 using the specified character array
func (g *SeqGenerator) intToBase64(number int64) string {
	var base64Bytes []byte
	for number > 0 {
		// Value to add is remainder with base size
		remainder := number % int64(g.BaseSize)
		// Add to remainder to start of array
		base64Bytes = append([]byte{CharArray[remainder]}, base64Bytes...)
		// Continue loop by dividing by base size
		number /= int64(g.BaseSize)
	}

	return string(base64Bytes)
}

// getSeqNumber increments and returns the current sequence number
func (g *SeqGenerator) getSeqNumber(db databaseConnector.DatabaseConnector) (int64, error) {
	seqNumber, err := db.GetNextSeqNumber()
	if err == nil {
		return seqNumber, nil
	}
	return -1, err
}
