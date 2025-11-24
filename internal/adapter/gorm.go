package adapter

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GormAdapter is an adapter for GORM ORM
type GormAdapter struct {
	db *gorm.DB
}

// PostgresConfig holds the configuration for PostgreSQL connection
type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewGormAdapter creates a new GormAdapter instance
func NewGormAdapter(db *gorm.DB) *GormAdapter {
	return &GormAdapter{
		db: db,
	}
}

// NewPostgresGormAdapter creates a new GormAdapter instance connected to a PostgreSQL database
func NewPostgresGormAdapter(config PostgresConfig) (*GormAdapter, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &GormAdapter{db: db}, nil
}

// FindAll retrieves all records that match the given conditions into dest with optional limit and offset
// limit and offset can be set to 0 for no limit/offset
func (g *GormAdapter) FindAll(dest interface{}, limit int, offset int, conditions ...interface{}) error {
	query := g.db
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	return query.Find(dest, conditions...).Error
}

// Close closes the database connection
func (g *GormAdapter) Close() error {
	sqlDB, err := g.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
