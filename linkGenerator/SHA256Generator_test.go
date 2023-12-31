package linkGenerator

import (
	"testing"
	"tinyURL/databaseConnector"
)

func TestSHA256Generator_GenerateLink(t *testing.T) {
	mockDB := &databaseConnector.MockDatabaseConnector{Data: make(map[string]string)}
	gen := &SHA256Generator{ShortLinkSize: 10, BaseSize: 64}

	// Test new link generation
	link := "https://example.com"
	shortLink, err := gen.GenerateLink(link, mockDB)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check if the link was correctly inserted into the mock database
	storedLink, err := mockDB.GetLink(shortLink)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if storedLink != link {
		t.Errorf("Expected %s, got %s", link, storedLink)
	}
}
