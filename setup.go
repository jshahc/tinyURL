package main

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"tinyURL/databaseConnector"
	"tinyURL/linkGenerator"
	"tinyURL/server"
)

type Config struct {
	Server        server.Config
	Database      databaseConnector.Config
	LinkGenerator linkGenerator.Config
}

func (s *TinyURL) readConfigFile() error {
	// Define a command-line flag for the config file path
	configFile := flag.String("config", "./resources/config/local-run-config.yaml", "Path to the config file")
	flag.Parse()

	viper.SetConfigFile(*configFile)
	// Read the config file
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
		return err
	}

	// Unmarshal the config file into a Config struct
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Error unmarshalling config: %s", err)
		return err
	}
	s.config = config
	return nil
}

// Setup initializes the components of the TinyURL application.
func (s *TinyURL) Setup() error {
	err := s.readConfigFile()
	if err != nil {
		return fmt.Errorf("unable to read config: %s", err)
	}

	// Initialize the ShortLinkGenerator based on config provided
	linkGeneratorConfig := s.config.LinkGenerator
	if linkGeneratorConfig.GeneratorType == "SHA256" {
		s.ShortLinkGenerator = &linkGenerator.SHA256Generator{ShortLinkSize: linkGeneratorConfig.ShortLinkSize, BaseSize: linkGeneratorConfig.Base}
	} else {
		s.ShortLinkGenerator = &linkGenerator.SeqGenerator{BaseSize: linkGeneratorConfig.Base, Counter: linkGeneratorConfig.StartingNumber}
	}

	// Initialize the database connector (using a PostgreSQL connector).
	s.URLStore = databaseConnector.NewPSQLConnector(s.config.Database)
	if err := s.URLStore.Init(); err != nil {
		return fmt.Errorf("unable to start database: %s", err)
	}

	// Initialize the HTTP server with the configured components.
	s.Server = server.HTTPServer{URLStore: s.URLStore, ShortLinkGenerator: s.ShortLinkGenerator, Config: s.config.Server}
	s.Server.Init()
	return nil
}
