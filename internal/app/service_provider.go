package app

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/mikhailsoldatkin/book_store/internal/closer"
	"github.com/mikhailsoldatkin/book_store/internal/config"
)

// serviceProvider holds application dependencies.
type serviceProvider struct {
	config   *config.Config
	dbClient *gorm.DB
}

// newServiceProvider creates a new instance of serviceProvider.
func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// Config returns the loaded application configuration.
func (s *serviceProvider) Config() *config.Config {
	if s.config == nil {
		cfg, err := config.Load()
		if err != nil {
			log.Fatal(err)
		}
		s.config = cfg
	}

	return s.config
}

// DBClient returns the database client instance.
func (s *serviceProvider) DBClient() *gorm.DB {
	if s.dbClient == nil {
		s.initDBClient()
	}
	return s.dbClient
}

// initDBClient initializes the database connection.
func (s *serviceProvider) initDBClient() {
	db, err := gorm.Open(postgres.Open(s.Config().DB.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get database instance: %v", err)
	}

	closer.Add(sqlDB.Close)

	s.dbClient = db
	log.Println("connected to database successfully")
}
