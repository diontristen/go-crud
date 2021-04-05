package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"encoding/json"

	"github.com/golang-migrate/migrate/v4"
	pgmigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // This is needed to apply migrations
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

// USAGE: Initialize the connection to the database.
func newDatabase(connString string) (*sqlx.DB, error) {
	//
	// database/sql vs native pgx interface:
	// - still no viable alternative to github.com/jmoiron/sqlx for native pgx
	//   (named parameters instead of $1, $2, ...)
	// So: database/sql is the choice
	//

	// https://medium.com/avitotech/how-to-work-with-postgres-in-go-bad2dabd13e4
	// The checklist:
	// - github.com/jackc/pgx
	// - Configure limits for the connection pool size
	// - TODO: Collect the connection pool metrics with DB.Stats()
	// - TODO: Log what is happening in the driver
	// - Use the Simple Query mode to avoid problems with prepared statements in the transactional mode PgBouncer
	// - Use an up-to-date version of PgBouncer
	// - Not to use request cancellation from the application side

	connConfig, err := pgx.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf(`pgx.ParseConfig error for "%s": %w`, connString, err)
	}

	// :TRICKY: without this one on, the error arises on requests with arguments like
	// "could not determine data type of parameter $1 (SQLSTATE 42P18)"
	connConfig.RuntimeParams["standard_conforming_strings"] = "on"
	connConfig.PreferSimpleProtocol = true

	connStr := stdlib.RegisterConnConfig(connConfig)

	driverName := "pgx"

	db, err := sqlx.Open(driverName, connStr)
	if err != nil {
		return nil, err
	}

	maxDBConnections := 100
	db.SetMaxOpenConns(maxDBConnections)
	db.SetMaxIdleConns(maxDBConnections)
	db.SetConnMaxLifetime(5 * time.Minute)

	migrationError := upMigrations(db)

	if migrationError != nil {
		return nil, migrationError
	}

	return db, nil
}

// USAGE: Migrates the SQL File to the database.
func upMigrations(db *sqlx.DB) error {
	driver, driverError := pgmigrate.WithInstance(db.DB, &pgmigrate.Config{})
	if driverError != nil {
		return driverError
	}

	m, migrateError := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)

	if migrateError != nil {
		return migrateError
	}

	upError := m.Up()
	if upError != nil && upError != migrate.ErrNoChange {
		return upError
	}
	return nil
}

// AppContext - App context type to be used with API handlers to integrate with database, configuration, and other stuff.
type AppContext struct {
	DB     *sqlx.DB
	Config Config
}

// Config is a type to storing application configuration
// This is base on the project requirements
type Config struct {
	PostgresConn      string `json:"postgres_conn"`
	EmailFrom         string `json:"email_from"`
	EmailPassword     string `json:"email_password"`
	EmailSMTPHost     string `json:"email_smtp_host"`
	EmailStmpPort     string `json:"email_stmp_port"`
	EmailSMTPUsername string `json:"email_smtp_username"`
	SecretKey         string `json:"secret_key"`

	AppURL string  `json:"app_url"`
	Listen *string `json:"listen"` // address and port to listen to, default to :5000

	SentryDSN string `json:"sentry_dsn"`
	DebugFlag bool   `json:"debug_flag"`
}

// USAGE: Initialize the Config from the config.json
func NewConfig() (Config, error) {
	config := Config{}

	jsonFile, jsonError := os.Open("config.json")

	if jsonError != nil {
		return config, jsonError
	}

	defer jsonFile.Close()

	byteValue, byteError := ioutil.ReadAll(jsonFile)

	if byteError != nil {
		return config, byteError
	}

	if err := json.Unmarshal(byteValue, &config); err != nil {
		return config, err
	}

	return config, nil
}

// NewAppContext creates a new application context to be used in API
// App Context contains the DB Instance and Config from config.json
func NewAppContext(config Config) (*AppContext, error) {
	db, dbError := newDatabase(config.PostgresConn)
	if dbError != nil {
		return nil, dbError
	}

	return &AppContext{
		DB:     db,
		Config: config,
	}, nil
}

// USAGE: Close the Instance of the DB
func (ac *AppContext) Close() {
	ac.DB.Close()
}

// USAGE: Retrieve the APP Secret Key for JWT
// USE this when we need authentication
func (ac *AppContext) GetSecretKey() []byte {
	return []byte(ac.Config.SecretKey)
}
