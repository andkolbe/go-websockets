package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/jackc/pgx/v4"
)

// how we connect our app to the database
// create a struct so we can add other dbs in the future if we want
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

// create a database pool for postgres
func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	// set some parameters on the db connection pool that will stop it from growing out of control
	d.SetMaxOpenConns(10) // never have more than 10 db conncetions open at a time
	d.SetMaxIdleConns(5)
	d.SetConnMaxLifetime(5 * time.Minute) // 5 min

	dbConn.SQL = d

	err = testDB(d) // tries to ping the database
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

// tries to ping the database
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

func NewDatabase(dsn string) (*sql.DB, error) {
	// connect to db
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// tests database connection
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}