package databaseConnector

type MockDatabaseConnector struct {
	Data map[string]string
}

func (m *MockDatabaseConnector) Init() error {
	return nil
}

func (m *MockDatabaseConnector) Close() error {
	return nil
}

func (m *MockDatabaseConnector) GetLink(shortLink string) (string, error) {
	// Check if short link exists in database, if yes return corresponding existing link
	for link, val := range m.Data {
		if val == shortLink {
			return link, nil
		}
	}
	return "", nil
}

func (m *MockDatabaseConnector) GetShortLink(link string) (string, error) {
	if val, ok := m.Data[link]; ok {
		return val, nil
	}
	return "", nil
}

func (m *MockDatabaseConnector) InsertLink(shortLink, link string) error {
	m.Data[link] = shortLink
	return nil
}
