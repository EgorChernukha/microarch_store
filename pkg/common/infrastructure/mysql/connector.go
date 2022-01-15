package mysql

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql"                   // provides MySQL driver
	_ "github.com/golang-migrate/migrate/v4/source/file" // provides filesystem source
)

const (
	dbDriverName = "mysql"
)

type DSN struct {
	User     string
	Password string
	Host     string
	Database string
}

type Client interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)

	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row

	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
}

type Transaction interface {
	Client
	Commit() error
	Rollback() error
}

type TransactionalClient interface {
	Client
	BeginTransaction() (Transaction, error)
}

type DBClient interface {
	TransactionalClient
	Close() error
	Ping() error
}

func NewClient(db *sqlx.DB) DBClient {
	return &client{DB: db}
}

func (dsn *DSN) String() string {
	return fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true", dsn.User, dsn.Password, dsn.Host, dsn.Database)
}

type Config struct {
	MaxConnections     int
	ConnectionLifetime time.Duration
	ConnectTimeout     time.Duration // 0 means default timeout (15 seconds)
}

type Connector interface {
	Open(dsn DSN, cfg Config) error
	MigrateUp(dsn DSN, migrationsDir string) error
	Client() TransactionalClient
	Close() error
}

func NewConnector() Connector {
	return &connector{}
}

type connector struct {
	db DBClient
}

type dbClient interface {
	SetMaxOpenConns(maxConnections int)
	SetConnMaxLifetime(d time.Duration)
	Ping() error
	Close() error
}

type client struct {
	*sqlx.DB
}

func (c *client) BeginTransaction() (Transaction, error) {
	return c.Beginx()
}

func (c *connector) Open(dsn DSN, cfg Config) error {
	var err error
	db, err := openDBX(dsn, cfg)
	if err != nil {
		return errors.Wrap(err, "failed to open database")
	}
	c.db = NewClient(db)
	return errors.WithStack(err)
}

func (c *connector) MigrateUp(dsn DSN, migrationsDir string) (err error) {
	// Db connections will be closed when migration object is closed, so new connection must be opened
	db, err := openDB(dsn, Config{MaxConnections: 1, ConnectionLifetime: time.Minute})
	if err != nil {
		return errors.WithStack(err)
	}
	m, err := createMigrator(db, migrationsDir)
	if err != nil {
		return errors.WithStack(err)
	}
	// noinspection GoUnhandledErrorResult
	defer m.Close()

	err = m.Up()
	if err == migrate.ErrNoChange {
		return nil
	}

	return errors.Wrap(err, "failed to migrate")
}

func (c *connector) Client() TransactionalClient {
	return c.db
}

func (c *connector) Close() error {
	err := c.db.Close()
	return errors.Wrap(err, "failed to disconnect")
}

func createMigrator(db *sql.DB, migrationsDir string) (*migrate.Migrate, error) {
	migrationsURL, err := makeMigrationsURL(migrationsDir)
	if err != nil {
		return nil, err
	}
	driver, err := createMigrationDriver(db)
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsURL, dbDriverName, driver)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create migrator")
	}
	return m, nil
}

func makeMigrationsURL(migrationsDir string) (string, error) {
	// if already url with scheme just return
	if u, err := url.Parse(migrationsDir); err == nil && u.Scheme != "" {
		return migrationsDir, nil
	}

	_, err := os.Stat(migrationsDir)
	if err != nil {
		return "", errors.Wrapf(err, "cannot use migrations from %s", migrationsDir)
	}
	migrationsDir, err = filepath.Abs(migrationsDir)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return fmt.Sprintf("file://%s", migrationsDir), nil
}

func createMigrationDriver(db *sql.DB) (driver database.Driver, err error) {
	err = backoff.Retry(func() error {
		var tryError error
		driver, tryError = mysql.WithInstance(db, &mysql.Config{})
		return tryError
	}, newExponentialBackOff(0))
	return driver, errors.Wrapf(err, "cannot create migrations driver")
}

func openDB(dsn DSN, cfg Config) (*sql.DB, error) {
	db, err := sql.Open(dbDriverName, dsn.String())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open database")
	}
	err = setupDB(db, cfg)
	return db, errors.WithStack(err)
}

func openDBX(dsn DSN, cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open(dbDriverName, dsn.String())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open database")
	}
	err = setupDB(db, cfg)
	return db, errors.WithStack(err)
}

func setupDB(db dbClient, cfg Config) error {
	// Limit max connections count,
	//  next goroutine will wait once reached limit.
	db.SetMaxOpenConns(cfg.MaxConnections)
	// Limits the maximum amount of time the connection may be reused
	// This value must be lower than wait_timeout value on MySQL
	db.SetConnMaxLifetime(cfg.ConnectionLifetime)

	err := backoff.Retry(func() error {
		tryError := db.Ping()
		return tryError
	}, newExponentialBackOff(cfg.ConnectTimeout))
	if err != nil {
		dbCloseErr := db.Close()
		if dbCloseErr != nil {
			err = errors.Wrap(err, dbCloseErr.Error())
		}
		return errors.Wrapf(err, "failed to ping database")
	}
	return nil
}

func newExponentialBackOff(timeout time.Duration) *backoff.ExponentialBackOff {
	exponentialBackOff := backoff.NewExponentialBackOff()
	const maxReconnectWaitingTime = 15 * time.Second
	if timeout != 0 {
		exponentialBackOff.MaxElapsedTime = timeout
	} else {
		exponentialBackOff.MaxElapsedTime = maxReconnectWaitingTime
	}
	return exponentialBackOff
}
