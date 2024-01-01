package databaseConnector

// DatabaseConnector is an interface for interacting with a generic database.
type DatabaseConnector interface {
	// Init initializes the database connection.
	Init() error

	// Close closes the database connection.
	Close() error

	// GetLink retrieves the original link associated with the given short link.
	GetLink(shortLink string) (string, error)

	// GetShortLink retrieves the short link associated with the given link.
	GetShortLink(link string) (string, error)

	// InsertLink inserts a new short link and its associated original link into the database.
	InsertLink(shortLink string, link string) error

	// GetNextSeqNumber retrieves the next sequence number from the database.
	GetNextSeqNumber() (int64, error)
}

// NewPSQLConnector returns a new instance of a PostgreSQL database connector.
func NewPSQLConnector(config Config) DatabaseConnector {
	return &psqlConnector{config: config}
}
