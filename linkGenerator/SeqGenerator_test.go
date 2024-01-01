package linkGenerator

import (
	"testing"
	"tinyURL/databaseConnector"
)

func TestSeqGenerator_GenerateLink(t *testing.T) {
	mockDB := &databaseConnector.MockDatabaseConnector{Data: make(map[string]string), Counter: 10000000}
	gen := &SeqGenerator{BaseSize: 64, Counter: 10000000}

	// Test new link generation
	link := "https://example.com"
	expectedShortLink := "c9Q1"
	shortLink, err := gen.GenerateLink(link, mockDB)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if shortLink != expectedShortLink {
		t.Errorf("Expected %s, got %s", expectedShortLink, shortLink)
	}

	// Test existing link retrieval
	existingShortLink, err := gen.GenerateLink(link, mockDB)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if existingShortLink != shortLink {
		t.Errorf("Expected %s, got %s", shortLink, existingShortLink)
	}
}

func TestSeqGenerator_intToBase64(t *testing.T) {
	gen := &SeqGenerator{BaseSize: 64}
	examples := map[int64]string{
		1:               "1",
		1_000_000_000:   "xcie0",
		100_000_000_000: "1T8TkW0",
	}
	for number, expected := range examples {
		result := gen.intToBase64(number)
		if result != expected {
			t.Errorf("Expected %s, got %s. Testing for: %d", expected, result, number)
		}
	}
}

func TestSeqGenerator_getSeqNumber(t *testing.T) {
	mockDB := &databaseConnector.MockDatabaseConnector{Data: make(map[string]string), Counter: 10000000}
	gen := &SeqGenerator{BaseSize: 64, Counter: 10000000}

	expectedSeqNumber := int64(10000001)
	seqNumber, err := gen.getSeqNumber(mockDB)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if seqNumber != expectedSeqNumber {
		t.Errorf("Expected %d, got %d", expectedSeqNumber, seqNumber)
	}
}
