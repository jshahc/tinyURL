package linkGenerator

import (
	"sync"
	"tinyURL/databaseConnector"
)

// Maximum number of links supported, used for partitioning
var limit = 100_000_000

// SeqGenerator is a LinkGenerator implementation using sequential numbers and base64 encoding
type SeqGenerator struct {
	BaseSize           int        // Base size for base64 conversion
	lock               sync.Mutex // Mutex for thread safety
	Counter            int        // Counter for generating unique sequence numbers, used only when counter is not partitioned
	PartitionedCounter []int      // Partitioned counter for generating unique sequence numbers
}

func (g *SeqGenerator) Init(numPartitions int) {
	g.PartitionedCounter = make([]int, numPartitions)
	for i := 0; i < numPartitions; i++ {
		g.PartitionedCounter[i] = (i * limit) / numPartitions
	}
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
		base64Bytes = append([]byte{CharArray[remainder]}, base64Bytes...)
		// Continue loop by dividing by base size
		number /= g.BaseSize
	}

	return string(base64Bytes)
}

// getSeqNumber increments and returns the current sequence number
func (g *SeqGenerator) getSeqNumber() int {
	// Partitioning is not enabled
	if len(g.PartitionedCounter) == 0 {
		g.Counter += 1
		return g.Counter
	}

	// Partitioning is enabled, pick a random partition and increment value
	randomPartition := GenerateRandomInt(0, len(g.PartitionedCounter)-1)
	g.PartitionedCounter[randomPartition] += 1
	return g.PartitionedCounter[randomPartition]
}
