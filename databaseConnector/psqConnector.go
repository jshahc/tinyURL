package databaseConnector

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/patrickmn/go-cache"
)

// ShortLinkSchema represents the short_links table structure.
type ShortLinkSchema struct {
	ShortLink    string `db:"short_link"`
	OriginalLink string `db:"original_link"`
}

// psqlConnector is a PostgreSQL database connector implementing the DatabaseConnector interface.
type psqlConnector struct {
	Database *sqlx.DB    // Database connection
	Cache    cache.Cache // In-memory cache, Both for ShortLink -> OriginalLink and OriginalLink -> ShortLink
	config   Config      // Database configuration
}

// Config represents the configuration of a PostgreSQL database.
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
	Cache    CacheConfig
}

type CacheConfig struct {
	Enable         bool
	ExpirationTime time.Duration
}

// Init initializes the PostgreSQL database connection.
func (psql *psqlConnector) Init() error {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		psql.config.Host, psql.config.Port, psql.config.Username, psql.config.Password, psql.config.DBName)

	// Open a PostgreSQL database connection using sqlx.
	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		return err
	}
	psql.Database = db

	// Initialize a cache for storing short links.
	if psql.config.Cache.Enable {
		psql.Cache = *cache.New(psql.config.Cache.ExpirationTime, psql.config.Cache.ExpirationTime)
	}

	return nil
}

// GetLink retrieves the original link associated with the given short link from the database.
func (psql *psqlConnector) GetLink(shortLink string) (string, error) {

	// Check if short link is in cache
	cacheEnabled := psql.config.Cache.Enable
	if cacheEnabled {
		if data, found := psql.Cache.Get(shortLink); found {
			return data.(string), nil
		}
	}

	// Execute query on database
	var data ShortLinkSchema
	query := "SELECT original_link FROM short_links WHERE short_link = $1 LIMIT 1"
	err := psql.Database.Get(&data, query, shortLink)
	if err != nil {
		return "", err
	}

	// Set Value in cache
	if cacheEnabled {
		psql.Cache.Set(shortLink, data.OriginalLink, psql.config.Cache.ExpirationTime)
	}

	return data.OriginalLink, nil
}

// GetShortLink retrieves the short link associated with the given link from the database.
func (psql *psqlConnector) GetShortLink(link string) (string, error) {

	// Check if original link is in cache
	cacheEnabled := psql.config.Cache.Enable
	if cacheEnabled {
		if data, found := psql.Cache.Get(link); found {
			return data.(string), nil
		}
	}

	// Execute query on database
	var data ShortLinkSchema
	query := "SELECT short_link FROM short_links WHERE original_link = $1 LIMIT 1"
	err := psql.Database.Get(&data, query, link)
	if err != nil {
		return "", err
	}

	// Set Value in cache
	if cacheEnabled {
		psql.Cache.Set(link, data.ShortLink, psql.config.Cache.ExpirationTime)
	}

	return data.ShortLink, nil
}

// InsertLink inserts a new short link and its associated original link into the database.
func (psql *psqlConnector) InsertLink(shortLink string, link string) error {
	query := "INSERT INTO short_links (short_link, original_link) VALUES ($1, $2)"
	_, err := psql.Database.Exec(query, shortLink, link)
	return err
}

// Close closes the PostgreSQL database connection.
func (psql *psqlConnector) Close() error {
	err := psql.Database.Close()
	if err != nil {
		return err
	}
	return nil
}
